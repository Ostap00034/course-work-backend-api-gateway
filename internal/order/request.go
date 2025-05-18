// internal/order/request.go
package order

type createOrderRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float32 `json:"price" binding:"required"`
	Address     string  `json:"address" binding:"required"`
	Longitude   string  `json:"longitude" binding:"required"`
	Latitude    string  `json:"latitude" binding:"required"`
	CategoryId  string  `json:"category_id"`
	ClientId    string  `json:"client_id" binding:"required"`
}

type getOrdersRequest struct {
	CategoriesIds []string `form:"categories_ids"`
	Status        string   `form:"status"`
	ClientId      string   `form:"client_id"`
	MasterId      string   `form:"master_id"`
}

type updateOrderRequest struct {
	Title       string  `json:"title,omitempty" binding:"omitempty"`
	Description string  `json:"description,omitempty" binding:"omitempty"`
	Price       float32 `json:"price,omitempty" binding:"omitempty"`
	Address     string  `json:"address,omitempty" binding:"omitempty"`
	Longitude   string  `json:"longitude,omitempty" binding:"omitempty"`
	Latitude    string  `json:"latitude,omitempty" binding:"omitempty"`
	Status      string  `json:"status,omitempty" binding:"omitempty"`
	CategoryId  string  `json:"category_id,omitempty" binding:"omitempty"`
	ClientId    string  `json:"client_id,omitempty" binding:"omitempty"`
	MasterId    string  `json:"master_id,omitempty" binding:"omitempty"`
}

type getMyOrdersRequest struct {
	UserId        string   `form:"user_id" binding:"required"`
	Status        string   `form:"status"`
	CategoriesIds []string `form:"categories_ids"`
}

type getMyFinishedOrdersRequest struct {
	UserId string `form:"user_id" binding:"required"`
}
