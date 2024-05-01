package app_error

import "errors"

var (
	ErrNotFound          = errors.New("not found")
	ErrDatabase          = errors.New("database error")
	ErrHTTPRequestFailed = errors.New("http request failed")
)
