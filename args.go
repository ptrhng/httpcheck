package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	separatorHeader      = ":"
	separatorDataString  = "="
	separatorDataRawJSON = ":="
)

var (
	// a longer separtor appear first so that they can be detected before
	// a shorter separator that is a substring of the former.
	// example:
	// ":=" detected before ":"
	separators = []string{
		separatorDataRawJSON,
		separatorHeader,
		separatorDataString,
	}

	methods = []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
	}
)

func findSeparator(s string) string {
	found := ""
	startAt := len(s)
	for _, v := range separators {
		idx := strings.Index(s, v)
		if idx >= 0 && idx < startAt {
			found = v
			startAt = idx
		}
	}

	return found
}

// ParseArgs parses args and update options.
func ParseArgs(args []string, opts *Options) error {
	if !slices.Contains(methods, args[0]) {
		args = append([]string{opts.Method}, args...)
	}
	opts.Method = args[0]
	opts.URL = args[1]
	if !strings.HasPrefix(opts.URL, "http://") && !strings.HasPrefix(opts.URL, "https://") {
		opts.URL = "http://" + opts.URL
	}

	for _, arg := range args[2:] {
		sep := findSeparator(arg)
		tokens := strings.SplitN(arg, sep, 2)
		if len(tokens) == 1 {
			tokens = append(tokens, "")
		}
		logrus.Debugf("separator: '%s', toekns: %v", sep, tokens)
		k, v := tokens[0], tokens[1]

		switch sep {
		case separatorHeader:
			opts.Header.Add(k, v)
		case separatorDataString:
			if opts.IsForm {
				opts.FormData.Add(k, v)
			} else {
				opts.Data[k] = v
			}
		case separatorDataRawJSON:
			if opts.IsForm {
				return fmt.Errorf("cannot use json value type '%s' with --form", arg)
			}

			var o any
			if err := json.Unmarshal([]byte(v), &o); err != nil {
				return fmt.Errorf("'%s' is not a valid json", v)
			}
			opts.Data[k] = o
		default:
			return fmt.Errorf("'%s' is not a valid request item", arg)
		}
	}

	return nil
}
