package main

import (
	"net/http"
	"time"
)

const contentTypeJSON = "application/json"

// NewDefaultOptions returns default httpcheck options.
func NewDefaultOptions() *Options {
	return &Options{
		Method: http.MethodGet,
		Header: http.Header{
			"Content-Type": []string{contentTypeJSON},
		},
		Data:        make(map[string]any),
		timeout:     time.Second * 10,
		maxBodySize: 1024,
	}
}

// Options configures httpstat.
type Options struct {
	Method  string
	URL     string
	Header  http.Header
	Data    map[string]any
	timeout time.Duration

	ShowBody    bool
	maxBodySize int
}
