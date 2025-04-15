package common

import "errors"

var (
	ErrNoItems = errors.New("Items must have at least one item")
)
