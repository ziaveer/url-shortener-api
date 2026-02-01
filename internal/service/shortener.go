package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	"github.com/ziaulhaq/url-shortener/internal/repository"
)

type ShortenerService struct {
	repo repository.URLRepository
}

func NewShortenerService(repo repository.URLRepository) *ShortenerService {
	return &ShortenerService{
		repo: repo,
	}
}

func (s *ShortenerService) Shorten(ctx context.Context, longURL string) (string, error) {
	code := generateCode(6)

	err := s.repo.Save(ctx, code, longURL)
	if err != nil {
		return "", err
	}

	return code, nil
}

func (s *ShortenerService) Resolve(ctx context.Context, code string) (string, bool, error) {
	return s.repo.FindByCode(ctx, code)
}

func generateCode(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)[:n]
}
