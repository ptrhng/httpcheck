package main

import (
	"net/http"
	"net/url"
	"time"
)

// NewDefaultOptions returns default httpcheck options.
func NewDefaultOptions() *Options {
	return &Options{
		Method:      http.MethodGet,
		Header:      http.Header{},
		Data:        make(map[string]any),
		FormData:    url.Values{},
		timeout:     time.Second * 10,
		maxBodySize: 1024,
	}
}

// Options configures httpstat.
type Options struct {
	Method         string
	URL            string
	Header         http.Header
	FormData       url.Values
	Data           map[string]any
	timeout        time.Duration
	FollowRedirect bool
	IsForm         bool

	ShowBody    bool
	maxBodySize int
}
