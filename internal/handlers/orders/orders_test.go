package orders

import (
	"context"
	"errors"
	"go-pet-shop/internal/handlers/orders/mocks"
	"go-pet-shop/internal/models"
	"go-pet-shop/internal/storage"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/mock"
)

func TestCreateOrder_Success(t *testing.T) {
	ordersMock := mocks.NewOrders(t)
	ordersMock.On("CreateOrder", mock.Anything, models.Order{UserID: 1}).Return(1, nil)

	body := `{"user_id": 1}`

	req := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), ordersMock)
	handler.CreateOrder(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}
}

func TestCreateOrder_BadRequest(t *testing.T) {
	ordersMock := mocks.NewOrders(t)

	body := `{"user_id": -4}`

	req := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), ordersMock)
	handler.CreateOrder(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestCreateOrder_Error(t *testing.T) {
	ordersMock := mocks.NewOrders(t)
	ordersMock.On("CreateOrder", mock.Anything, models.Order{UserID: 1}).Return(0, errors.New("DB error"))

	body := `{"user_id": 1}`

	req := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), ordersMock)
	handler.CreateOrder(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

func TestAddOrderItem_Success(t *testing.T) {
	ordersMock := mocks.NewOrders(t)
	ordersMock.On("AddOrderItem", mock.Anything, models.OrderItem{OrderID: 1, ProductID: 2, Quantity: 3}).Return(nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")

	body := `{"product_id": 2, "quantity": 3}`
	req := httptest.NewRequest(http.MethodPost, "/orders/1/items", strings.NewReader(body))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	handler := New(slog.Default(), ordersMock)
	handler.AddOrderItem(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}
}

func TestAddOrderItem_BadRequest(t *testing.T) {
	ordersMock := mocks.NewOrders(t)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")

	body := `{"product_id": 2, "quantity": -3}`
	req := httptest.NewRequest(http.MethodPost, "/orders/1/items", strings.NewReader(body))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	handler := New(slog.Default(), ordersMock)
	handler.AddOrderItem(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestAddOrderItem_Error(t *testing.T) {
	ordersMock := mocks.NewOrders(t)
	ordersMock.On("AddOrderItem", mock.Anything, models.OrderItem{OrderID: 1, ProductID: 2, Quantity: 3}).Return(errors.New("DB error"))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")

	body := `{"product_id": 2, "quantity": 3}`
	req := httptest.NewRequest(http.MethodPost, "/orders/1/items", strings.NewReader(body))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	handler := New(slog.Default(), ordersMock)
	handler.AddOrderItem(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

func TestGetOrderByID_Success(t *testing.T) {
	ordersMock := mocks.NewOrders(t)
	ordersMock.On("GetOrderByID", mock.Anything, 1).Return(models.Order{UserID: 1, TotalPrice: 2500, CreatedAt: time.Now()}, nil)
	ordersMock.On("GetOrderItemsByOrderID", mock.Anything, 1).Return([]models.OrderItem{}, nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")

	req := httptest.NewRequest(http.MethodGet, "/orders/1", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), ordersMock)
	handler.GetOrderByID(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestGetOrderByID_BadRequest(t *testing.T) {
	ordersMock := mocks.NewOrders(t)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "-1")

	req := httptest.NewRequest(http.MethodGet, "/orders/1", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), ordersMock)
	handler.GetOrderByID(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestGetOrderByID_NotFound(t *testing.T) {
	ordersMock := mocks.NewOrders(t)
	ordersMock.On("GetOrderByID", mock.Anything, 1).Return(models.Order{}, storage.ErrNotFound)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")

	req := httptest.NewRequest(http.MethodGet, "/orders/1", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), ordersMock)
	handler.GetOrderByID(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
}

func TestGetOrderByID_Error(t *testing.T) {
	ordersMock := mocks.NewOrders(t)
	ordersMock.On("GetOrderByID", mock.Anything, 1).Return(models.Order{}, errors.New("DB error"))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")

	req := httptest.NewRequest(http.MethodGet, "/orders/1", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), ordersMock)
	handler.GetOrderByID(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

func TestGetOrdersByUserEmail_Success(t *testing.T) {
	ordersMock := mocks.NewOrders(t)
	ordersMock.On("GetOrdersByUserEmail", mock.Anything, "email@mail.ru").Return([]models.Order{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/users/orders?email=email@mail.ru", nil)

	w := httptest.NewRecorder()
	handler := New(slog.Default(), ordersMock)
	handler.GetOrdersByUserEmail(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestGetOrdersByUserEmail_BadRequest(t *testing.T) {
	ordersMock := mocks.NewOrders(t)

	req := httptest.NewRequest(http.MethodGet, "/users/orders?email=", nil)

	w := httptest.NewRecorder()
	handler := New(slog.Default(), ordersMock)
	handler.GetOrdersByUserEmail(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestGetOrdersByUserEmail_Error(t *testing.T) {
	ordersMock := mocks.NewOrders(t)
	ordersMock.On("GetOrdersByUserEmail", mock.Anything, "email@mail.ru").Return(nil, errors.New("DB error"))

	req := httptest.NewRequest(http.MethodGet, "/users/orders?email=email@mail.ru", nil)

	w := httptest.NewRecorder()
	handler := New(slog.Default(), ordersMock)
	handler.GetOrdersByUserEmail(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}
