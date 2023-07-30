package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"l2_wb/develop/dev05/greputil"
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

func main() {
	var set greputil.Settings
	flag.IntVar(&set.A, "A", set.A, "\"after\" печатать +N строк после совпадения")
	flag.IntVar(&set.B, "B", set.B, "\"before\" печатать +N строк до совпадения")
	flag.IntVar(&set.C, "C", set.C, "\"context\" (A+B) печатать ±N строк вокруг совпадения")
	flag.IntVar(&set.Count, "c", set.Count, "\"count\" (количество строк)")
	flag.BoolVar(&set.IgnoreCase, "i", set.IgnoreCase, "\"ignore-case\" (игнорировать регистр)")
	flag.BoolVar(&set.Invert, "v", set.Invert, "\"invert\" (вместо совпадения, исключать)")
	flag.BoolVar(&set.Fixed, "F", set.Fixed, "\"fixed\", точное совпадение со строкой, не паттерн")
	flag.BoolVar(&set.LineNum, "n", set.LineNum, "\"line num\", печатать номер строки")
	flag.BoolVar(&set.FileName, "f", set.FileName, "\"file name\", выводить имя файла")
	flag.BoolVar(&set.Highlighting, "h", set.Highlighting, "\"highlighting\", выделять совпадения")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Использование grep:\n")
		fmt.Fprintf(os.Stderr, "grep [ФЛАГ]... ШАБЛОН [ФАЙЛ]...\n")
		fmt.Fprintf(os.Stderr, "Если в качестве первого файла указан \"-\", то будет считан STDIN\n")
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
			fmt.Fprintf(os.Stderr, "Не могу считать STDIN: %v\n", err)
			os.Exit(1)
		}

		rows, err := greputil.GrepBytes(data, pattern, set)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Не могу грепнуть: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(rows)
		os.Exit(0)
	}

	for _, file := range files {
		if set.FileName {
			fmt.Printf("Файл %q:\n", file)
		}

		data, err := os.ReadFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Не могу прочитать файл: %v\n", err)
			continue
		}

		rows, err := greputil.GrepBytes(data, pattern, set)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Не могу грепнуть: %v\n", err)
			continue
		}
		fmt.Print(rows)
	}
}
