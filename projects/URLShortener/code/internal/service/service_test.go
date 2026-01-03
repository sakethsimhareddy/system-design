package service

import (
	"errors"
	"testing"

	"url-shortener/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of the Repository interface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Save(shortURL model.ShortURL) (model.ShortURL, error) {
	args := m.Called(shortURL)
	return args.Get(0).(model.ShortURL), args.Error(1)
}

func (m *MockRepository) Get(shortCode string) (string, error) {
	args := m.Called(shortCode)
	return args.String(0), args.Error(1)
}

func TestShortenURL(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewURLService(mockRepo)

	originalURL := "http://example.com"

	// Expect Save to be called. We use mock.AnythingOfType for the ShortURL argument
	// because ID and CreatedAt are generated inside the service.
	mockRepo.On("Save", mock.AnythingOfType("model.ShortURL")).Return(model.ShortURL{}, nil)

	shortCode, err := service.ShortenURL(originalURL)

	assert.NoError(t, err)
	assert.NotEmpty(t, shortCode)
	assert.Len(t, shortCode, 6)

	mockRepo.AssertExpectations(t)
}

func TestShortenURL_Error(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewURLService(mockRepo)

	originalURL := "http://example.com"
	expectedError := errors.New("database error")

	mockRepo.On("Save", mock.AnythingOfType("model.ShortURL")).Return(model.ShortURL{}, expectedError)

	shortCode, err := service.ShortenURL(originalURL)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Empty(t, shortCode)

	mockRepo.AssertExpectations(t)
}

func TestGetOriginalURL(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewURLService(mockRepo)

	shortCode := "abcdef"
	originalURL := "http://example.com"

	mockRepo.On("Get", shortCode).Return(originalURL, nil)

	result, err := service.GetOriginalURL(shortCode)

	assert.NoError(t, err)
	assert.Equal(t, originalURL, result)

	mockRepo.AssertExpectations(t)
}

func TestGetOriginalURL_Error(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewURLService(mockRepo)

	shortCode := "invalid"
	expectedError := errors.New("not found")

	mockRepo.On("Get", shortCode).Return("", expectedError)

	result, err := service.GetOriginalURL(shortCode)

	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.Equal(t, expectedError, err)

	mockRepo.AssertExpectations(t)
}
