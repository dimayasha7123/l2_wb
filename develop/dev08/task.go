package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*
=== Взаимодействие с ОС ===

Необходимо реализовать свой собственный UNIX-шелл-утилиту с поддержкой ряда простейших команд:

- cd <args> - смена директории (в качестве аргумента могут быть то-то и то)
- pwd - показать путь до текущего каталога
- echo <args> - вывод аргумента в STDOUT
- kill <args> - "убить" процесс, переданный в качестве аргумента (пример: такой-то пример)
- ps - выводит общую информацию по запущенным процессам в формате *такой-то формат*

Так же требуется поддерживать функционал fork/exec-команд

Дополнительно необходимо поддерживать конвейер на пайпах (linux pipes, пример cmd1 | cmd2 | .... | cmdN).

*Шелл — это обычная консольная программа, которая будучи запущенной, в интерактивном сеансе выводит некое приглашение
в STDOUT и ожидает ввода пользователя через STDIN. Дождавшись ввода, обрабатывает команду согласно своей логике
и при необходимости выводит результат на экран. Интерактивный сеанс поддерживается до тех пор, пока не будет введена
команда выхода (например \quit).

*/

func main() {
	reader := bufio.NewReader(os.Stdin)

	commands := map[string]commandType{
		"pwd":  pwd,
		"cd":   cd,
		"echo": echo,
		"ls":   ls,
		"ps":   ps,
		"kill": kill,
	}

	for {
		fmt.Printf("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		input = input[:len(input)-1]
		if len(input) == 0 {
			continue
		}

		split := strings.Split(input, " ")
		if len(split) == 0 {
			continue
		}
		commandName := split[0]
		args := split[1:]

		if commandName == "quit" {
			return
		}

		command, ok := commands[commandName]
		if !ok {
			fmt.Fprintf(os.Stderr, "Command %q not found\n", commandName)
			continue
		}

		outBytes, errBytes := command(args)
		if outBytes != nil {
			fmt.Printf("%s\n", string(outBytes))
		}
		if errBytes != nil {
			fmt.Fprintf(os.Stderr, "%s\n", string(errBytes))
		}
	}

}

type commandType func(args []string) ([]byte, []byte)

func pwd(args []string) ([]byte, []byte) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, []byte(err.Error())
	}
	return []byte(wd), nil
}

func cd(args []string) ([]byte, []byte) {
	if len(args) != 1 {
		return nil, []byte(fmt.Sprintf("must be exact 1 argument, but get %d", len(args)))
	}
	err := os.Chdir(args[0])
	if err != nil {
		return nil, []byte(err.Error())
	}
	return nil, nil
}

func echo(args []string) ([]byte, []byte) {
	return []byte(strings.Join(args, " ")), nil
}

func ls(args []string) ([]byte, []byte) {
	if len(args) > 1 {
		return nil, []byte(fmt.Sprintf("must be 0 or 1 argument, but get %d", len(args)))
	}
	arg := "."
	if len(args) == 1 {
		arg = args[0]
	}
	dirs, err := os.ReadDir(arg)
	if err != nil {
		return nil, []byte(err.Error())
	}
	entries := make([]string, 0, len(dirs))
	for _, dir := range dirs {
		entries = append(entries, dir.Name())
	}

	return []byte(strings.Join(entries, " ")), nil
}

func ps(args []string) ([]byte, []byte) {
	files, err := os.ReadDir("/proc")
	if err != nil {
		return nil, []byte("can't read processes file")
	}

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%-5s %-10s %-10s\n", "PID", "TTY", "CMD"))
	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		pid, err := strconv.Atoi(file.Name())
		if err != nil {
			continue
		}
		tty, err := getTTY(pid)
		if err != nil {
			continue
		}
		cmd, err := getCMD(pid)
		if err != nil {
			continue
		}
		sb.WriteString(fmt.Sprintf("%-5d %-10s %-10s\n", pid, tty, cmd))
	}

	return []byte(sb.String()), nil
}

func getTTY(pid int) (string, error) {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return "", fmt.Errorf("can't read file: %v", err)
	}
	fields := strings.Fields(string(data))
	if len(fields) < 7 {
		return "", fmt.Errorf("not enought fields")
	}
	ttyNum, err := strconv.Atoi(fields[6])
	if err != nil {
		return "", fmt.Errorf("can't convert tty num")
	}
	if ttyNum == 0 {
		return "?", nil
	}
	return fmt.Sprintf("/dev/tty%d", ttyNum), nil
}

func getCMD(pid int) (string, error) {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
	if err != nil {
		return "", fmt.Errorf("can't read file: %v", err)
	}
	return strings.ReplaceAll(string(data), "\x00", " "), nil
}

func kill(args []string) ([]byte, []byte) {
	if len(args) != 1 {
		return nil, []byte(fmt.Sprintf("must be exact 1 argument, but get %d", len(args)))
	}

	pid, err := strconv.Atoi(os.Args[1])
	if err != nil {
		return nil, []byte("pid not valid")
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return nil, []byte("process not found")
	}

	err = process.Kill()
	if err != nil {
		return nil, []byte(fmt.Sprint("can't kill process:", err))
	}

	return []byte(fmt.Sprintf("process %d killed", pid)), nil
}
