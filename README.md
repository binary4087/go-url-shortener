# Go URL Shortener

A simple, fast, in-memory URL shortener implemented in Go.

## Usage

1. Run the server:
   ```bash
   go run main.go
   ```

2. Shorten a URL:
   `http://localhost:8080/shorten?url=https://google.com`

3. Use the short link:
   `http://localhost:8080/r/<hash>`