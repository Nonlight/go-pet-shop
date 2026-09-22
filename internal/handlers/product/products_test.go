package product

import (
	"context"
	"errors"
	"go-pet-shop/internal/handlers/product/mocks"
	"go-pet-shop/internal/models"
	"go-pet-shop/internal/storage"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/mock"
)

func TestGetAllProducts_Success(t *testing.T) {
	productMock := mocks.NewProducts(t)
	productMock.On("GetAllProducts", mock.Anything).Return([]models.Product{{ID: 1, Name: "Cat Food"}}, nil)

	// Создаем HTTP-запрос GET /products
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	// Создаем хендлер с мок-хранилищем
	handler := New(slog.Default(), productMock)

	// Вызываем метод GetAllProducts, который является http.HandlerFunc
	handler.GetAllProducts(w, req)

	// Проверяем HTTP-код
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}
func TestGetAllProducts_Error(t *testing.T) {
	productMock := mocks.NewProducts(t)
	productMock.On("GetAllProducts", mock.Anything).Return(nil, errors.New("DB error"))

	// Создаем запрос
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), productMock)
	handler.GetAllProducts(w, req)

	// Ожидаем HTTP 500
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// =======================
// Get Product by ID
// =======================

func TestGetProductByID_Success(t *testing.T) {
	productMock := mocks.NewProducts(t)
	productMock.On("GetProductByID", mock.Anything, 1).Return(models.Product{ID: 1, Name: "Cat food"}, nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req := httptest.NewRequest(http.MethodGet, "/products/1", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	handler := New(slog.Default(), productMock)
	handler.GetProductByID(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

}

func TestGetProductByID_NotFound(t *testing.T) {
	productMock := mocks.NewProducts(t)
	productMock.On("GetProductByID", mock.Anything, 1).Return(models.Product{}, storage.ErrNotFound)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req := httptest.NewRequest(http.MethodGet, "/products/1", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	handler := New(slog.Default(), productMock)
	handler.GetProductByID(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
}

func TestGetProductByID_Error(t *testing.T) {
	productMock := mocks.NewProducts(t)
	productMock.On("GetProductByID", mock.Anything, 1).Return(models.Product{}, errors.New("DB error"))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req := httptest.NewRequest(http.MethodGet, "/products/1", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	handler := New(slog.Default(), productMock)
	handler.GetProductByID(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// =======================
// Create Product
// =======================

func TestCreateProduct_Success(t *testing.T) {
	productMock := mocks.NewProducts(t)
	productMock.On("CreateProduct", mock.Anything, models.Product{
		Name:  "Cat Food",
		Price: 12.5,
		Stock: 10},
	).Return(1, nil)

	body := `{"name":"Cat Food","price":12.5,"stock":10}`

	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), productMock)

	handler.CreateProduct(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}

}

func TestCreateProduct_BadRequest(t *testing.T) {
	productMock := mocks.NewProducts(t)

	body := `{"name":"Cat Food","price":12.5,"stock":10`
	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), productMock)

	handler.CreateProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestCreateProduct_Fail(t *testing.T) {
	productMock := mocks.NewProducts(t)
	productMock.On("CreateProduct", mock.Anything, mock.Anything).Return(0, errors.New("DB error"))

	body := `{"name":"Cat Food","price":12.5,"stock":10}`
	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), productMock)
	handler.CreateProduct(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// =======================
// Update Product
// =======================

func TestUpdateProduct_Success(t *testing.T) {
	productMock := mocks.NewProducts(t)
	productMock.On("UpdateProduct", mock.Anything, models.Product{
		ID:    1,
		Name:  "Cat Food",
		Price: 15.5,
		Stock: 20},
	).Return(nil)

	body := `{"name":"Cat Food","price":15.5,"stock":20}`
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")

	req := httptest.NewRequest(http.MethodPut, "/products/1", strings.NewReader(body))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	handler := New(slog.Default(), productMock)
	handler.UpdateProduct(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestUpdateProduct_BadRequest(t *testing.T) {
	productMock := mocks.NewProducts(t)

	body := `{"name":"Cat Food","price":15.5,"stock":20`
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")

	req := httptest.NewRequest(http.MethodPut, "/products/1", strings.NewReader(body))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), productMock)
	handler.UpdateProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestUpdateProduct_Fail(t *testing.T) {
	productMock := mocks.NewProducts(t)
	productMock.On("UpdateProduct", mock.Anything, mock.Anything).Return(errors.New("DB error"))

	body := `{"name":"Cat Food","price":15.5,"stock":20}`
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")

	req := httptest.NewRequest(http.MethodPut, "/products/1", strings.NewReader(body))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	handler := New(slog.Default(), productMock)
	handler.UpdateProduct(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// =======================
// Delete Product
// =======================

func TestDeleteProduct_Success(t *testing.T) {
	productMock := mocks.NewProducts(t)
	productMock.On("DeleteProduct", mock.Anything, 1).Return(nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	handler := New(slog.Default(), productMock)

	handler.DeleteProduct(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestDeleteProduct_BadRequest(t *testing.T) {
	productMock := mocks.NewProducts(t)

	rctx := chi.NewRouteContext()

	req := httptest.NewRequest(http.MethodDelete, "/products/", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	handler := New(slog.Default(), productMock)
	handler.DeleteProduct(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestDeleteProduct_Fail(t *testing.T) {
	productMock := mocks.NewProducts(t)
	productMock.On("DeleteProduct", mock.Anything, mock.Anything).Return(errors.New("DB error"))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	handler := New(slog.Default(), productMock)

	handler.DeleteProduct(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}
