package analytics

import (
	"encoding/json"
	"errors"
	"go-pet-shop/internal/handlers/analytics/mocks"
	"go-pet-shop/internal/models"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetUserOrderHistory_Success(t *testing.T) {
	ordersMock := mocks.NewAnalytics(t)
	createdAt := time.Now()
	ordersMock.On("GetUserOrderHistory", mock.Anything, "email@mail.ru").Return([]models.OrderDetail{{OrderID: 1, CreatedAt: createdAt, ProductName: "Cat Food",
		Quantity: 2, Price: 200, TotalPrice: 400, TransactionStatus: "success"}}, nil)

	req := httptest.NewRequest(http.MethodGet, "/users/history?email=email@mail.ru", nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), ordersMock)
	handler.GetUserOrderHistory(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var got []models.OrderDetail
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	expected := []models.OrderDetail{
		{
			OrderID:           1,
			CreatedAt:         createdAt,
			ProductName:       "Cat Food",
			Quantity:          2,
			Price:             200,
			TotalPrice:        400,
			TransactionStatus: "success",
		},
	}

	expected[0].CreatedAt = expected[0].CreatedAt.Round(0)
	got[0].CreatedAt = got[0].CreatedAt.Round(0)

	assert.Equal(t, expected, got)
}

func TestGetUserOrderHistory_BadRequest(t *testing.T) {
	ordersMock := mocks.NewAnalytics(t)

	req := httptest.NewRequest(http.MethodGet, "/users/history", nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), ordersMock)
	handler.GetUserOrderHistory(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestGetUserOrderHistory_Error(t *testing.T) {
	ordersMock := mocks.NewAnalytics(t)
	ordersMock.On("GetUserOrderHistory", mock.Anything, "email@mail.ru").Return([]models.OrderDetail{}, errors.New("DB error"))

	req := httptest.NewRequest(http.MethodGet, "/users/history?email=email@mail.ru", nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), ordersMock)
	handler.GetUserOrderHistory(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

func TestGetPopularProducts_Success(t *testing.T) {
	ordersMock := mocks.NewAnalytics(t)
	ordersMock.On("GetPopularProducts", mock.Anything).Return([]models.PopularProduct{{ProductID: 1, ProductName: "Cat Food", TotalSold: 400}}, nil)

	req := httptest.NewRequest(http.MethodGet, "/products/popular", nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), ordersMock)
	handler.GetPopularProducts(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var got []models.PopularProduct
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	expected := []models.PopularProduct{
		{ProductID: 1, ProductName: "Cat Food", TotalSold: 400},
	}

	assert.Equal(t, expected, got)
}

func TestGetPopularProducts_Error(t *testing.T) {
	ordersMock := mocks.NewAnalytics(t)
	ordersMock.On("GetPopularProducts", mock.Anything).Return(nil, errors.New("DB error"))

	req := httptest.NewRequest(http.MethodGet, "/products/popular", nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), ordersMock)
	handler.GetPopularProducts(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}
