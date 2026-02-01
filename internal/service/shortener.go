package service

import (
	"crypto/rand"
	"encoding/base64"
	"sync"
)

type ShortenerService struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewShortenerService() *ShortenerService {
	return &ShortenerService{
		data: make(map[string]string),
	}
}

func (s *ShortenerService) Shorten(longURL string) string {
	code := generateCode(6)

	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[code] = longURL
	return code
}

func (s *ShortenerService) Resolve(code string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	url, ok := s.data[code]
	return url, ok
}

func generateCode(n int) string {
	b := make([]byte, n)
	rand.Read(b)

	// URL safe string
	return base64.RawURLEncoding.EncodeToString(b)[:n]
}
