package main

import "fmt"

/*
	Реализовать паттерн «комманда».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Command_pattern
*/

/*
	ПРИМЕНИМОСТЬ:
1.	Параметризация объектов выполняемым действием
2.	Необходимость поддержки очереди запросов (создание, добавление, удаление, логгирование и т.д.)


	ПЛЮСЫ:
1.	Отделение объекта, инициирующего операцию от объекта, который будет ее исполнять
2.	Команды - это объекты, которые можно удобно обрабатывать
3.	Возможность составления более сложных команд из простых
4.	Простота добавления новых команд

	МИНУСЫ:
1. 	Сложность реализации. Реализация паттерна команда может быть сложной,
	особенно если необходимо поддержать отмену и повторение операций.
2. 	Недостаточная гибкость. При необходимости изменения поведения команды, может
	потребоваться изменить много других классов.

	ПРИМЕРЫ:
1.	Реализация одних и тех же действий из разных элементов UI (сохранение документа по горячим клавишам,
	при выходе в диалоговом окне, либо же по нажатию кнопке в основном окне)
2.	Поддержка отмены действий по Ctrl+Z (каждая команда представляет собой отдельную операцию,
	которая может быть отменена или повторена)
*/

type button struct {
	command command
}

func (b *button) press() {
	b.command.execute()
}

type command interface {
	execute()
}

type offCommand struct {
	device device
}

func (c *offCommand) execute() {
	c.device.off()
}

type onCommand struct {
	device device
}

func (c *onCommand) execute() {
	c.device.on()
}

type device interface {
	on()
	off()
}

type tv struct {
	isRunning bool
}

func (t *tv) on() {
	t.isRunning = true
	fmt.Println("Turning tv on")
}

func (t *tv) off() {
	t.isRunning = false
	fmt.Println("Turning tv off")
}

func main() {
	tv := &tv{}
	onCommand := &onCommand{
		device: tv,
	}
	offCommand := &offCommand{
		device: tv,
	}
	onButton := &button{
		command: onCommand,
	}
	onButton.press()
	offButton := &button{
		command: offCommand,
	}
	offButton.press()
}
