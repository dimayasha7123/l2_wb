package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
)

/*
=== Утилита grep ===

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type settings struct {
	A          int
	B          int
	C          int  // Перекрывает A и B, если не равно нулю +
	Count      int  // При Invert Count будет считать количество не подошедших строк
	IgnoreCase bool // Не будет работать с fixed +
	Invert     bool // При Invert A, B, C не учитываются
	Fixed      bool // Изменяет поиск по regex на поиск подстрок +
	LineNum    bool
}

// ТЕСТ КЕЙСЫ:
// когда строка в самом начале
// когда между строкой и началом меньше строк чем В
// когда между строкой и началом строк ровно столько, сколько В
// когда между строкой и началом больше строк, чем В
// предыдущие 4 пункта, но между двумя подходящими строками
// предыдущие 5 пунктов, но с В и концом, вместо начала
// перекрытие А и В (как первые 4 кейса)
// С перекрывает А и В если не равно нулю
// при инверт А, Б и С не учитываются
// каунт считает количество подошедших строк
// каунт считает количество неподошедших строк при инверт
// каунт не считает строки, окружающие ответ
// игнор_кейс базовый
// игнор_кейс не работает с фиксед
// работа фикседа

func (s settings) Validate() error {
	errs := make([]error, 0, 4)
	if s.A < 0 {
		errs = append(errs, errors.New("A не может быть меньше нуля"))
	}
	if s.B < 0 {
		errs = append(errs, errors.New("B не может быть меньше нуля"))
	}
	if s.C < 0 {
		errs = append(errs, errors.New("C не может быть меньше нуля"))
	}
	if s.Count < 0 {
		errs = append(errs, errors.New("COUNT не может быть меньше нуля"))
	}
	return errors.Join(errs...)
}

type grepRows []grepRow

type grepRow struct {
	LineNum      int
	Data         []byte
	MatchIndexes [][]int
}

func grep(data [][]byte, pattern string, set settings) (grepRows, error) {
	after := set.A
	before := set.B
	if set.C > 0 {
		after = set.C
		before = set.C
	}

	if set.IgnoreCase && !set.Fixed {
		pattern = ignoreCasePrefix + pattern
	}
	if set.Fixed {
		pattern = regexp.QuoteMeta(pattern)
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("не могу скомпилировать регулярное выражение %s: %v", pattern, err)
	}

	ret := make(grepRows, 0, len(data))
	count := 0
	for i, dataRow := range data {
		if set.Count != 0 && count >= set.Count {
			break
		}
		indexPairs := re.FindAllIndex(dataRow, -1)
		if set.Invert == (len(indexPairs) == 0) {
			row := grepRow{
				LineNum:      i + 1,
				Data:         dataRow,
				MatchIndexes: indexPairs,
			}
			ret = append(ret, row)
			count++
		}
	}

	if !set.Invert && (after != 0 || before != 0) {
		last := 0
		newRet := make(grepRows, 0, len(data))
		for i, row := range ret {
			j := row.LineNum - 1
			leftBorder := max(last, j-before)
			rightBorder := min(j+after+1, len(data))
			if i+1 < len(ret) {
				rightBorder = min(rightBorder, ret[i+1].LineNum-1)
			}
			for k := leftBorder; k < rightBorder; k++ {
				newRow := grepRow{
					LineNum: k + 1,
					Data:    data[k],
				}
				if k == j {
					newRow.MatchIndexes = row.MatchIndexes
				}
				newRet = append(newRet, newRow)
			}
			last = rightBorder
		}
		ret = newRet
	}

	return ret, nil
}

func main() {
	var set settings
	flag.IntVar(&set.A, "A", set.A, "\"after\" печатать +N строк после совпадения")
	flag.IntVar(&set.B, "B", set.B, "\"before\" печатать +N строк до совпадения")
	flag.IntVar(&set.C, "C", set.C, "\"context\" (A+B) печатать ±N строк вокруг совпадения")
	flag.IntVar(&set.Count, "c", set.Count, "\"count\" (количество строк)")
	flag.BoolVar(&set.IgnoreCase, "i", set.IgnoreCase, "\"ignore-case\" (игнорировать регистр)")
	flag.BoolVar(&set.Invert, "v", set.Invert, "\"invert\" (вместо совпадения, исключать)")
	flag.BoolVar(&set.Fixed, "F", set.Fixed, "\"fixed\", точное совпадение со строкой, не паттерн")
	flag.BoolVar(&set.LineNum, "n", set.LineNum, "\"line num\", печатать номер строки")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Использование grep:\n")
		fmt.Fprintf(os.Stderr, "grep [ФЛАГ]... ШАБЛОН [ФАЙЛ]...\n")
		fmt.Fprintf(os.Stderr, "Если в качестве первого файла указан \"-\", то будет считан STDIN")
		flag.PrintDefaults()
	}
	flag.Parse()

	pattern := flag.Arg(0)
	if len(pattern) == 0 {
		fmt.Fprintf(os.Stderr, "ШАБЛОН не может быть пустым\n")
		os.Exit(2)
	}

	if len(flag.Args()) < 2 {
		fmt.Fprintf(os.Stderr, "Не указаны файлы или \"-\" для чтения из STDIN\n")
		os.Exit(2)
	}

	err := set.Validate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Не верно указаны флаги: %v\n", err)
		os.Exit(2)
	}

	files := flag.Args()[1:]
	if files[0] == "-" {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Не могу считать из STDIN: %v\n", err)
			os.Exit(1)
		}

		rows, err := grep(bytes.Split(data, []byte{byte(stringSeparator)}), pattern, set)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Не могу грепнуть: %v\n", err)
			os.Exit(1)
		}
		fmt.Println()
		printRows(rows, set)
		os.Exit(0)
	}

	for _, file := range files {
		fmt.Printf("Файл %q:\n", file)

		data, err := os.ReadFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Не могу прочитать файл: %v\n", err)
			continue
		}

		rows, err := grep(bytes.Split(data, []byte{byte(stringSeparator)}), pattern, set)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Не могу грепнуть: %v\n", err)
			continue
		}
		printRows(rows, set)
	}
}

func printRows(rows grepRows, set settings) {
	var maxNumLen int
	for _, row := range rows {
		maxNumLen = max(maxNumLen, len(fmt.Sprint(row.LineNum)))
	}

	lineNumPattern := fmt.Sprintf("%%%dd: ", maxNumLen)
	for _, row := range rows {
		if set.LineNum {
			setBlueColor()
			fmt.Printf(lineNumPattern, row.LineNum)
			setDefaultColor()
		}
		i := 0
		for _, indexes := range row.MatchIndexes {
			left := indexes[0]
			right := indexes[1]

			fmt.Print(string(row.Data[i:left]))
			setRedColor()
			fmt.Print(string(row.Data[left:right]))
			setDefaultColor()

			i = right
		}
		fmt.Println(string(row.Data[i:]))
	}
}

const (
	colorReset       = "\033[0m"
	colorRed         = "\033[31m"
	colorBlue        = "\033[34m"
	stringSeparator  = '\n'
	ignoreCasePrefix = "(?i)"
)

func setDefaultColor() {
	fmt.Print(colorReset)
}

func setRedColor() {
	fmt.Print(colorRed)
}

func setBlueColor() {
	fmt.Print(colorBlue)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
