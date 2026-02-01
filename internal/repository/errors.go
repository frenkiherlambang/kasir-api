package repository

import "errors"

// ErrNotFound is returned when an entity is not found.
var ErrNotFound = errors.New("not found")
