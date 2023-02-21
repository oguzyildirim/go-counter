package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var Create func(key string) error

var Find func() (string, error)

type limiterRepositoryMock struct{}

func (r limiterRepositoryMock) Create(key string) error {
	return Create(key)
}

func (r limiterRepositoryMock) Find() (string, error) {
	return Find()
}

func TestNewRateLimitMiddleware(t *testing.T) {

	t.Run("Middleware: Rate limit", func(t *testing.T) {
		mockRepo := limiterRepositoryMock{}
		Create = func(key string) error {
			return nil
		}
		Find = func() (string, error) {
			now := time.Now()
			formatted := now.Format(time.RFC1123)
			request := formatted + "1.1.1.1" + "\n"
			var data string
			for i := 0; i < 7; i++ {
				data = data + request
			}
			return data, nil
		}

		middleware := NewRateLimitMiddleware(mockRepo)

		doWork := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "hello world")
		})

		server := httptest.NewServer(middleware.Handle(doWork))

		client := server.Client()

		req, _ := http.NewRequest("GET", server.URL, nil)
		req.Header.Set("X-Real-IP", "1.1.1.1")

		resp, _ := client.Do(req)
		if got, want := resp.StatusCode, http.StatusTooManyRequests; got != want {
			t.Errorf("expected %d to be %d", got, want)
		}
	})

	t.Run("Middleware: no limit", func(t *testing.T) {
		mockRepo := limiterRepositoryMock{}
		Create = func(key string) error {
			return nil
		}
		Find = func() (string, error) {
			now := time.Now()
			formatted := now.Format(time.RFC1123)
			request := formatted + "1.2.1.1" + "\n"
			var data string
			for i := 0; i < 2; i++ {
				data = data + request
			}
			return data, nil
		}

		middleware := NewRateLimitMiddleware(mockRepo)

		doWork := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "hello world")
		})

		server := httptest.NewServer(middleware.Handle(doWork))

		client := server.Client()

		req, _ := http.NewRequest("GET", server.URL, nil)
		req.Header.Set("X-Real-IP", "1.2.1.1")

		resp, _ := client.Do(req)
		if got, want := resp.StatusCode, http.StatusOK; got != want {
			t.Errorf("expected %d to be %d", got, want)
		}
	})
}
