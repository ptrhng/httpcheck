package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func goldenAssert(t *testing.T, filename, actual string) {
	fpath := filepath.Clean("testdata/" + filename)
	b, err := os.ReadFile(fpath)
	require.NoError(t, err)
	assert.Equal(t, string(b), actual)
}

func TestPrintResult_success(t *testing.T) {
	cases := []struct {
		name   string
		opts   []PrintOption
		result *Result
	}{
		{
			name: "http",
			result: &Result{
				URL:         "http://1.1.1.1",
				RemoteAddr:  "1.1.1.1:80",
				LocalAddr:   "192.168.1.1:63917",
				HTTPVersion: "HTTP/1.1",
				Status:      "200",
				Headers: []Header{
					{Name: "Content-Ranges", Value: "bytes"},
					{Name: "Expires", Value: "-1"},
					{Name: "Server", Value: "test"},
				},
				Output:                 "testdata/response_body.txt",
				MetricDNSLookup:        10,
				MetricTCPConnection:    10,
				MetricServerProcessing: 10,
				MetricContentTransfer:  10,
			},
		},
		{
			name: "https",
			result: &Result{
				URL:         "https://1.1.1.1",
				RemoteAddr:  "1.1.1.1:443",
				LocalAddr:   "192.168.1.1:63917",
				HTTPVersion: "HTTP/2.0",
				Status:      "200",
				Headers: []Header{
					{Name: "Content-Ranges", Value: "bytes"},
					{Name: "Expires", Value: "-1"},
					{Name: "Server", Value: "test"},
				},
				Output:                 "testdata/response_body.txt",
				MetricDNSLookup:        10,
				MetricTCPConnection:    10,
				MetricTLSHandshake:     10,
				MetricServerProcessing: 10,
				MetricContentTransfer:  10,
			},
		},
		{
			name: "showbody",
			opts: []PrintOption{WithShowBody(true), WithMaxBodySize(100)},
			result: &Result{
				URL:         "https://1.1.1.1",
				RemoteAddr:  "1.1.1.1:443",
				LocalAddr:   "192.168.1.1:63917",
				HTTPVersion: "HTTP/2.0",
				Status:      "200",
				Headers: []Header{
					{Name: "Content-Ranges", Value: "bytes"},
					{Name: "Expires", Value: "-1"},
					{Name: "Server", Value: "test"},
				},
				Output:                 "testdata/response_body.txt",
				MetricDNSLookup:        10,
				MetricTCPConnection:    10,
				MetricTLSHandshake:     10,
				MetricServerProcessing: 10,
				MetricContentTransfer:  10,
			},
		},
		{
			name: "showbody_truncated",
			opts: []PrintOption{WithShowBody(true), WithMaxBodySize(5)},
			result: &Result{
				URL:         "https://1.1.1.1",
				RemoteAddr:  "1.1.1.1:443",
				LocalAddr:   "192.168.1.1:63917",
				HTTPVersion: "HTTP/2.0",
				Status:      "200",
				Headers: []Header{
					{Name: "Content-Ranges", Value: "bytes"},
					{Name: "Expires", Value: "-1"},
					{Name: "Server", Value: "test"},
				},
				Output:                 "testdata/response_body.txt",
				MetricDNSLookup:        10,
				MetricTCPConnection:    10,
				MetricTLSHandshake:     10,
				MetricServerProcessing: 10,
				MetricContentTransfer:  10,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			err := PrintResult(tc.result, append(tc.opts, WithOut(buf), WithNoColor())...)
			require.NoError(t, err)
			goldenAssert(t, tc.name+".golden", buf.String())
		})
	}
}
