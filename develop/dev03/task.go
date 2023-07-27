package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

/*
=== Утилита sort ===

Отсортировать строки (man sort)
Основное

Поддержать ключи

-k — указание колонки для сортировки
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительное

Поддержать ключи

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учётом суффиксов

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	var set settings
	flag.IntVar(&set.Column, "k", 1, "указание колонки для сортировки")
	flag.BoolVar(&set.ByNumeric, "n", false, "сортировать по числовому значению")
	flag.BoolVar(&set.Reverse, "r", false, "сортировать в обратном порядке")
	flag.BoolVar(&set.UniqOnly, "u", false, "не выводить повторяющиеся строки")
	flag.BoolVar(&set.ByMonth, "M", false, "сортировать по названию месяца")
	flag.BoolVar(&set.IgnoreTailSpaces, "b", false, "игнорировать хвостовые пробелы")
	flag.BoolVar(&set.Check, "c", false, "проверять отсортированы ли данные")
	flag.BoolVar(&set.ByNumericSuffix, "h", false, "сортировать по числовому значению с учётом суффиксов")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Использование sort:\n")
		fmt.Fprintf(os.Stderr, "sort [ФЛАГ]... [ФАЙЛ]...\n")
		fmt.Fprintf(os.Stderr, "Если в качестве первого файла указан \"-\", то будет считан STDIN,"+
			" иначе будут считаны, смержены и отсорчены данные всех файлов")
		flag.PrintDefaults()
	}
	flag.Parse()

	if set.Column < 1 {
		fmt.Fprintf(os.Stderr, "K не может быть меньше 1\n")
		os.Exit(2)
	}
	set.Column--

	files := flag.Args()
	if len(files) == 0 {
		fmt.Fprintf(os.Stderr, "Не указаны файлы или \"-\" для чтения из STDIN\n")
		os.Exit(2)
	}

	var data []byte
	if files[0] == "-" {
		var err error
		data, err = io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Не могу считать из STDIN: %v\n", err)
			os.Exit(1)
		}
	} else {
		for _, file := range files {
			fileData, err := os.ReadFile(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Не могу прочитать файл: %v\n", err)
				continue
			}

			if len(data) != 0 && data[len(data)-1] != endLine {
				data = append(data, endLine)
			}
			data = append(data, fileData...)
		}
	}

	output, err := sortUtility(data, set)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	fmt.Println(set.Column)
	fmt.Println(string(output))
}

type settings struct {
	Column           int  // номер колонки (параметр сортировки)
	ByNumeric        bool // параметр сортировки
	Reverse          bool // в конце реверснуть
	UniqOnly         bool // в конце добавлять в мапу и проверять, был ли он
	ByMonth          bool // параметр сортировки
	IgnoreTailSpaces bool // где-то при сплите обрезать пробелы
	Check            bool // в начале смотреть, если тру, то переходить к проверке, а не к сортировке
	ByNumericSuffix  bool // параметр сортировки
}

func sortUtility(data []byte, set settings) ([]byte, error) {
	sd := newSortData(data, set)
	if set.Check {
		var sorted bool
		if set.Reverse {
			sorted = sort.IsSorted(sort.Reverse(sd))
		} else {
			sorted = sort.IsSorted(sd)
		}

		if sorted {
			return nil, nil
		}
		return nil, errors.New("неправильный порядок")
	}

	if set.Reverse {
		sort.Sort(sort.Reverse(sd))
	} else {
		sort.Sort(sd)
	}

	return sd.toBytes(), nil
}

type sortData struct {
	data [][][]byte
	set  settings
}

func newSortData(data []byte, set settings) *sortData {
	retData := make([][][]byte, 0)
	lines := bytes.Split(data, []byte{byte(endLine)})
	for _, line := range lines {
		if set.IgnoreTailSpaces {
			line = bytes.TrimRight(line, string(delimiter))
		}
		words := bytes.Split(line, []byte{byte(delimiter)})
		retData = append(retData, words)
	}

	return &sortData{
		data: retData,
		set:  set,
	}
}

func (s *sortData) toBytes() []byte {
	lines := make([]string, 0, len(s.data))
	for _, line := range s.data {
		lines = append(lines, string(bytes.Join(line, []byte{byte(delimiter)})))
	}
	if s.set.UniqOnly {
		uLines := make([]string, 0, len(lines))
		dubs := make(map[string]struct{})

		for _, line := range lines {
			_, ok := dubs[line]
			if ok {
				continue
			}
			uLines = append(uLines, line)
			dubs[line] = struct{}{}
		}

		lines = uLines
	}
	return []byte(strings.Join(lines, string(endLine)))
}

func (s *sortData) Len() int      { return len(s.data) }
func (s *sortData) Swap(i, j int) { s.data[i], s.data[j] = s.data[j], s.data[i] }
func (s *sortData) Less(i, j int) bool {
	iLine := s.data[i]
	jLine := s.data[j]
	var (
		iCell []byte
		jCell []byte
	)
	if s.set.Column < len(iLine) {
		iCell = iLine[s.set.Column]
	}
	if s.set.Column < len(jLine) {
		jCell = jLine[s.set.Column]
	}
	if iCell == nil {
		return true
	}
	if jCell == nil {
		return false
	}

	switch {
	case s.set.ByNumeric:
		iNum, err := strconv.ParseFloat(string(iCell), 64)
		if err != nil {
			return true
		}
		jNum, err := strconv.ParseFloat(string(jCell), 64)
		if err != nil {
			return false
		}
		return iNum < jNum

	case s.set.ByNumericSuffix:
		iNum, err := parseHumanNumeric(string(iCell))
		if err != nil {
			return true
		}
		jNum, err := parseHumanNumeric(string(jCell))
		if err != nil {
			return false
		}
		return iNum < jNum

	case s.set.ByMonth:
		iMonth := parseMonth(string(iCell))
		jMonth := parseMonth(string(jCell))
		return iMonth < jMonth

	}

	return string(iCell) < string(jCell)
}

const (
	endLine      = '\n'
	delimiter    = ' '
	monthAbbrLen = 3
)

func parseMonth(month string) int {
	month = strings.ToLower(month)
	if len(month) < monthAbbrLen {
		return 0
	}
	month = month[:monthAbbrLen]
	switch month {
	case "jan":
		return 1
	case "feb":
		return 2
	case "mar":
		return 3
	case "apr":
		return 4
	case "may":
		return 5
	case "jun":
		return 6
	case "jul":
		return 7
	case "aug":
		return 8
	case "sep":
		return 9
	case "oct":
		return 10
	case "nov":
		return 11
	case "dec":
		return 12
	}
	return 0
}

func parseHumanNumeric(s string) (float64, error) {
	if len(s) == 0 {
		return 0, errors.New("empty string")
	}
	lastByte := s[len(s)-1]
	noSuffix := unicode.IsDigit(rune(lastByte))
	if noSuffix {
		return strconv.ParseFloat(s, 64)
	}

	num, err := strconv.ParseFloat(s[:len(s)-1], 64)
	if err != nil {
		return num, err
	}

	suffix := unicode.ToLower(rune(lastByte))
	power := 0
	switch suffix {
	case 'k':
		power = 3
	case 'm':
		power = 6
	case 'g':
		power = 9
	case 't':
		power = 12
	case 'p':
		power = 15
	case 'e':
		power = 18
	case 'z':
		power = 21
	case 'y':
		power = 24
	}

	return num * math.Pow(10, float64(power)), nil
}
