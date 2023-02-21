package service

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Oguzyildirim/go-counter/internal"
)

const windowLimit = 1

// CounterRepository defines the datastore handling persisting Counter records
type CounterRepository interface {
	Create() error
	Find() (string, error)
}

// Counter defines the application service in charge of interacting with Counters
type Counter struct {
	repo CounterRepository
}

// NewCounter
func NewCounter(repo CounterRepository) *Counter {
	return &Counter{
		repo: repo,
	}
}

// Create insert a new record
func (c *Counter) Create() error {
	err := c.repo.Create()
	if err != nil {
		return fmt.Errorf("repo create: %w", err)
	}

	return nil
}

// Find reads the record
func (c *Counter) Find() (string, error) {
	num, err := c.repo.Find()
	if err != nil {
		return "", fmt.Errorf("repo create: %w", err)
	}

	composed, err := compose(num)
	if err != nil {
		return "", fmt.Errorf("compose: %w", err)
	}
	count := strconv.Itoa(composed)

	return count, nil
}

func compose(result string) (int, error) {
	var count int
	var times []time.Time

	scanner := bufio.NewScanner(strings.NewReader(result))
	for scanner.Scan() {
		data := scanner.Text()
		t, err := time.Parse(time.RFC1123, data)
		if err != nil {
			return 0, fmt.Errorf("time.Parse: %w", err)
		}
		times = append(times, t)
	}
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("NewScanner: %w", err)
	}

	now := time.Now()
	window := now.Add(time.Minute * -windowLimit)
	for i := len(times) - 1; i >= 0; i-- {
		if internal.InTimeSpan(window, now, times[i]) {
			count = count + 1
		} else {
			break
		}
	}

	return count, nil
}
