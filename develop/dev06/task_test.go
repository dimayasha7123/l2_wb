package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	space = ' '
	tab   = '\t'
)

func Test_cut(t *testing.T) {
	type args struct {
		ss     []string
		fields []bool
		del    rune
		sep    bool
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "base example",
			args: args{
				ss: []string{
					"aaaaa 9999 JAN",
					"ccccc 1111 FEB",
					"bbbbb 5555 SEP",
				},
				fields: []bool{true, false, true},
				del:    space,
			},
			want: []string{
				"aaaaa JAN",
				"ccccc FEB",
				"bbbbb SEP",
			},
		},
		{
			name: "sep=false",
			args: args{
				ss: []string{
					"aaaaa 9999 JAN",
					"ccccc 1111 FEB",
					"bbbbb 5555 SEP",
					"one_word",
					"",
				},
				fields: []bool{true, false, true},
				del:    space,
				sep:    false,
			},
			want: []string{
				"aaaaa JAN",
				"ccccc FEB",
				"bbbbb SEP",
				"one_word",
				"",
			},
		},
		{
			name: "sep=false",
			args: args{
				ss: []string{
					"aaaaa 9999 JAN",
					"ccccc 1111 FEB",
					"bbbbb 5555 SEP",
					"one_word",
					"",
				},
				fields: []bool{true, false, true},
				del:    space,
				sep:    true,
			},
			want: []string{
				"aaaaa JAN",
				"ccccc FEB",
				"bbbbb SEP",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cut(tt.args.ss, tt.args.fields, tt.args.del, tt.args.sep)

			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_cutString(t *testing.T) {
	type args struct {
		s      string
		fields []bool
		del    rune
		buf    bytes.Buffer
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "nil fiedls",
			args: args{
				s:      "dimya masha petya",
				fields: nil,
				del:    space,
			},
			want: "",
		},
		{
			name: "empty fiedls",
			args: args{
				s:      "dimya masha petya",
				fields: []bool{},
				del:    space,
			},
			want: "",
		},
		{
			name: "false fiedls",
			args: args{
				s:      "dimya masha petya",
				fields: []bool{false},
				del:    space,
			},
			want: "",
		},
		{
			name: "empty string",
			args: args{
				s:      "",
				fields: []bool{false, true, true},
				del:    space,
			},
			want: "",
		},
		{
			name: "no dels",
			args: args{
				s:      "abcdefg",
				fields: []bool{true, true, true},
				del:    space,
			},
			want: "abcdefg",
		},
		{
			name: "no dels and no first field",
			args: args{
				s:      "abcdefg",
				fields: []bool{false, true},
				del:    space,
			},
			want: "",
		},
		{
			name: "len(fields) = len(split(s)",
			args: args{
				s:      "dima masha petya galya",
				fields: []bool{false, true, false, true},
				del:    space,
			},
			want: "masha galya",
		},
		{
			name: "len(fields) = len(split(s) with another separator",
			args: args{
				s:      "dima\tmasha petya\tgalya 7123",
				fields: []bool{false, true, true, true},
				del:    tab,
			},
			want: "masha petya\tgalya 7123",
		},
		{
			name: "len(fields) < len(split(s)",
			args: args{
				s:      "dima masha petya galya",
				fields: []bool{true, false, true},
				del:    space,
			},
			want: "dima petya",
		},
		{
			name: "len(fields) > len(split(s)",
			args: args{
				s:      "dima masha petya galya",
				fields: []bool{true, false, false, false, false, true},
				del:    space,
			},
			want: "dima",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cutString(tt.args.s, tt.args.fields, tt.args.del, tt.args.buf)

			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_increaseRetToSize(t *testing.T) {
	type args struct {
		ret  []bool
		size int
	}
	tests := []struct {
		name string
		args args
		want []bool
	}{
		{
			name: "nil input to 0 size",
			args: args{
				ret:  nil,
				size: 0,
			},
			want: []bool{},
		},
		{
			name: "empty input to 0 size",
			args: args{
				ret:  []bool{},
				size: 0,
			},
			want: []bool{},
		},
		{
			name: "nil input to 3 size",
			args: args{
				ret:  nil,
				size: 3,
			},
			want: []bool{false, false, false},
		},
		{
			name: "empty input to 3 size",
			args: args{
				ret:  []bool{},
				size: 3,
			},
			want: []bool{false, false, false},
		},
		{
			name: "base input to bigger size",
			args: args{
				ret:  []bool{true, false, true},
				size: 4,
			},
			want: []bool{true, false, true, false},
		},
		{
			name: "base input to same size",
			args: args{
				ret:  []bool{true, false, true},
				size: 3,
			},
			want: []bool{true, false, true},
		},
		{
			name: "base input to smaller size",
			args: args{
				ret:  []bool{true, false, true},
				size: 2,
			},
			want: []bool{true, false, true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := increaseRetToSize(tt.args.ret, tt.args.size)

			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_parseFields(t *testing.T) {
	var tests = []struct {
		name    string
		s       string
		want    []bool
		wantErr error
	}{
		{
			name:    "empty",
			s:       "",
			want:    []bool{},
			wantErr: nil,
		},
		{
			name:    "just number",
			s:       "3",
			want:    []bool{false, false, true},
			wantErr: nil,
		},
		{
			name:    "few numbers",
			s:       "3,1,4",
			want:    []bool{true, false, true, true},
			wantErr: nil,
		},
		{
			name:    "segment",
			s:       "3-4",
			want:    []bool{false, false, true, true},
			wantErr: nil,
		},
		{
			name:    "numbers and segments",
			s:       "3,1,5-6,2,8-9",
			want:    []bool{true, true, true, false, true, true, false, true, true},
			wantErr: nil,
		},
		{
			name:    "not number",
			s:       "dimya",
			want:    nil,
			wantErr: errFieldMustBeInteger,
		},
		{
			name:    "not number in segment",
			s:       "43-vasya",
			want:    nil,
			wantErr: errFieldMustBeInteger,
		},
		{
			name:    "few tire",
			s:       "3-5-7",
			want:    nil,
			wantErr: errBadFieldsFormat,
		},
		{
			name:    "equal numbers in segment",
			s:       "8-8",
			want:    nil,
			wantErr: errLeftColumnMoreOrEqualRightColumn,
		},
		{
			name:    "bad numbers in segment",
			s:       "8-3",
			want:    nil,
			wantErr: errLeftColumnMoreOrEqualRightColumn,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseFields(tt.s)

			assert.Equal(t, tt.want, got)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
