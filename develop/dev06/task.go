package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*
=== Утилита cut ===

Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Дополнительно:
-i - "input" - файл для чтения (если не передан, читает stdin)

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	var (
		fields    = ""
		delimiter = "\t"
		separated = false
		inputFile = ""
	)
	flag.StringVar(&fields, "f", fields, "required fields in this format N (only column N counting from)"+
		"or N-M (from N to M), can combine formats using comma, like this \"1, 3-5, 7\", so it will be 1, 3, 4, 5, 7 columns")
	flag.StringVar(&delimiter, "d", delimiter, "delimiter for splitting")
	flag.BoolVar(&separated, "s", separated, "only lines with delimiters")
	flag.StringVar(&inputFile, "i", inputFile, "path to input file, if empty - using STDIN")
	flag.Parse()

	if len(fields) == 0 {
		fmt.Fprintf(os.Stderr, "fields (flag \"f\") must be non empty\n")
		os.Exit(2)
	}

	if len(delimiter) == 0 {
		fmt.Fprintf(os.Stderr, "delimiter (flag \"d\") must be non empty\n")
		os.Exit(2)
	}
	if len(delimiter) != 1 {
		fmt.Fprintf(os.Stderr, "delimiter length (flag \"d\") must be exact one (only one symbol)\n")
		os.Exit(2)
	}
	delimiterRune := []rune(delimiter)[0]

	fieldsSlice, err := parseFields(fields)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't parse fields: %v\n", err)
		os.Exit(2)
	}

	// p.s. Меня чет переклинило и я сделал чтение ввода из первого аргумента, а не STDIN (
	var input string
	if inputFile == "" {
		input = flag.Arg(0)
	} else {
		data, err := os.ReadFile(inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "can't read file %s: %v\n", inputFile, err)
			os.Exit(1)
		}
		input = string(data)
	}

	if len(input) == 0 {
		fmt.Fprintf(os.Stderr, "empty input\n")
		os.Exit(2)
	}

	ss := strings.Split(input, "\n")
	out := cut(ss, fieldsSlice, delimiterRune, separated)
	for _, s := range out {
		fmt.Println(s)
	}
}

var (
	errBadFieldsFormat                  = errors.New("fields must be in format \"N\" or \"N-M\"")
	errFieldMustBeInteger               = errors.New("fields must be integers")
	errLeftColumnMoreOrEqualRightColumn = errors.New("left column mest be less right column")
)

func parseFields(fields string) ([]bool, error) {
	if len(fields) == 0 {
		return []bool{}, nil
	}
	fieldGroups := strings.Split(fields, ",")
	ret := make([]bool, 0)
	for _, fg := range fieldGroups {
		if strings.Contains(fg, "-") {
			fgs := strings.Split(fg, "-")
			if len(fgs) != 2 {
				return nil, errBadFieldsFormat
			}

			left, err := strconv.Atoi(fgs[0])
			if err != nil {
				return nil, errFieldMustBeInteger
			}
			right, err := strconv.Atoi(fgs[1])
			if err != nil {
				return nil, errFieldMustBeInteger
			}

			if left >= right {
				return nil, errLeftColumnMoreOrEqualRightColumn
			}

			left--
			ret = increaseRetToSize(ret, right)
			for i := left; i < right; i++ {
				ret[i] = true
			}
		} else {
			f, err := strconv.Atoi(fg)
			if err != nil {
				return nil, errFieldMustBeInteger
			}

			ret = increaseRetToSize(ret, f)
			ret[f-1] = true
		}
	}
	return ret, nil
}

func increaseRetToSize(ret []bool, size int) []bool {
	if ret == nil {
		return make([]bool, size)
	}
	if len(ret) >= size {
		return ret
	}
	newRet := make([]bool, size)
	copy(newRet, ret)
	return newRet
}

func cut(ss []string, fields []bool, del rune, sep bool) []string {
	ret := make([]string, 0, len(ss))
	buf := bytes.Buffer{}
	for _, s := range ss {
		if sep && !strings.ContainsRune(s, del) {
			continue
		}
		ret = append(ret, cutString(s, fields, del, buf))
	}
	return ret
}

func cutString(s string, fields []bool, del rune, buf bytes.Buffer) string {
	buf.Reset()
	ss := strings.Split(s, string(del))
	first := true
	for i := 0; i < min(len(ss), len(fields)); i++ {
		if !fields[i] {
			continue
		}
		if !first {
			buf.WriteRune(del)
		}
		first = false
		buf.WriteString(ss[i])
	}
	return buf.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
