package storage

import "errors"

var (
	ErrNoConnection = errors.New("cannot connect to database")
)
