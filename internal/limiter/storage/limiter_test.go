package storage

import (
	"os"
	"testing"
	"time"
)

func TestNewLimiter(t *testing.T) {
	test := NewLimiter("test")
	if test.db == nil {
		t.Errorf("unable to create new limiter")
	}
}

func TestCreate(t *testing.T) {
	// t.Parallel()

	t.Run("Create: OK", func(t *testing.T) {
		// t.Parallel()

		teardown()
		defer teardown()

		err := NewLimiter("test-limiter").Create("key")
		if err != nil {
			t.Fatalf("expected no error, got %s", err)
		}
	})
}

func TestFind(t *testing.T) {
	// t.Parallel()

	t.Run("Find: OK", func(t *testing.T) {
		// t.Parallel()

		teardown()
		defer teardown()

		limiter := NewLimiter("test-limiter")
		limiter.Create("testkey")
		data, err := limiter.Find()
		if err != nil {
			t.Fatalf("expected no error, got %s", err)
		}

		now := time.Now()
		formatted := now.Format(time.RFC1123)
		actualData := formatted + "testkey" + "\n"

		if data != actualData {
			t.Fatalf("expected result does not match: %s %s", data, actualData)
		}
	})
}

// TODO run test parallel
// https://github.com/mattetti/filebuffer
func teardown() {
	os.Remove("test-limiter")
}
