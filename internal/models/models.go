package models

import "time"

type Product struct {
	ID    int
	Name  string
	Price float64
	Stock int // количество на складе
}

type User struct {
	ID    int
	Name  string
	Email string
}

type Order struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	TotalPrice float64   `json:"total_price"`
	CreatedAt  time.Time `json:"created_at"`
}

type OrderItem struct {
	ID        int `json:"id"`
	OrderID   int `json:"order_id"`
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type Transaction struct {
	ID        int       `json:"id"`
	OrderID   int       `json:"order_id"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type OrderDetail struct {
	OrderID           int       `json:"order_id"`
	CreatedAt         time.Time `json:"created_at"`
	ProductName       string    `json:"product_name"`
	Quantity          int       `json:"quantity"`
	Price             float64   `json:"price"`
	TotalPrice        float64   `json:"total_price"`
	TransactionStatus string    `json:"transaction_status"`
}

type PopularProduct struct {
	ProductID   int    `json:"product_id"`
	ProductName string `json:"product_name"`
	TotalSold   int    `json:"total_sold"`
}
