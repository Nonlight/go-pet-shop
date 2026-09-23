package postgres

import (
	"context"
	"errors"
	"fmt"
	"go-pet-shop/internal/models"
	"go-pet-shop/internal/storage"

	"github.com/jackc/pgx/v5"
)

func (s *Storage) CreateOrder(ctx context.Context, order models.Order) (int, error) {
	const fn = "storage.postgres.order.CreateOrder"

	var id int
	err := s.db.QueryRow(ctx, `INSERT INTO orders (user_id, total_price) VALUES ($1, $2) RETURNING id`,
		order.UserID, order.TotalPrice).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}
	return id, nil
}

func (s *Storage) AddOrderItem(ctx context.Context, orderItem models.OrderItem) error {
	const fn = "storage.postgres.order.AddOrderItem"

	_, err := s.db.Exec(ctx, `INSERT INTO order_items(order_id, product_id, quantity) VALUES ($1, $2, $3)`,
		orderItem.OrderID, orderItem.ProductID, orderItem.Quantity)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	return nil
}

func (s *Storage) GetOrderByID(ctx context.Context, id int) (models.Order, error) {
	const fn = "storage.postgres.order.GetOrderByID"

	var order models.Order
	if err := s.db.QueryRow(ctx, `SELECT id, user_id, total_price, created_at FROM orders WHERE id = $1`, id).Scan(
		&order.ID,
		&order.UserID,
		&order.TotalPrice,
		&order.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Order{}, storage.ErrNotFound
		}
		return models.Order{}, fmt.Errorf("%s: %w", fn, err)
	}

	return order, nil
}

func (s *Storage) GetOrdersByUserEmail(ctx context.Context, email string) ([]models.Order, error) {
	const fn = "storage.postgres.order.GetOrdersByUserEmail"

	rows, err := s.db.Query(ctx, `SELECT o.id, o.user_id, o.total_price, o.created_at FROM users u JOIN orders o on o.user_id = u.id WHERE u.email = $1`, email)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.ID, &order.UserID, &order.TotalPrice, &order.CreatedAt); err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return orders, nil
}

func (s *Storage) GetOrderItemsByOrderID(ctx context.Context, orderID int) ([]models.OrderItem, error) {
	const fn = "storage.postgres.order.GetOrderItemsByOrderID"

	rows, err := s.db.Query(ctx, `SELECT product_id, quantity FROM order_items WHERE order_id = $1`, orderID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer rows.Close()

	var orderItems []models.OrderItem
	for rows.Next() {
		var orderItem models.OrderItem
		if err := rows.Scan(&orderItem.ProductID, &orderItem.Quantity); err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		orderItems = append(orderItems, orderItem)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	return orderItems, nil
}
