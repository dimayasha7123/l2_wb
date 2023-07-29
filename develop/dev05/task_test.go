package main

import (
	"fmt"
	"os/exec"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
// работа хайлайта (2 тест-кейса, сами совпадения и номера строк)
// работа флага вывода файла (один файл и несколько файлов со включенной или выключенной опцией)

const directory = "test_data"

func Test_e2e_grep(t *testing.T) {
	grepPath := path.Join(t.TempDir(), "task")
	compileErr := exec.Command("go1.20.1", "build", "-o", grepPath, "task.go").Run()
	require.NoError(t, compileErr, "can't compile grep: %v\n", compileErr)

	cmd := exec.Command(grepPath, "o.", "-")
	cmd.Stdin = strings.NewReader("London")

	output := strings.Builder{}
	cmd.Stdout = &output
	errOutput := strings.Builder{}
	cmd.Stderr = &errOutput

	err := cmd.Run()
	//require.NoError(t, err)

	fmt.Println("ERR:", err)
	fmt.Println("STDOUT:", []byte(output.String()))
	fmt.Println("STDOUT:", output.String())
	for _, r := range []rune(output.String()) {
		fmt.Printf("%c ", r)
	}
	fmt.Println()
	fmt.Println("STDERR:", errOutput.String())

	tests := []struct {
		name       string
		args       []string
		stdin      string
		wantStdout string
		wantStderr string
		errFunc    assert.ErrorAssertionFunc
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}
