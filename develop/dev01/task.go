package main

import (
	"fmt"
	"os"
	"time"

	"github.com/beevik/ntp"
)

/*
=== Базовая задача ===

Создать программу, печатающую точное время с использованием NTP библиотеки. Инициализировать как go module.
Использовать библиотеку https://github.com/beevik/ntp.
Написать программу, печатающую текущее время / точное время с использованием этой библиотеки.

Программа должна быть оформлена с использованием как go module.
Программа должна корректно обрабатывать ошибки библиотеки: распечатывать их в STDERR и возвращать ненулевой код выхода в OS.
Программа должна проходить проверки go vet и golint.
*/

const ntpAddress = "ru.pool.ntp.org"

// GetPreciseTime return precise current time
func GetPreciseTime() (time.Time, error) {
	return ntp.Time(ntpAddress)
}

func main() {
	preciseTime, err := GetPreciseTime()
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't get precise time: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(preciseTime.Format("2 Jan 2006 15:04:05.999999"))
}
