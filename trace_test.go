package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTraceHTTP(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		time.Sleep(time.Millisecond * 10)
		fmt.Fprint(rw, "data")
	}))
	defer svr.Close()

	opts := NewDefaultOptions()
	opts.URL = svr.URL
	r, err := Trace(context.Background(), opts)

	require.NoError(t, err)
	assert.GreaterOrEqual(t, r.MetricServerProcessing, int64(10))
	assert.True(t, strings.HasPrefix(r.HTTPVersion, "HTTP/"))
	assert.Equal(t, "200", r.Status)

	b, err := os.ReadFile(r.Output)
	require.NoError(t, err)
	assert.Equal(t, "data", string(b))
}

func TestRaceHTTP_redirect(t *testing.T) {
	svr1 := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprint(rw, "data")
	}))
	defer svr1.Close()
	svr2 := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		http.Redirect(rw, req, svr1.URL, http.StatusMovedPermanently)
	}))
	defer svr2.Close()

	opts := NewDefaultOptions()
	opts.URL = svr2.URL
	r, err := Trace(context.Background(), opts)

	require.NoError(t, err)
	assert.Equal(t, "301", r.Status)

	opts.FollowRedirect = true
	r, err = Trace(context.Background(), opts)
	require.NoError(t, err)
	assert.Equal(t, "200", r.Status)
}
