package dto

type CreateBookRequest struct {
	Name          string  `json:"name" binding:"required"`
	Author        string  `json:"author" binding:"required"`
	Genre         string  `json:"genre" binding:"required"`
	ISBN          string  `json:"isbn" binding:"required"`
	Description   string  `json:"description" binding:"required"`
	PublishedYear uint    `json:"published_year" binding:"required"`
	Stock         uint    `json:"stock" binding:"required"`
	Price         float64 `json:"price" binding:"required"`
}

type UpdateBookRequest struct {
	Name          string  `json:"name" `
	Author        string  `json:"author"`
	Genre         string  `json:"genre"`
	ISBN          string  `json:"isbn"`
	Description   string  `json:"description"`
	PublishedYear uint    `json:"published_year"`
	Stock         uint    `json:"stock"`
	Price         float64 `json:"price"`
}

type BookResponse struct {
	ID            uint    `json:"id"`
	Name          string  `json:"name" `
	Author        string  `json:"author" `
	Genre         string  `json:"genre" `
	ISBN          string  `json:"isbn" `
	Description   string  `json:"description"`
	PublishedYear uint    `json:"published_year"`
	Stock         uint    `json:"stock"`
	Price         float64 `json:"price"`
}
