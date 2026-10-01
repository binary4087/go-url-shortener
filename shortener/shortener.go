package shortener

import (
	"sync"
	"time"
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

type urlEntry struct {
	longURL   string
	expiresAt time.Time
}

type URLShortener struct {
	mu       sync.RWMutex
	urls     map[string]urlEntry
	reversed map[string]string
	hits     map[string]int
	counter  uint64
}

func New() *URLShortener {
	s := &URLShortener{
		urls:     make(map[string]urlEntry),
		reversed: make(map[string]string),
		hits:     make(map[string]int),
		counter:  100000, // Start at a higher number for consistent length
	}
	go s.startCleanupTimer()
	return s
}

func (s *URLShortener) startCleanupTimer() {
	ticker := time.NewTicker(10 * time.Minute)
	for range ticker.C {
		s.CleanupExpired()
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

	s.urls[short] = urlEntry{
		longURL:   longURL,
		expiresAt: time.Time{}, // No expiration by default
	}
	s.reversed[longURL] = short
	return short
}

func (s *URLShortener) ShortenWithExpiration(longURL string, duration time.Duration) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	if short, exists := s.reversed[longURL]; exists {
		return short
	}

	s.counter++
	short := encodeBase62(s.counter)

	s.urls[short] = urlEntry{
		longURL:   longURL,
		expiresAt: time.Now().Add(duration),
	}
	s.reversed[longURL] = short
	return short
}

func (s *URLShortener) Resolve(shortURL string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, exists := s.urls[shortURL]
	if !exists {
		return "", false
	}

	if !entry.expiresAt.IsZero() && time.Now().After(entry.expiresAt) {
		delete(s.urls, shortURL)
		delete(s.reversed, entry.longURL)
		delete(s.hits, shortURL)
		return "", false
	}

	s.hits[shortURL]++
	return entry.longURL, true
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
	for short, entry := range s.urls {
		stats = append(stats, URLStats{
			ShortURL: short,
			LongURL:  entry.longURL,
			Hits:     s.hits[short],
		})
	}
	return stats
}

func (s *URLShortener) CleanupExpired() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	count := 0
	now := time.Now()
	for short, entry := range s.urls {
		if !entry.expiresAt.IsZero() && now.After(entry.expiresAt) {
			delete(s.reversed, entry.longURL)
			delete(s.hits, short)
			delete(s.urls, short)
			count++
		}
	}
	return count
}