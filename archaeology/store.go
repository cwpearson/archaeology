package archaeology

import (
	"errors"
	"io"
)

type Store interface {
	MostRecent(path string) (io.Reader, error)
}

type LocalStore struct {
	root string
}

// MostRecent returns a reader corresponding to the most recent version of path
func (s *LocalStore) MostRecent(path string) (io.Reader, error) {
	return nil, errors.New("LocalStore.MostRecent Unimplemented")
}
