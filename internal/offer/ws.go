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

// wsMsg описывает общую обёртку для входящих сообщений
type wsMsg struct {
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data"`
}

// createOfferPayload — пэйлоад для создания оффера
type createOfferPayload struct {
	OrderId  string  `json:"order_id"`
	MasterId string  `json:"master_id"`
	Price    float32 `json:"price"`
}

// updateOfferPayload — пэйлоад для обновления оффера
type updateOfferPayload struct {
	OfferId string `json:"offer_id"`
	Status  string `json:"status"`
}

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

// OfferWsHandler возвращает Gin-хендлер WebSocket.
//   - hub        — менеджер подписок, у которого реализованы методы Subscribe, Unsubscribe и Broadcast.
//   - offerClient — gRPC-клиент OfferService.
//   - authClient  — gRPC-клиент AuthService для проверки токена.
func OfferWsHandler(
	hub *Hub,
	offerClient offerpbv1.OfferServiceClient,
	authClient authpbv1.AuthServiceClient,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "upgrade failed"})
			return
		}
		defer conn.Close()

		// 1) Авторизация по cookie "token"
		token, err := c.Cookie("token")
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

		// 2) Настройка ping/pong для keep-alive
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

		// 3) Основной цикл обработки сообщений
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
			// подписаться на обновления конкретного заказа
			case "subscribe":
				var sub struct {
					OrderId string `json:"order_id"`
				}
				if err := json.Unmarshal(m.Data, &sub); err == nil {
					hub.Subscribe(sub.OrderId, conn)
				}

			// создать новый оффер
			case "createOffer":
				var p createOfferPayload
				if err := json.Unmarshal(m.Data, &p); err != nil {
					conn.WriteJSON(gin.H{"action": "createOffer", "error": "bad data"})
					continue
				}
				grpcResp, err := offerClient.CreateOffer(
					context.Background(),
					&offerpbv1.CreateOfferRequest{
						OrderId:  p.OrderId,
						MasterId: p.MasterId,
						Price:    p.Price,
					},
				)
				if err != nil {
					st := status.Convert(err)
					conn.WriteJSON(gin.H{"action": "createOffer", "error": st.Message()})
					continue
				}
				// ответ инициатору
				conn.WriteJSON(gin.H{"action": "createOffer", "offer": grpcResp.Offer})
				// уведомить всех подписчиков заказа
				hub.Broadcast(p.OrderId, gin.H{"action": "offerCreated", "offer": grpcResp.Offer})

			// обновить статус существующего оффера
			case "updateOffer":
				var p updateOfferPayload
				if err := json.Unmarshal(m.Data, &p); err != nil {
					conn.WriteJSON(gin.H{"action": "updateOffer", "error": "bad data"})
					continue
				}
				grpcResp, err := offerClient.UpdateOffer(
					context.Background(),
					&offerpbv1.UpdateOfferRequest{
						Id:     p.OfferId,
						Status: p.Status,
					},
				)
				if err != nil {
					st := status.Convert(err)
					conn.WriteJSON(gin.H{"action": "updateOffer", "error": st.Message()})
					continue
				}
				// ответ инициатору
				conn.WriteJSON(gin.H{"action": "updateOffer", "offer": grpcResp.Offer})
				// и рассылка всем подписчикам по заказу
				hub.Broadcast(p.OfferId, gin.H{"action": "offerUpdated", "offer": grpcResp.Offer})

			default:
				conn.WriteJSON(gin.H{"error": "unknown action"})
			}
		}

		// при отключении клиента отписать от всех заказов
		hub.mu.Lock()
		for orderID := range hub.subs {
			hub.Unsubscribe(orderID, conn)
		}
		hub.mu.Unlock()
	}
}
