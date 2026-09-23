package orders

import (
	"context"
	"errors"
	"go-pet-shop/internal/models"
	"go-pet-shop/internal/storage"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

//go:generate go run github.com/vektra/mockery/v2 --name=Orders
type Orders interface {
	CreateOrder(ctx context.Context, order models.Order) (int, error)
	AddOrderItem(ctx context.Context, orderItem models.OrderItem) error
	GetOrderByID(ctx context.Context, id int) (models.Order, error)
	GetOrdersByUserEmail(ctx context.Context, email string) ([]models.Order, error)
	GetOrderItemsByOrderID(ctx context.Context, orderID int) ([]models.OrderItem, error)
}
type Handler struct {
	log     *slog.Logger
	storage Orders
}

func New(log *slog.Logger, storage Orders) *Handler {
	return &Handler{
		log:     log,
		storage: storage,
	}
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.order.CreateOrder"

	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	log.Info("Creating new order", slog.String("url", r.URL.String()))

	var order models.Order
	if err := render.DecodeJSON(r.Body, &order); err != nil {
		log.Error("failed to decode body", slog.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Invalid JSON payload",
		})
		return
	}

	if order.UserID <= 0 {
		log.Error("user id is invalid")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Invalid id format",
		})
		return
	}

	order.TotalPrice = 0

	id, err := h.storage.CreateOrder(r.Context(), order)
	if err != nil {
		log.Error("failed to create order", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to create order",
		})
		return
	}

	log.Info(
		"Order created successfully",
		slog.Int("id", id),
		slog.Int("user_id", order.UserID),
		slog.String("url", r.URL.String()))

	order.ID = id
	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, map[string]interface{}{
		"status": "Order created successfully",
		"id":     id,
		"order":  order,
	})
}
func (h *Handler) AddOrderItem(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.order.AddOrderItem"

	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		log.Error("empty id parameter")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Order id is required",
		})
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error("invalid id parameter")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Invalid id parameter",
		})
		return
	}

	if id <= 0 {
		log.Error("invalid id parameter")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Invalid id parameter",
		})
		return
	}

	var orderItem models.OrderItem
	if err := render.DecodeJSON(r.Body, &orderItem); err != nil {
		log.Error("failed to decode body", slog.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Invalid JSON payload",
		})
		return
	}

	if orderItem.ProductID <= 0 {
		log.Error("product id is invalid")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Invalid JSON payload",
		})
		return
	}

	if orderItem.Quantity <= 0 {
		log.Error("quantity is invalid")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Invalid JSON payload",
		})
		return
	}

	orderItem.OrderID = id

	if err := h.storage.AddOrderItem(r.Context(), orderItem); err != nil {
		log.Error("failed to add order item", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to create order item",
		})
		return
	}

	log.Info("Order item created successfully",
		slog.String("url", r.URL.String()))

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, map[string]interface{}{
		"status":     "Order item created successfully",
		"order_id":   orderItem.OrderID,
		"product_id": orderItem.ProductID,
		"quantity":   orderItem.Quantity,
	})
}

func (h *Handler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.order.GetOrderByID"

	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		log.Error("empty id parameter")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Id is required",
		})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error("invalid id parameter")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Invalid id parameter",
		})
		return
	}

	if id <= 0 {
		log.Error("invalid id parameter")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Invalid id parameter",
		})
		return
	}

	order, err := h.storage.GetOrderByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			log.Warn("order not found", slog.Int("id", id))
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, map[string]string{
				"error":   "Not found",
				"message": "Order not found",
			})
			return
		}

		log.Error("failed to get order", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to retrieve order",
		})
		return
	}

	orderItem, err := h.storage.GetOrderItemsByOrderID(r.Context(), id)
	if err != nil {
		log.Error("failed to get order", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to retrieve order",
		})
		return
	}

	type orderResponse struct {
		Order      models.Order
		OrderItems []models.OrderItem
	}

	log.Info("Order retrieved successfully",
		slog.String("url", r.URL.String()))

	render.JSON(w, r, orderResponse{Order: order, OrderItems: orderItem})
}

func (h *Handler) GetOrdersByUserEmail(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.order.GetOrdersByUserEmail"

	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	email := r.URL.Query().Get("email")
	if email == "" {
		log.Error("email is required")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Email is required",
		})
		return
	}

	orders, err := h.storage.GetOrdersByUserEmail(r.Context(), email)
	if err != nil {
		log.Error("failed to get orders", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to retrieve orders",
		})
		return
	}

	log.Info("Retrieved orders successfully",
		slog.String("url", r.URL.String()))

	render.JSON(w, r, orders)
}
