package dto

type AddCartItemRequest struct {
	BookId   uint `json:"book_id"  binding:"required"`
	Quantity uint `json:"quantity"  binding:"required"`
}

type CartItem struct {
	Book     BookResponse `json:"book"`
	Quantity uint         `json:"quantity"`
	SubTotal float64      `json:"subTotal"`
}

type CartResponse struct {
	Items []CartItem `json:"items"`
	Total float64    `json:"total"`
}
