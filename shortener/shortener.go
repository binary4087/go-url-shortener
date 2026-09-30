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
	hits     map[string]int
}

func New() *URLShortener {
	return &URLShortener{
		urls:     make(map[string]string),
		reversed: make(map[string]string),
		hits:     make(map[string]int),
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
	s.mu.Lock()
	defer s.mu.Unlock()

	long, exists := s.urls[shortURL]
	if exists {
		s.hits[shortURL]++
	}
	return long, exists
}

func (s *URLShortener) GetHits(shortURL string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.hits[shortURL]
}