package converters

import (
	"fmt"
	"time"
)

// StrToDate function
func StrToDate(str string) (time.Time, error) {
	date, err := time.Parse(time.DateOnly, str)
	if err != nil {
		return date, fmt.Errorf("can't parse date: %v", err)
	}
	return date, nil
}

// DateToStr function
func DateToStr(date time.Time) string {
	return date.Format(time.DateOnly)
}
