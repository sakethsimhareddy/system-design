package service

import (
	"math/rand"
	"time"

	"url-shortener/internal/model"
	"url-shortener/internal/repository"
)

type Service interface {
	ShortenURL(originalURL string) (string, error)
	GetOriginalURL(shortCode string) (string, error)
}

type urlService struct {
	repo repository.Repository
}

func NewURLService(repo repository.Repository) Service {
	return &urlService{
		repo: repo,
	}
}

func (s *urlService) ShortenURL(originalURL string) (string, error) {
	shortCode := generateShortCode(6)
	shortURL := model.ShortURL{
		ID:          shortCode, // Using shortCode as ID for simplicity
		OriginalURL: originalURL,
		ShortCode:   shortCode,
		CreatedAt:   time.Now(),
	}

	if _,err := s.repo.Save(shortURL); err != nil {
		return "", err
	}

	return shortCode, nil
}

func (s *urlService) GetOriginalURL(shortCode string) (string, error) {

	shortURL, err := s.repo.Get(shortCode)
	if err != nil {
		return "", err
	}
	return shortURL, nil
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateShortCode(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
