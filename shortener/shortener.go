package shortener

import (
	"sync"
)

const base62Chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func encodeBase62(n uint64) string {
	if n == 0 {
		return string(base62Chars[0])
	}
	var res []byte
	for n > 0 {
		res = append(res, base62Chars[n%62])
		n /= 62
	}
	// Reverse the slice
	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}
	return string(res)
}

type URLStats struct {
	ShortURL string
	LongURL  string
	Hits     int
}

type URLShortener struct {
	mu       sync.RWMutex
	urls     map[string]string
	reversed map[string]string
	hits     map[string]int
	counter  uint64
}

func New() *URLShortener {
	return &URLShortener{
		urls:     make(map[string]string),
		reversed: make(map[string]string),
		hits:     make(map[string]int),
		counter:  100000, // Start at a higher number for consistent length
	}
}

func (s *URLShortener) Shorten(longURL string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	if short, exists := s.reversed[longURL]; exists {
		return short
	}

	s.counter++
	short := encodeBase62(s.counter)

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

func (s *URLShortener) ListAll() []URLStats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := make([]URLStats, 0, len(s.urls))
	for short, long := range s.urls {
		stats = append(stats, URLStats{
			ShortURL: short,
			LongURL:  long,
			Hits:     s.hits[short],
		})
	}
	return stats
}