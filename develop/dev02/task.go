package main

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

/*
=== Задача на распаковку ===

Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
	- "a4bc2d5e" => "aaaabccddddde"
	- "abcd" => "abcd"
	- "45" => "" (некорректная строка)
	- "" => ""
Дополнительное задание: поддержка escape - последовательностей
	- qwe\4\5 => qwe45 (*)
	- qwe\45 => qwe44444 (*)
	- qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.

Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

var (
	errBadEscaping = errors.New("bad string escaping")
	errBadFormat   = errors.New("bad string format")
)

const escapeRune = '\\'

type part struct {
	Symbol rune
	Count  int
}

func (p part) isSymbol() bool {
	return p.Symbol != rune(0)
}

func (p part) isCount() bool {
	return !p.isSymbol()
}

func unpack(s string) (string, error) {
	parts := make([]part, 0, len(s))

	wasSlash := false
	for _, r := range s {
		switch {
		case wasSlash:
			parts = append(parts, part{Symbol: r})
			wasSlash = false
		case r == escapeRune:
			wasSlash = true
		case unicode.IsDigit(r):
			parts = append(parts, part{Count: int(r) - int('0')})
		default:
			parts = append(parts, part{Symbol: r})
		}
	}
	if wasSlash {
		return "", errBadEscaping
	}

	if len(parts) > 0 && parts[0].isCount() {
		return "", errBadFormat
	}

	sb := strings.Builder{}
	for i := 0; i < len(parts); i++ {
		p := parts[i]
		if p.isCount() {
			if i != len(parts)-1 && parts[i+1].isCount() {
				return "", errBadFormat
			}
			continue
		}

		count := 1
		if i != len(parts)-1 {
			next := parts[i+1]
			if next.isCount() {
				count = next.Count
			}
		}
		sb.Grow(count)
		for j := 0; j < count; j++ {
			sb.WriteRune(p.Symbol)
		}
	}

	return sb.String(), nil
}

func main() {
	fmt.Println(unpack(``))
}
