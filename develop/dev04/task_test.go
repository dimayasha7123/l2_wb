package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getHash(t *testing.T) {
	tests := []struct {
		s    string
		want string
	}{
		{
			s:    "абвг",
			want: "1_1_1_1" + strings.Repeat("_", 29),
		},
		{
			s:    "",
			want: strings.Repeat("_", 32),
		},
		{
			s:    "д",
			want: "____1" + strings.Repeat("_", 28),
		},
		{
			s:    "абвгабваба",
			want: "4_3_2_1" + strings.Repeat("_", 29),
		},
		{
			s:    "абвгабвабагггггггггггг",
			want: "4_3_2_13" + strings.Repeat("_", 29),
		},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			got := getHash(tt.s)

			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_getHash_Equals(t *testing.T) {
	tests := []struct {
		s1    string
		s2    string
		equal bool
	}{
		{
			s1:    "абвг",
			s2:    "багв",
			equal: true,
		},
		{
			s1:    "а",
			s2:    "а",
			equal: true,
		},
		{
			s1:    "",
			s2:    "",
			equal: true,
		},
		{
			s1:    "аб",
			s2:    "вг",
			equal: false,
		},
		{
			s1:    "абвгабваба",
			s2:    "гввбббаааа",
			equal: true,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s equal %s: %v", tt.s1, tt.s2, tt.equal),
			func(t *testing.T) {
				h1 := getHash(tt.s1)
				h2 := getHash(tt.s2)

				assert.Equal(t, h1 == h2, tt.equal)
			},
		)
	}
}

func Test_makeAnagrams(t *testing.T) {
	tests := []struct {
		name  string
		words []string
		want  map[string][]string
	}{
		{
			name:  "nil words",
			words: nil,
			want:  make(map[string][]string),
		},
		{
			name:  "empty words",
			words: []string{},
			want:  make(map[string][]string),
		},
		{
			name:  "one word",
			words: []string{"дом"},
			want:  make(map[string][]string),
		},
		{
			name:  "base example",
			words: []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик"},
			want: map[string][]string{
				"листок": {"листок", "слиток", "столик"},
				"пятак":  {"пятак", "пятка", "тяпка"},
			},
		},
		{
			name:  "base example and one word",
			words: []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "дом"},
			want: map[string][]string{
				"листок": {"листок", "слиток", "столик"},
				"пятак":  {"пятак", "пятка", "тяпка"},
			},
		},
		{
			name:  "base example in different registers",
			words: []string{"пЯтАк", "ПяТкА", "тяпка", "ЛИСток", "слиток", "стоЛИК"},
			want: map[string][]string{
				"листок": {"листок", "слиток", "столик"},
				"пятак":  {"пятак", "пятка", "тяпка"},
			},
		},
		{
			name:  "base example with double words",
			words: []string{"пятак", "пятак", "пятка", "тяпка", "листок", "слиток", "слиток", "слиток", "столик"},
			want: map[string][]string{
				"листок": {"листок", "слиток", "столик"},
				"пятак":  {"пятак", "пятка", "тяпка"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := makeAnagrams(tt.words)

			assert.Equal(t, tt.want, got)
		})
	}
}
