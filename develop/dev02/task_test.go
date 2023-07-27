package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_unpack(t *testing.T) {
	tests := []struct {
		s       string
		want    string
		wantErr error
		errFunc assert.ErrorAssertionFunc
	}{
		{
			s:       "a4bc2d5e",
			want:    "aaaabccddddde",
			errFunc: assert.NoError,
		},
		{
			s:       "abcd",
			want:    "abcd",
			errFunc: assert.NoError,
		},
		{
			s:       "45",
			want:    "",
			errFunc: assert.Error,
		},
		{
			s:       "",
			want:    "",
			errFunc: assert.NoError,
		},
		{
			s:       `qwe\4\5`,
			want:    "qwe45",
			errFunc: assert.NoError,
		},
		{
			s:       `qwe\45`,
			want:    "qwe44444",
			errFunc: assert.NoError,
		},
		{
			s:       `qwe\\5`,
			want:    `qwe\\\\\`,
			errFunc: assert.NoError,
		},
		{
			s:       `1`,
			want:    "",
			errFunc: assert.Error,
		},
		{
			s:       `\1`,
			want:    "1",
			errFunc: assert.NoError,
		},
		{
			s:       `qwerty\`,
			want:    "",
			wantErr: errBadEscaping,
		},
		{
			s:       `\qwerty`,
			want:    "qwerty",
			errFunc: assert.NoError,
		},
		{
			s:       `4qwerty`,
			want:    "",
			wantErr: errBadFormat,
		},
		{
			s:       `abcd22efg`,
			want:    "",
			wantErr: errBadFormat,
		},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			got, err := unpack(tt.s)

			if tt.errFunc != nil {
				tt.errFunc(t, err)
			}
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
