package users

import (
	"context"
	"errors"
	"go-pet-shop/internal/handlers/users/mocks"
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

func TestGetAllUsers_Success(t *testing.T) {
	usersMock := mocks.NewUsers(t)
	usersMock.On("GetAllUsers", mock.Anything).Return([]models.User{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), usersMock)
	handler.GetAllUsers(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestGetAllUsers_Error(t *testing.T) {
	usersMock := mocks.NewUsers(t)
	usersMock.On("GetAllUsers", mock.Anything).Return(nil, errors.New("DB error"))

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), usersMock)
	handler.GetAllUsers(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

func TestGetUserByEmail_Success(t *testing.T) {
	usersMock := mocks.NewUsers(t)
	usersMock.On("GetUserByEmail", mock.Anything, "email@mail.ru").Return(models.User{ID: 1, Name: "name", Email: "email@mail.ru"}, nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("email", "email@mail.ru")
	req := httptest.NewRequest(http.MethodGet, "/users/email@mail.ru", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	handler := New(slog.Default(), usersMock)
	handler.GetUserByEmail(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	usersMock := mocks.NewUsers(t)
	usersMock.On("GetUserByEmail", mock.Anything, "email@mail.ru").Return(models.User{}, storage.ErrNotFound)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("email", "email@mail.ru")
	req := httptest.NewRequest(http.MethodGet, "/users/email@mail.ru", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	handler := New(slog.Default(), usersMock)
	handler.GetUserByEmail(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
}

func TestGetUserByEmail_Error(t *testing.T) {
	usersMock := mocks.NewUsers(t)
	usersMock.On("GetUserByEmail", mock.Anything, "email@mail.ru").Return(models.User{}, errors.New("DB error"))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("email", "email@mail.ru")
	req := httptest.NewRequest(http.MethodGet, "/users/email@mail.ru", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	handler := New(slog.Default(), usersMock)
	handler.GetUserByEmail(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

func TestCreateUser_Success(t *testing.T) {
	usersMock := mocks.NewUsers(t)
	usersMock.On("CreateUser", mock.Anything, models.User{Name: "name", Email: "email@mail.ru"}).Return(1, nil)

	body := `{"name":"name", "email":"email@mail.ru"}`

	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), usersMock)
	handler.CreateUser(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}
}

func TestCreateUser_BadRequest(t *testing.T) {
	usersMock := mocks.NewUsers(t)

	body := `{"name":"name", "email":"email@mail.ru"`
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), usersMock)
	handler.CreateUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestCreateUser_Error(t *testing.T) {
	usersMock := mocks.NewUsers(t)
	usersMock.On("CreateUser", mock.Anything, models.User{Name: "name", Email: "email@mail.ru"}).Return(0, errors.New("DB error"))

	body := `{"name":"name", "email":"email@mail.ru"}`

	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), usersMock)
	handler.CreateUser(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}
