package analytics

import (
	"context"
	"go-pet-shop/internal/models"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

//go:generate go run github.com/vektra/mockery/v2 --name=Analytics
type Analytics interface {
	GetUserOrderHistory(ctx context.Context, email string) ([]models.OrderDetail, error)
	GetPopularProducts(ctx context.Context) ([]models.PopularProduct, error)
}

type Handler struct {
	log     *slog.Logger
	storage Analytics
}

func New(log *slog.Logger, storage Analytics) *Handler {
	return &Handler{
		log:     log,
		storage: storage,
	}
}

func (h *Handler) GetUserOrderHistory(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.analytics.GetUserOrderHistory"

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

	orderDetails, err := h.storage.GetUserOrderHistory(r.Context(), email)
	if err != nil {
		log.Error("failed to get orders", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to retrieve orders",
		})
		return
	}

	log.Info("Retrieved order history successfully",
		slog.String("url", r.URL.String()))

	render.JSON(w, r, orderDetails)
}

func (h *Handler) GetPopularProducts(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.analytics.GetPopularProducts"

	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())))

	popularProducts, err := h.storage.GetPopularProducts(r.Context())
	if err != nil {
		log.Error("failed to get products", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to retrieve products",
		})
		return
	}
	log.Info("Retrieved products successfully",
		slog.String("url", r.URL.String()))

	render.JSON(w, r, popularProducts)
}
