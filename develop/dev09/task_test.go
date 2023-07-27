package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_wget(t *testing.T) {
	tests := []struct {
		name          string
		reqURL        string
		wantSubstring string
		errFunc       assert.ErrorAssertionFunc
	}{
		{
			name:          "empty reqURL",
			reqURL:        "",
			wantSubstring: "",
			errFunc:       assert.Error,
		},
		{
			name:          "bad reqURL",
			reqURL:        " https:/\\habra",
			wantSubstring: "",
			errFunc:       assert.Error,
		},
		{
			name:          "good example 1 (go doc)",
			reqURL:        "https://pkg.go.dev/net/http",
			wantSubstring: "Package http provides HTTP client and server implementations.",
			errFunc:       assert.NoError,
		},
		{
			name:          "good example 2 (wiki)",
			reqURL:        "https://ru.wikipedia.org/wiki/Chmod",
			wantSubstring: "команда для изменения прав доступа",
			errFunc:       assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := wget(tt.reqURL)

			assert.True(t, strings.Contains(string(got), tt.wantSubstring))
			tt.errFunc(t, err)
		})
	}
}
