package middleware

import (
	"bufio"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Oguzyildirim/go-counter/internal"
)

const (
	HeaderRateLimitLimit     = "X-RateLimit-Limit"
	HeaderRateLimitRemaining = "X-RateLimit-Remaining"
	HeaderRateLimitReset     = "X-RateLimit-Reset"

	HeaderRetryAfter = "Retry-After"
)

const (
	interval = 20 * time.Second
	limit    = 5
)

// LimiterRepository defines the datastore handling persisting Limiter records
type LimiterRepository interface {
	Create(key string) error
	Find() (string, error)
}

// RateLimitMiddleware is a handler/mux that can wrap other middlware to implement HTTP
type RateLimitMiddleware struct {
	db LimiterRepository
}

// NewRateLimitMiddleware creates a new middleware suitable for use as an HTTP handler.
func NewRateLimitMiddleware(db LimiterRepository) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		db: db,
	}
}

// Handle returns the HTTP handler as a middleware.
func (m *RateLimitMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key, err := resolveIP(r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		data, err := m.db.Find()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		resetTime, remaining, err := compose(data, key)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set(HeaderRateLimitLimit, strconv.FormatUint(limit, 10))
		w.Header().Set(HeaderRateLimitRemaining, strconv.Itoa(remaining))
		w.Header().Set(HeaderRateLimitReset, resetTime)

		if remaining < 1 {
			w.Header().Set(HeaderRetryAfter, resetTime)
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		m.db.Create(key)

		next.ServeHTTP(w, r)
	})
}

func compose(data string, key string) (string, int, error) {
	var count int
	var times []time.Time
	var keys []string
	var startTime time.Time

	scanner := bufio.NewScanner(strings.NewReader(data))
	for scanner.Scan() {
		data := scanner.Text()
		t, err := time.Parse(time.RFC1123, data[:29])
		if err != nil {
			return "", 0, fmt.Errorf("time.Parse: %w", err)
		}
		times = append(times, t)
		keys = append(keys, data[29:])
	}
	if err := scanner.Err(); err != nil {
		return "", 0, fmt.Errorf("NewScanner failed %w", err)
	}

	now := time.Now()
	window := now.Add(-interval)
	for i := len(times) - 1; i >= 0; i-- {
		if internal.InTimeSpan(window, now, times[i]) && key == keys[i] {
			count = count + 1
			startTime = times[i]
		}
	}
	resetTime := startTime.Add(interval)
	formatted := resetTime.Format(time.RFC1123)
	remaining := limit - count
	return formatted, remaining, nil
}
