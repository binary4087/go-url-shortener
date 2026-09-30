package shortener

import "testing"

func TestShortener(t *testing.T) {
	s := New()
	long := "https://www.example.com/some/very/long/path"
	short := s.Shorten(long)

	if len(short) == 0 {
		t.Error("Expected non-empty short URL")
	}

	resolved, exists := s.Resolve(short)
	if !exists || resolved != long {
		t.Errorf("Expected %s, got %s", long, resolved)
	}

	short2 := s.Shorten(long)
	if short != short2 {
		t.Errorf("Expected consistent mapping, got %s and %s", short, short2)
	}
}