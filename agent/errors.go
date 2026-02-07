package agent

import "errors"

var (
	ErrToolNotFound = errors.New("tool not found")
	ErrToolInput    = errors.New("invalid tool input")
)
