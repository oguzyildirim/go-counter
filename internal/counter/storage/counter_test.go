package storage

import (
	"os"
	"testing"
	"time"
)

func TestNewCounter(t *testing.T) {
	test := NewCounter("test")
	if test.db == nil {
		t.Errorf("unable to create new counter")
	}
}

func TestCreate(t *testing.T) {
	// t.Parallel()

	t.Run("Create: OK", func(t *testing.T) {
		// t.Parallel()

		teardown()
		defer teardown()

		err := NewCounter("test-counter").Create()
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

		counter := NewCounter("test-counter")
		counter.Create()
		data, err := counter.Find()
		if err != nil {
			t.Fatalf("expected no error, got %s", err)
		}

		now := time.Now()
		formatted := now.Format(time.RFC1123)
		actualData := formatted + "\n"

		if data != actualData {
			t.Fatalf("expected result does not match: %s %s", data, actualData)
		}
	})
}

// TODO run test parallel
// https://github.com/mattetti/filebuffer
func teardown() {
	os.Remove("test-counter")
}
