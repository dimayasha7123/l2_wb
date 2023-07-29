package greputil

import (
	"bytes"
	"fmt"
	"regexp"
)

// GrepBytes used to grep []byte by pattern with some settings
func GrepBytes(data []byte, pattern string, set Settings) (GrepRows, error) {
	return GrepSplitBytes(bytes.Split(data, []byte{byte(stringSeparator)}), pattern, set)
}

// GrepSplitBytes used to grep split data [][]byte by pattern with some settings
func GrepSplitBytes(data [][]byte, pattern string, set Settings) (GrepRows, error) {
	err := set.Validate()
	if err != nil {
		return GrepRows{}, fmt.Errorf("set не валидны: %v", err)
	}

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
		return GrepRows{}, fmt.Errorf("не могу скомпилировать регулярное выражение %s: %v", pattern, err)
	}

	rows := make([]grepRow, 0, len(data))
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
			rows = append(rows, row)
			count++
		}
	}

	if !set.Invert && (after != 0 || before != 0) {
		last := 0
		newRet := make([]grepRow, 0, len(data))
		for i, row := range rows {
			j := row.LineNum - 1
			leftBorder := max(last, j-before)
			rightBorder := min(j+after+1, len(data))
			if i+1 < len(rows) {
				rightBorder = min(rightBorder, rows[i+1].LineNum-1)
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
		rows = newRet
	}

	return GrepRows{
		set:  set,
		rows: rows,
	}, nil
}
