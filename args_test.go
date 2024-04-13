package main

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseArg_defaults(t *testing.T) {
	args := []string{"www.example.com"}
	opts := NewDefaultOptions()
	err := ParseArgs(args, opts)

	require.NoError(t, err)
	assert.Equal(t, "http://www.example.com", opts.URL)
	assert.Equal(t, 1, len(opts.Header))
	assert.Equal(t, "application/json", opts.Header.Get("Content-Type"))
}

func TestParseArg_post(t *testing.T) {
	args := []string{http.MethodPost, "https://www.example.com"}
	opts := NewDefaultOptions()
	err := ParseArgs(args, opts)

	require.NoError(t, err)
	assert.Equal(t, http.MethodPost, opts.Method)
}

func TestParseArg_header(t *testing.T) {
	args := []string{
		"https://www.example.com",
		"k:v1",
		"k:v2",
	}
	opts := NewDefaultOptions()
	err := ParseArgs(args, opts)

	require.NoError(t, err)
	assert.Equal(t, []string{"v1", "v2"}, opts.Header.Values("k"))
}

func TestParseArg_json(t *testing.T) {
	args := []string{
		"https://www.example.com",
		"k:=[1, 2, 3]",
	}
	opts := NewDefaultOptions()
	err := ParseArgs(args, opts)

	require.NoError(t, err)

	b, err := json.Marshal(opts.Data)
	require.NoError(t, err)
	assert.Equal(t, `{"k":[1,2,3]}`, string(b))
}

func TestPArseArg_errors(t *testing.T) {
	cases := []struct {
		name string
		args []string
	}{
		{
			name: "incomplete flag",
			args: []string{"www.example.com", "-"},
		},
		{
			name: "unknown request item",
			args: []string{"www.example.com", "unknown"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			opts := NewDefaultOptions()
			err := ParseArgs(tc.args, opts)
			require.Error(t, err)
		})
	}
}
