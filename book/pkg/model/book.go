package model

// Book defines the book metadata.
type Book struct {
	ID    string `json:"id"`
	Title string `json:"title" binding:"required"`
}

type BookEvent struct {
	ID    string `json:"id"`
	Title string `json:"title" binding:"required"`
}
