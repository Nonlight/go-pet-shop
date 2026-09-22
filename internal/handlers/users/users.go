package users

import (
	"context"
	"errors"
	"go-pet-shop/internal/models"
	"go-pet-shop/internal/storage"
	"log/slog"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

//go:generate go run github.com/vektra/mockery/v2 --name=Users
type Users interface {
	GetAllUsers(ctx context.Context) ([]models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	CreateUser(ctx context.Context, user models.User) (int, error)
}

type Handler struct {
	log     *slog.Logger
	storage Users
}

func New(log *slog.Logger, storage Users) *Handler {
	return &Handler{
		log:     log,
		storage: storage,
	}
}

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.users.GetAllUsers"

	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	items, err := h.storage.GetAllUsers(r.Context())
	if err != nil {
		log.Error("failed to fetch users", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to retrieve users",
		})
		return
	}

	log.Info("fetched users",
		slog.String("url", r.URL.String()),
		slog.Int("count", len(items)),
	)

	render.JSON(w, r, items)
}

func (h *Handler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.users.GetUserByEmail"

	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	emailStr := chi.URLParam(r, "email")
	if emailStr == "" {
		log.Error("empty email")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "User email is required",
		})
		return
	}

	item, err := h.storage.GetUserByEmail(r.Context(), emailStr)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			log.Warn("user not found", slog.String("email", emailStr))
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, map[string]string{
				"error":   "Not found",
				"message": "User not found",
			})
			return
		}
		log.Error("failed to get user", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to retrieve user",
		})
		return
	}
	log.Info("Retrieved user successfully",
		slog.String("url", r.URL.String()))

	render.JSON(w, r, item)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.users.CreateUser"

	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	log.Info("Creating new user", slog.String("url", r.URL.String()))

	var user models.User
	if err := render.DecodeJSON(r.Body, &user); err != nil {
		log.Error("failed to decode body", slog.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Invalid JSON payload",
		})
		return
	}

	if user.Name == "" {
		log.Error("user name is empty")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "User name is required",
		})
		return
	}

	if !govalidator.IsEmail(user.Email) {
		log.Error("user email is invalid")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Invalid email format",
		})
		return
	}

	id, err := h.storage.CreateUser(r.Context(), user)
	if err != nil {
		log.Error("failed to create user", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to create user",
		})
		return
	}

	log.Info(
		"User created successfully",
		slog.Int("id", id),
		slog.String("name", user.Name),
		slog.String("url", r.URL.String()))

	user.ID = id
	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, map[string]interface{}{
		"status": "User created successfully",
		"id":     id,
		"user":   user,
	})
}
