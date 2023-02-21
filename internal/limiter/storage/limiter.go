package storage

import (
	"time"

	"github.com/Oguzyildirim/go-counter/internal"
	"github.com/Oguzyildirim/go-counter/tools/db"
)

// Limiter represents the repository used for interacting with Limiter records
type Limiter struct {
	db *db.Driver
}

// NewLimiter instantiates the Limiter repository
func NewLimiter(dir string) *Limiter {
	return &Limiter{
		db: db.New(dir),
	}
}

// Create inserts a new ıser record
func (l *Limiter) Create(key string) error {
	now := time.Now()
	formatted := now.Format(time.RFC1123)
	data := formatted + key + "\n"
	err := l.db.Insert(data)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "create failed")
	}
	return nil
}

// Find
func (l *Limiter) Find() (string, error) {
	result, err := l.db.Get()
	if err != nil {
		return "", internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "find failed")
	}

	return result, nil
}
