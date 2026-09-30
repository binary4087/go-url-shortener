package shortener

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"sync"
)

type URLShortener struct {
	mu       sync.RWMutex
	urls     map[string]string
	reversed map[string]string
}

func New() *URLShortener {
	return &URLShortener{
		urls:     make(map[string]string),
		reversed: make(map[string]string),
	}
}

func (s *URLShortener) Shorten(longURL string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	if short, exists := s.reversed[longURL]; exists {
		return short
	}

	hash := sha256.Sum256([]byte(longURL))
	short := base64.URLEncoding.EncodeToString(hash[:])[:8]

	// Handle potential collisions
	for {
		if _, exists := s.urls[short]; !exists {
			break
		}
		short = short + fmt.Sprintf("%x", hash[8])[:1]
	}

	s.urls[short] = longURL
	s.reversed[longURL] = short
	return short
}

func (s *URLShortener) Resolve(shortURL string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	long, exists := s.urls[shortURL]
	return long, exists
}