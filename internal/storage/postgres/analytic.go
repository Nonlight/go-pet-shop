package postgres

import (
	"context"
	"fmt"
	"go-pet-shop/internal/models"
)

func (s *Storage) GetUserOrderHistory(ctx context.Context, email string) ([]models.OrderDetail, error) {
	const fn = "storage.postgres.analytic.getUserOrderHistory"

	rows, err := s.db.Query(ctx, `SELECT oi.order_id, o.created_at, p.name, 
oi.quantity, p.price, o.total_price, t.status 
FROM orders o 
JOIN users u ON o.user_id = u.id
JOIN order_items oi ON o.id = oi.order_id
JOIN products p ON oi.product_id = p.id
JOIN transactions t ON t.order_id = o.id
WHERE u.email = $1`, email)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer rows.Close()

	var ordersDetail []models.OrderDetail
	for rows.Next() {
		var orderDetail models.OrderDetail
		err := rows.Scan(&orderDetail.OrderID, &orderDetail.CreatedAt, &orderDetail.ProductName,
			&orderDetail.Quantity, &orderDetail.Price, &orderDetail.TotalPrice, &orderDetail.TransactionStatus)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		ordersDetail = append(ordersDetail, orderDetail)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return ordersDetail, nil
}

func (s *Storage) GetPopularProducts(ctx context.Context) ([]models.PopularProduct, error) {
	const fn = "storage.postgres.analytic.getPopularProducts"

	rows, err := s.db.Query(ctx, `SELECT p.id, p.name, SUM(o.quantity) AS total_sold
FROM products p 
JOIN order_items o ON p.id = o.product_id
GROUP BY p.id, p.name
ORDER BY total_sold DESC`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer rows.Close()

	var popularProducts []models.PopularProduct
	for rows.Next() {
		var product models.PopularProduct
		if err := rows.Scan(&product.ProductID, &product.ProductName, &product.TotalSold); err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		popularProducts = append(popularProducts, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return popularProducts, nil
}
