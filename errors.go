package ratelimit

import "errors"

var (
	ErrLimitExceeded = errors.New("ratelimit: limit exceeded")
)
