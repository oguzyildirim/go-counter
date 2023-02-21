package db

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/Oguzyildirim/go-counter/internal"
)

// trying to mock file
var fs fileSystem = osFS{}

type fileSystem interface {
	OpenFile(name string, flag int, perm os.FileMode) (file, error)
	WriteString(s string) (n int, err error)
}

type file interface {
	io.Closer
	io.Reader
	io.ReaderAt
	io.Seeker
	WriteString(s string) (n int, err error)
}

// osFS implements fileSystem using the local disk.
type osFS struct{}

func (osFS) OpenFile(name string, flag int, perm os.FileMode) (file, error) {
	return os.OpenFile(name, flag, perm)
}

func (osFS) WriteString(s string) (n int, err error) {
	file, _ := os.OpenFile("test", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	return file.WriteString(s)
}

func TestNew(t *testing.T) {
	test := New("test")
	if test.dir != "test" {
		t.Errorf("unable to create new driver %s", test.dir)
	}
}

func TestInsert(t *testing.T) {

	oldFs := fs
	mfs := &osFS{}
	fs = mfs
	defer func() {
		fs = oldFs
	}()

	t.Run("Insert: OK", func(t *testing.T) {
		teardown()
		defer teardown()

		driver := New("test")
		data := "testdata1"
		if err := driver.Insert(data); err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	})

	t.Run("Insert: Err", func(t *testing.T) {
		teardown()
		defer teardown()

		defer func() { recover() }()

		driver := New("test")
		data := 123
		err := driver.Insert(data)

		if err == nil {
			t.Fatalf("expected error, got not value")
		}

		var ierr *internal.Error
		if !errors.As(err, &ierr) || ierr.Code() != internal.ErrorCodeInvalidArgument {
			t.Fatalf("expected %T error, got %T : %v", ierr, err, err)
		}
	})
}

func TestGet(t *testing.T) {
	teardown()
	defer teardown()

	oldFs := fs
	// Create and "install" mocked fs:
	mfs := &osFS{}
	fs = mfs
	// Make sure fs is restored after this test:
	defer func() {
		fs = oldFs
	}()

	driver := New("test")
	data := "testdata1"
	driver.Insert(data)
	_, err := driver.Get()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

// workaround
func teardown() {
	os.Remove("test")
}
