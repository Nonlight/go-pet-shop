package postgres

import (
	"context"
	"go-pet-shop/internal/models"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func TestPlaceOrder_Success(t *testing.T) {
	ctx := context.Background()

	if err := godotenv.Load("../../../.env"); err != nil {
		t.Fatal(err)
	}

	databaseURL := os.Getenv("TEST_DATABASE_URL")
	db, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		t.Fatal(err)
	}

	storage := &Storage{db: db}
	if _, err = db.Exec(ctx, `TRUNCATE transactions, order_items, orders, products, users RESTART IDENTITY CASCADE`); err != nil {
		t.Fatal(err)
	}
	if _, err = db.Exec(ctx, `INSERT INTO users (name, email) VALUES ($1, $2)`,
		"Test User", "test@mail.ru"); err != nil {
		t.Fatal(err)
	}

	var productID int
	if err := db.QueryRow(ctx, `INSERT INTO products (name, price, stock) VALUES ($1, $2, $3) RETURNING id`,
		"Test Product", 100, 10).Scan(&productID); err != nil {
		t.Fatal(err)
	}

	orderID, err := storage.PlaceOrder(ctx, "test@mail.ru",
		[]models.OrderItem{{ProductID: productID, Quantity: 3}})
	if err != nil {
		t.Fatal(err)
	}

	var totalPrice float64
	if err := db.QueryRow(ctx, `SELECT total_price FROM orders WHERE id = $1`,
		orderID).Scan(&totalPrice); err != nil {
		t.Fatal(err)
	}
	if totalPrice != 300 {
		t.Errorf("got %v, want %v", totalPrice, 300)
	}

	var stock int
	if err := db.QueryRow(ctx, `SELECT stock FROM products WHERE id = $1`, productID).Scan(&stock); err != nil {
		t.Fatal(err)
	}
	if stock != 7 {
		t.Errorf("got %d, want %d", stock, 7)
	}

	var orderItem models.OrderItem
	if err := db.QueryRow(ctx, `SELECT product_id, quantity FROM order_items WHERE order_id = $1`,
		orderID).Scan(&orderItem.ProductID, &orderItem.Quantity); err != nil {
		t.Fatal(err)
	}
	if orderItem.ProductID != productID {
		t.Errorf("got product_id %d, want %d", orderItem.ProductID, productID)
	}
	if orderItem.Quantity != 3 {
		t.Errorf("got quantity %d, want %d", orderItem.Quantity, 3)
	}

	var transaction struct {
		Amount float64
		Status string
	}

	if err := db.QueryRow(ctx, `SELECT amount, status FROM transactions WHERE order_id = $1`,
		orderID).Scan(&transaction.Amount, &transaction.Status); err != nil {
		t.Fatal(err)
	}
	if transaction.Amount != 300 {
		t.Errorf("got transaction amount %v, want %v", transaction.Amount, 300)
	}
	if transaction.Status != "completed" {
		t.Errorf("got transaction status %q, want %q", transaction.Status, "completed")
	}
}

func TestPlaceOrder_Rollback(t *testing.T) {
	ctx := context.Background()

	if err := godotenv.Load("../../../.env"); err != nil {
		t.Fatal(err)
	}

	databaseURL := os.Getenv("TEST_DATABASE_URL")
	db, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		t.Fatal(err)
	}

	storage := &Storage{db: db}
	if _, err = db.Exec(ctx,
		`TRUNCATE transactions, order_items, orders, products, users RESTART IDENTITY CASCADE`,
	); err != nil {
		t.Fatal(err)
	}
	if _, err = db.Exec(ctx,
		`INSERT INTO users (name, email) VALUES ($1, $2)`,
		"Test User", "email@mail.ru",
	); err != nil {
		t.Fatal(err)
	}

	var product1ID int
	if err := db.QueryRow(ctx, `INSERT INTO products (name, price, stock) VALUES ($1, $2, $3) RETURNING id`,
		"Product 1", 100, 10).Scan(&product1ID); err != nil {
		t.Fatal(err)
	}

	var product2ID int
	if err := db.QueryRow(ctx, `INSERT INTO products (name, price, stock) VALUES ($1, $2, $3) RETURNING id`,
		"Product 2", 200, 1).Scan(&product2ID); err != nil {
		t.Fatal(err)
	}

	_, err = storage.PlaceOrder(ctx, "email@mail.ru", []models.OrderItem{
		{
			ProductID: product1ID,
			Quantity:  2,
		},
		{
			ProductID: product2ID,
			Quantity:  5,
		},
	})
	if err == nil {
		t.Fatal("expected PlaceOrder to fail")
	}

	var stock1 int
	if err := db.QueryRow(ctx, `SELECT stock FROM products WHERE id = $1`,
		product1ID).Scan(&stock1); err != nil {
		t.Fatal(err)
	}
	if stock1 != 10 {
		t.Errorf("got product 1 stock %d, want %d", stock1, 10)
	}

	var stock2 int
	if err := db.QueryRow(ctx, `SELECT stock FROM products WHERE id = $1`,
		product2ID).Scan(&stock2); err != nil {
		t.Fatal(err)
	}
	if stock2 != 1 {
		t.Errorf("got product 2 stock %d, want %d", stock2, 1)
	}

	var orderCount int
	if err := db.QueryRow(ctx, `SELECT COUNT (*)  FROM orders`).Scan(&orderCount); err != nil {
		t.Fatal(err)
	}
	if orderCount != 0 {
		t.Errorf("got %d orders, want %d", orderCount, 0)
	}

	var orderItemsCount int
	if err := db.QueryRow(ctx, `SELECT COUNT (*)FROM order_items`).Scan(&orderItemsCount); err != nil {
		t.Fatal(err)
	}
	if orderItemsCount != 0 {
		t.Errorf("got %d order items, want %d", orderItemsCount, 0)
	}

	var transactionsCount int
	if err := db.QueryRow(ctx, `SELECT COUNT (* )FROM transactions`).Scan(&transactionsCount); err != nil {
		t.Fatal(err)
	}

	if transactionsCount != 0 {
		t.Errorf("got %d transactions, want %d", transactionsCount, 0)
	}
}
