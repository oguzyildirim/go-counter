package storage

import (
	"time"

	"github.com/Oguzyildirim/go-counter/internal"
	"github.com/Oguzyildirim/go-counter/tools/db"
)

// Counter represents the repository used for interacting with Counter records
type Counter struct {
	db *db.Driver
}

// NewCounter instantiates the Counter repository
func NewCounter(dir string) *Counter {
	return &Counter{
		db: db.New(dir),
	}
}

// Create inserts a new ıser record
func (c *Counter) Create() error {
	now := time.Now()
	formatted := now.Format(time.RFC1123)
	data := formatted + "\n"
	err := c.db.Insert(data)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "create failed")
	}
	return nil
}

// Find
func (c *Counter) Find() (string, error) {
	result, err := c.db.Get()
	if err != nil {
		return "", internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "find failed")
	}

	return result, nil
}
