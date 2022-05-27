package domain

import "errors"

var (
	ErrInvalidDependency = errors.New("provided dependency is nil")

	ErrInvalidToken = errors.New("provided jwt is invalid")
)
