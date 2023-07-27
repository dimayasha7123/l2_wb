package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getDataFromString(s string) [][]byte {
	return bytes.Split([]byte(s), []byte{byte(stringSeparator)})
}

func Test_grep(t *testing.T) {
	type args struct {
		data    [][]byte
		pattern string
		set     settings
	}
	tests := []struct {
		name    string
		args    args
		want    grepRows
		errFunc assert.ErrorAssertionFunc
	}{
		{
			name: "default settings test",
			args: args{
				data:    getDataFromString("London\n2\n3\n4\nOoh\noh\nOH\n8\n9\nLondo."),
				pattern: "o.",
				set:     settings{},
			},
			want: grepRows{
				{
					LineNum:      1,
					Data:         []byte("London"),
					MatchIndexes: [][]int{{1, 3}, {4, 6}},
				},
				{
					LineNum:      5,
					Data:         []byte("Ooh"),
					MatchIndexes: [][]int{{1, 3}},
				},
				{
					LineNum:      6,
					Data:         []byte("oh"),
					MatchIndexes: [][]int{{0, 2}},
				},
				{
					LineNum:      10,
					Data:         []byte("Londo."),
					MatchIndexes: [][]int{{1, 3}, {4, 6}},
				},
			},
			errFunc: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := grep(tt.args.data, tt.args.pattern, tt.args.set)

			assert.Equal(t, tt.want, got)
			tt.errFunc(t, err)
		})
	}
}
