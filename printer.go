package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
)

const tpl = `Connected to {{ cyan .RemoteAddr }} from {{ .LocalAddr }}

{{ green .HTTPVersion }} {{ cyan .Status }}
{{ range $key, $value := .Header }}
{{- cyan $key }}: {{ join $value ";" | gray }}
{{ end }}

{{- if .ShowBody }}
    {{- if gt .BodySize .BodyMaxSize }}
{{ .BodyString }}{{ cyan "..." }}

{{ green "Body" }} is truncated ({{ .BodyMaxSize }} out of {{ .BodySize }}), stored in: {{ .Output }}
    {{- else }}
{{ .BodyString }}
    {{- end }}
{{- else }}
{{ green "Body" }} stored in: {{ .Output }}
{{- end }}

{{- if .IsHTTPS }}

  DNS Lookup   TCP Connection   TLS Handshake   Server Processing   Content Transfer
[{{fmta .DNSLookup}}   | {{fmta .TCPConnection}}      | {{fmta .TLSHandshake}}     |   {{fmta .ServerProcessing}}       |  {{fmta .ContentTransfer}}       ]
	     |                |               |                   |                  |
    namelookup:{{fmtb .DNSLookup}}      |               |                   |                  |
			connect:{{fmtb .Connect}}     |                   |                  |
				    pretransfer:{{fmtb .PreTransfer}}         |                  |
						      starttransfer:{{fmtb .StartTransfer}}        |
										 total:{{fmtb .Total}}
{{- else }}

  DNS Lookup   TCP Connection   Server Processing   Content Transfer
[{{fmta .DNSLookup}}   | {{fmta .TCPConnection}}      |   {{fmta .ServerProcessing}}       |  {{fmta .ContentTransfer}}       ]
	     |                |                   |                  |
    namelookup:{{fmtb .DNSLookup}}      |                   |                  |
			connect:{{fmtb .Connect}}         |                  |
				      starttransfer:{{fmtb .StartTransfer}}        |
							         total:{{fmtb .Total}}
{{ end }}
`

func fmta(d int64) string {
	return cyan(fmt.Sprintf("%7dms", d))
}

func fmtb(d int64) string {
	return cyan(fmt.Sprintf("%-9s", strconv.Itoa(int(d))+"ms"))
}

func cyan(s string) string {
	return fmt.Sprintf("\033[36m%s\033[0m", s)
}

func gray(s string) string {
	return fmt.Sprintf("\033[38;5;245m%s\033[0m", s)
}

func green(s string) string {
	return fmt.Sprintf("\033[32m%s\033[0m", s)
}

type printOptions struct {
	showBody    bool
	maxBodySize int
	out         io.Writer
}

// PrintOption configures PrintResult.
type PrintOption func(*printOptions)

// WithShowBody configures PrintResult to show response body in the output.
func WithShowBody(s bool) PrintOption {
	return func(opts *printOptions) {
		opts.showBody = s
	}
}

// WithMaxBodySize configures PrintResult to limit the size of the body shown
// in the output.
func WithMaxBodySize(n int) PrintOption {
	return func(opts *printOptions) {
		opts.maxBodySize = n
	}
}

// WithOut configures PrintResult to write the result to the provided destination.
func WithOut(w io.Writer) PrintOption {
	return func(opts *printOptions) {
		opts.out = w
	}
}

type data struct {
	RemoteAddr  string
	LocalAddr   string
	HTTPVersion string
	Status      string
	Header      http.Header

	BodyString string
	BodySize   int64
	Output     string

	BodyMaxSize int
	ShowBody    bool
	IsHTTPS     bool

	DNSLookup        int64
	TCPConnection    int64
	TLSHandshake     int64
	ServerProcessing int64
	ContentTransfer  int64

	Connect       int64
	PreTransfer   int64
	StartTransfer int64
	Total         int64
}

// PrintResult writes the result.
func PrintResult(r *Result, opts ...PrintOption) error {
	options := &printOptions{
		out: os.Stdout,
	}
	for _, o := range opts {
		o(options)
	}

	f, err := os.Open(r.Output)
	if err != nil {
		return err
	}
	info, err := f.Stat()
	if err != nil {
		return err
	}
	bodySize := info.Size()
	body := make([]byte, min(int64(options.maxBodySize), bodySize))
	if _, err := f.Read(body); err != nil {
		return err
	}

	namelookup := r.MetricDNSLookup
	connect := namelookup + r.MetricTCPConnection
	preTransfer := connect + r.MetricTLSHandshake
	startTransfer := preTransfer + r.MetricServerProcessing
	total := startTransfer + r.MetricContentTransfer

	d := data{
		RemoteAddr:  r.RemoteAddr,
		LocalAddr:   r.LocalAddr,
		HTTPVersion: r.HTTPVersion,
		Status:      r.Status,
		Header:      r.Header,
		Output:      r.Output,

		BodyString:  string(body),
		BodySize:    bodySize,
		BodyMaxSize: options.maxBodySize,
		ShowBody:    options.showBody,
		IsHTTPS:     strings.HasPrefix(r.URL, "https://"),

		DNSLookup:        r.MetricDNSLookup,
		TCPConnection:    r.MetricTCPConnection,
		TLSHandshake:     r.MetricTLSHandshake,
		ServerProcessing: r.MetricServerProcessing,
		ContentTransfer:  r.MetricContentTransfer,

		Connect:       connect,
		PreTransfer:   preTransfer,
		StartTransfer: startTransfer,
		Total:         total,
	}

	funcs := template.FuncMap{
		"join":  strings.Join,
		"cyan":  cyan,
		"gray":  gray,
		"green": green,
		"fmta":  fmta,
		"fmtb":  fmtb,
	}
	tmpl, err := template.New("result").Funcs(funcs).Parse(tpl)
	if err != nil {
		return err
	}

	return tmpl.Execute(os.Stdout, d)
}
