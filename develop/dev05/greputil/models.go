package greputil

import (
	"errors"
	"fmt"
	"strings"

	"l2_wb/develop/dev05/color"
)

// Settings used for set grep settings
type Settings struct {
	A            int
	B            int
	C            int  // Перекрывает A и B, если не равно нулю +
	Count        int  // При Invert Count будет считать количество не подошедших строк
	IgnoreCase   bool // Не будет работать с fixed +
	Invert       bool // При Invert A, B, C не учитываются
	Fixed        bool // Изменяет поиск по regex на поиск подстрок +
	LineNum      bool
	FileName     bool
	Highlighting bool
}

// Validate used to validate settings
func (s Settings) Validate() error {
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

const (
	stringSeparator  = '\n'
	ignoreCasePrefix = "(?i)"
)

// GrepRows is GrepBytes or GrepSplitBytes result
type GrepRows struct {
	set  Settings
	rows []grepRow
}

type grepRow struct {
	LineNum      int
	Data         []byte
	MatchIndexes [][]int
}

func (gr GrepRows) String() string {
	var maxNumLen int
	for _, row := range gr.rows {
		maxNumLen = max(maxNumLen, len(fmt.Sprint(row.LineNum)))
	}

	sb := &strings.Builder{}
	lineNumPattern := fmt.Sprintf("%%%dd: ", maxNumLen)
	cw := color.NewWrapper(gr.set.Highlighting)
	for _, row := range gr.rows {
		if gr.set.LineNum {
			cw.SetBlue(sb)
			fmt.Fprintf(sb, lineNumPattern, row.LineNum)
			cw.Reset(sb)
		}
		i := 0
		for _, indexes := range row.MatchIndexes {
			left := indexes[0]
			right := indexes[1]

			fmt.Fprint(sb, string(row.Data[i:left]))
			cw.SetRed(sb)
			fmt.Fprint(sb, string(row.Data[left:right]))
			cw.Reset(sb)

			i = right
		}
		fmt.Fprintln(sb, string(row.Data[i:]))
	}

	return sb.String()
}
