// internal/offer/ws.go
package offer

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	authpbv1 "github.com/Ostap00034/course-work-backend-api-specs/gen/go/auth/v1"
	offerpbv1 "github.com/Ostap00034/course-work-backend-api-specs/gen/go/offer/v1"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"google.golang.org/grpc/status"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type wsMsg struct {
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data"`
}

type createOfferPayload struct {
	OrderId  string  `json:"order_id"`
	MasterId string  `json:"master_id"`
	Price    float32 `json:"price"`
}

// Новый пэйлоад для обновления предложения
type updateOfferPayload struct {
	OfferId string `json:"offer_id"`
	Status  string `json:"status"`
}

func OfferWsHandler(offerClient offerpbv1.OfferServiceClient, authClient authpbv1.AuthServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "upgrade failed"})
			return
		}
		defer conn.Close()

		// 1) Auth по cookie
		token, err := c.Cookie("token") // убедись, что тут то же имя, что и в SetCookie
		if err != nil {
			conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "no auth"))
			return
		}
		if _, err := authClient.ValidateToken(c, &authpbv1.ValidateTokenRequest{Token: token}); err != nil {
			conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "unauthorized"))
			return
		}

		// 2) Ping/Pong keep-alive
		conn.SetReadDeadline(time.Now().Add(pongWait))
		conn.SetPongHandler(func(string) error {
			conn.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})
		ticker := time.NewTicker(pingPeriod)
		defer ticker.Stop()
		go func() {
			for range ticker.C {
				conn.SetWriteDeadline(time.Now().Add(writeWait))
				if conn.WriteMessage(websocket.PingMessage, nil) != nil {
					return
				}
			}
		}()

		// 3) Основной loop
		for {
			_, raw, err := conn.ReadMessage()
			if err != nil {
				break
			}
			var m wsMsg
			if json.Unmarshal(raw, &m) != nil {
				conn.WriteJSON(gin.H{"error": "invalid format"})
				continue
			}

			switch m.Action {

			case "createOffer":
				var p createOfferPayload
				if json.Unmarshal(m.Data, &p) != nil {
					conn.WriteJSON(gin.H{"action": "createOffer", "error": "bad data"})
					continue
				}
				resp, err := offerClient.CreateOffer(context.Background(), &offerpbv1.CreateOfferRequest{
					OrderId:  p.OrderId,
					MasterId: p.MasterId,
					Price:    p.Price,
				})
				if err != nil {
					st := status.Convert(err)
					conn.WriteJSON(gin.H{"action": "createOffer", "error": st.Message()})
				} else {
					conn.WriteJSON(gin.H{"action": "createOffer", "offer": resp.Offer})
				}

			case "listOffers":
				var tmp struct {
					OrderId string `json:"order_id"`
				}
				if json.Unmarshal(m.Data, &tmp) != nil {
					conn.WriteJSON(gin.H{"action": "listOffers", "error": "bad data"})
					continue
				}
				resp, err := offerClient.GetMyOrderOffers(context.Background(), &offerpbv1.GetMyOrderOffersRequest{
					OrderId: tmp.OrderId,
				})
				if err != nil {
					st := status.Convert(err)
					conn.WriteJSON(gin.H{"action": "listOffers", "error": st.Message()})
				} else {
					conn.WriteJSON(gin.H{"action": "listOffers", "offers": resp.Offers})
				}

			case "updateOffer":
				var p updateOfferPayload
				if json.Unmarshal(m.Data, &p) != nil {
					conn.WriteJSON(gin.H{"action": "updateOffer", "error": "bad data"})
					continue
				}

				resp, err := offerClient.UpdateOffer(context.Background(), &offerpbv1.UpdateOfferRequest{
					Id:     p.OfferId,
					Status: p.Status,
				})

				if err != nil {
					st := status.Convert(err)
					conn.WriteJSON(gin.H{"action": "updateOffer", "error": st.Message()})
				} else {
					conn.WriteJSON(gin.H{"action": "updateOffer", "offer": resp.Offer})
				}

			default:
				conn.WriteJSON(gin.H{"error": "unknown action"})
			}
		}
	}
}
