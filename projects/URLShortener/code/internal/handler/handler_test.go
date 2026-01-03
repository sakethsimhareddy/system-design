package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockService is a mock implementation of the Service interface
type MockService struct {
	mock.Mock
}

func (m *MockService) ShortenURL(originalURL string) (string, error) {
	args := m.Called(originalURL)
	return args.String(0), args.Error(1)
}

func (m *MockService) GetOriginalURL(shortCode string) (string, error) {
	args := m.Called(shortCode)
	return args.String(0), args.Error(1)
}

func TestShortenURL(t *testing.T) {
	mockSvc := new(MockService)
	handler := NewHandler(mockSvc)

	originalURL := "http://example.com"
	shortCode := "abcdef"
	mockSvc.On("ShortenURL", originalURL).Return(shortCode, nil)

	reqBody := `{"url": "http://example.com"}`
	req, _ := http.NewRequest("POST", "/shorten", bytes.NewBufferString(reqBody))
	rr := httptest.NewRecorder()

	handler.ShortenURL(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp ShortenResponse
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, shortCode, resp.ShortCode)

	mockSvc.AssertExpectations(t)
}

func TestShortenURL_InvalidBody(t *testing.T) {
	mockSvc := new(MockService)
	handler := NewHandler(mockSvc)

	reqBody := `invalid json`
	req, _ := http.NewRequest("POST", "/shorten", bytes.NewBufferString(reqBody))
	rr := httptest.NewRecorder()

	handler.ShortenURL(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestShortenURL_ServiceError(t *testing.T) {
	mockSvc := new(MockService)
	handler := NewHandler(mockSvc)

	originalURL := "http://example.com"
	mockSvc.On("ShortenURL", originalURL).Return("", errors.New("internal error"))

	reqBody := `{"url": "http://example.com"}`
	req, _ := http.NewRequest("POST", "/shorten", bytes.NewBufferString(reqBody))
	rr := httptest.NewRecorder()

	handler.ShortenURL(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	mockSvc.AssertExpectations(t)
}

func TestRedirect(t *testing.T) {
	mockSvc := new(MockService)
	handler := NewHandler(mockSvc)

	shortCode := "abcdef"
	originalURL := "http://example.com"
	mockSvc.On("GetOriginalURL", shortCode).Return(originalURL, nil)

	req, _ := http.NewRequest("GET", "/"+shortCode, nil)
	rr := httptest.NewRecorder()

	handler.Redirect(rr, req)

	assert.Equal(t, http.StatusFound, rr.Code)
	assert.Equal(t, originalURL, rr.Header().Get("Location"))

	mockSvc.AssertExpectations(t)
}

func TestRedirect_NotFound(t *testing.T) {
	mockSvc := new(MockService)
	handler := NewHandler(mockSvc)

	shortCode := "unknown"
	mockSvc.On("GetOriginalURL", shortCode).Return("", errors.New("not found"))

	req, _ := http.NewRequest("GET", "/"+shortCode, nil)
	rr := httptest.NewRecorder()

	handler.Redirect(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	mockSvc.AssertExpectations(t)
}
