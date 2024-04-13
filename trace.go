package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"maps"
	"net/http"
	"net/http/httptrace"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

// Result is the performance metric returned by Trace function.
type Result struct {
	URL         string
	RemoteAddr  string
	LocalAddr   string
	HTTPVersion string
	Status      string
	Header      http.Header

	Output string

	MetricDNSLookup        int64
	MetricTCPConnection    int64
	MetricTLSHandshake     int64
	MetricServerProcessing int64
	MetricContentTransfer  int64
}

func diffMills(t1, t2 time.Time) int64 {
	return t1.Sub(t2).Milliseconds()
}

func close(c io.Closer) {
	if err := c.Close(); err != nil {
		logrus.Warn(err)
	}
}

// Trace sends a request to the specified URL and returns
// a performance metirc.
func Trace(ctx context.Context, opts *Options) (*Result, error) {
	ctx, cancel := context.WithTimeout(ctx, opts.timeout)
	defer cancel()

	var body io.Reader
	if len(opts.Data) > 0 {
		b, err := json.Marshal(opts.Data)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(b)
	}

	req, err := http.NewRequestWithContext(ctx, opts.Method, opts.URL, body)
	if err != nil {
		return nil, err
	}
	req.Header = opts.Header

	r := &Result{
		URL: opts.URL,
	}
	var t0, t1, t2, t3, t4, t5, t6, t7, t8 time.Time
	trace := &httptrace.ClientTrace{
		DNSStart: func(di httptrace.DNSStartInfo) {
			t0 = time.Now()
		},
		DNSDone: func(di httptrace.DNSDoneInfo) {
			t1 = time.Now()
		},
		ConnectStart: func(network, addr string) {
			t2 = time.Now()
		},
		ConnectDone: func(network, addr string, err error) {
			if err != nil {
				logrus.Error(err)
				return
			}

			t3 = time.Now()
			r.RemoteAddr = addr

		},
		TLSHandshakeStart: func() {
			t4 = time.Now()
		},
		TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
			if err != nil {
				logrus.Error(err)
				return
			}
			t5 = time.Now()
		},
		GotConn: func(gci httptrace.GotConnInfo) {
			t6 = time.Now()
			r.LocalAddr = gci.Conn.LocalAddr().String()
		},
		GotFirstResponseByte: func() {
			t7 = time.Now()
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer close(resp.Body)

	f, err := os.CreateTemp("", "")
	if err != nil {
		return nil, err
	}
	defer close(f)
	if _, err := io.Copy(f, resp.Body); err != nil {
		return nil, err
	}
	t8 = time.Now()

	r.MetricDNSLookup = diffMills(t1, t0)
	r.MetricTCPConnection = diffMills(t3, t2)
	r.MetricTLSHandshake = diffMills(t5, t4)
	r.MetricServerProcessing = diffMills(t7, t6)
	r.MetricContentTransfer = diffMills(t8, t7)

	r.HTTPVersion = resp.Proto
	r.Header = maps.Clone(resp.Header)
	r.Status = strconv.Itoa(resp.StatusCode)
	r.Output = f.Name()

	return r, nil
}