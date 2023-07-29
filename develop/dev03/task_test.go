package main

import (
	"math"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const directory = "test_data"

func Test_sortUtility(t *testing.T) {
	type args struct {
		inputFile string
		set       settings
	}
	tests := []struct {
		name        string
		args        args
		wantFile    string
		wantErrFunc assert.ErrorAssertionFunc
	}{
		{
			name: "default settings",
			args: args{
				inputFile: path.Join(directory, "1.txt"),
			},
			wantFile:    path.Join(directory, "1-1a.txt"),
			wantErrFunc: assert.NoError,
		},
		{
			name: "not default column",
			args: args{
				inputFile: path.Join(directory, "1.txt"),
				set: settings{
					Column: 4,
				},
			},
			wantFile:    path.Join(directory, "1-2a.txt"),
			wantErrFunc: assert.NoError,
		},
		{
			name: "by numeric",
			args: args{
				inputFile: path.Join(directory, "1.txt"),
				set: settings{
					Column:    1,
					ByNumeric: true,
				},
			},
			wantFile:    path.Join(directory, "1-3a.txt"),
			wantErrFunc: assert.NoError,
		},
		{
			name: "by month",
			args: args{
				inputFile: path.Join(directory, "1.txt"),
				set: settings{
					Column:  2,
					ByMonth: true,
				},
			},
			wantFile:    path.Join(directory, "1-4a.txt"),
			wantErrFunc: assert.NoError,
		},
		{
			name: "by number with suffix",
			args: args{
				inputFile: path.Join(directory, "1.txt"),
				set: settings{
					Column:          3,
					ByNumericSuffix: true,
				},
			},
			wantFile:    path.Join(directory, "1-5a.txt"),
			wantErrFunc: assert.NoError,
		},
		{
			name: "reverse",
			args: args{
				inputFile: path.Join(directory, "1.txt"),
				set: settings{
					Column:  0,
					Reverse: true,
				},
			},
			wantFile:    path.Join(directory, "1-6a.txt"),
			wantErrFunc: assert.NoError,
		},
		{
			name: "not uniq",
			args: args{
				inputFile: path.Join(directory, "2.txt"),
				set: settings{
					Column:   0,
					UniqOnly: false,
				},
			},
			wantFile:    path.Join(directory, "2-1a.txt"),
			wantErrFunc: assert.NoError,
		},
		{
			name: "uniq",
			args: args{
				inputFile: path.Join(directory, "2.txt"),
				set: settings{
					Column:   0,
					UniqOnly: true,
				},
			},
			wantFile:    path.Join(directory, "2-2a.txt"),
			wantErrFunc: assert.NoError,
		},
		{
			name: "check good",
			args: args{
				inputFile: path.Join(directory, "3.txt"),
				set: settings{
					Column: 0,
					Check:  true,
				},
			},
			wantFile:    path.Join(directory, "empty.txt"),
			wantErrFunc: assert.NoError,
		},
		{
			name: "check bad",
			args: args{
				inputFile: path.Join(directory, "4.txt"),
				set: settings{
					Column: 0,
					Check:  true,
				},
			},
			wantFile:    path.Join(directory, "empty.txt"),
			wantErrFunc: assert.Error,
		},
		{
			name: "check reverse good",
			args: args{
				inputFile: path.Join(directory, "5.txt"),
				set: settings{
					Column:  0,
					Check:   true,
					Reverse: true,
				},
			},
			wantFile:    path.Join(directory, "empty.txt"),
			wantErrFunc: assert.NoError,
		},
		{
			name: "check reverse bad",
			args: args{
				inputFile: path.Join(directory, "4.txt"),
				set: settings{
					Column:  0,
					Check:   true,
					Reverse: true,
				},
			},
			wantFile:    path.Join(directory, "empty.txt"),
			wantErrFunc: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputData, err := os.ReadFile(tt.args.inputFile)
			require.NoError(t, err, "can't read input data file %s", tt.args.inputFile)
			wantData, err := os.ReadFile(tt.wantFile)
			require.NoError(t, err, "can't read want data file %s", tt.wantFile)

			got, err := sortUtility(inputData, tt.args.set)

			tt.wantErrFunc(t, err)
			assert.Equal(t, wantData, got, "WANT:\n%s\nGOT:\n%s", string(wantData), string(got))
		})
	}
}

func Test_parseHumanNumeric(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    float64
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "abracadabra",
			s:       "100abracadabra",
			want:    0,
			wantErr: assert.Error,
		},
		{
			name:    "usual number",
			s:       "1234",
			want:    1234,
			wantErr: assert.NoError,
		},
		{
			name:    "minus number",
			s:       "-200",
			want:    -200,
			wantErr: assert.NoError,
		},
		{
			name:    "zero number",
			s:       "0",
			want:    0,
			wantErr: assert.NoError,
		},
		{
			name:    "k suffix",
			s:       "12k",
			want:    12000,
			wantErr: assert.NoError,
		},
		{
			name:    "minus k suffix",
			s:       "-30k",
			want:    -30000,
			wantErr: assert.NoError,
		},
		{
			name:    "zero k suffix",
			s:       "0k",
			want:    0,
			wantErr: assert.NoError,
		},
		{
			name:    "empty string",
			s:       "",
			want:    0,
			wantErr: assert.Error,
		},
		{
			name:    "m suffix",
			s:       "26m",
			want:    26_000_000,
			wantErr: assert.NoError,
		},
		{
			name:    "bad suffix",
			s:       "123x",
			want:    0,
			wantErr: assert.Error,
		},
	}
	eps := math.Pow(10, -3)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseHumanNumeric(tt.s)

			tt.wantErr(t, err)
			assert.InDelta(t, tt.want, got, eps)
		})
	}
}

func Test_parseMonth(t *testing.T) {
	tests := []struct {
		name  string
		month string
		want  int
	}{
		{
			name:  "empty",
			month: "",
			want:  0,
		},
		{
			name:  "right format",
			month: "feb",
			want:  2,
		},
		{
			name:  "upper case",
			month: "Mar",
			want:  3,
		},
		{
			name:  "full name",
			month: "april",
			want:  4,
		},
		{
			name:  "full name upper",
			month: "JUNE",
			want:  6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseMonth(tt.month)

			assert.Equal(t, tt.want, got)
		})
	}
}
