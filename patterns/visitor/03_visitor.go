package main

import (
	"errors"
	"fmt"
)

/*
	Реализовать паттерн «посетитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Visitor_pattern
*/

/*
	ПРИМЕНИМОСТЬ:
0.	Необходимость определения новой(ых) операции(ий) с подклассами без значительных их изменений
1.	Наличие множества подклассов с различными интерфейсами, с которыми необходимо выполнять схожие операции,
	зависящие от этих подклассов
2.	В дополнение к предыдущему пункту. Иногда не хочется засорять подклассы новыми методами (из-за общей структуры объектов,
	необходимости поддержки или высокой частоты изменения/добавления новых операций)

	ПЛЮСЫ:
1.	Возможность добавления новых методов без значительных изменений подклассов
2.	Управление группой подклассов в одном месте
3.	Удобная обработка списков/иерархий из подклассов

	МИНУСЫ:
1.	Нарушение инкапсуляции объекта, если необходим доступ к приватным полям и методам
2.	Сложность поддержки при частом изменении структур подклассов
3.	Необходимость реализации метода visit(visitor) у подклассов (но только один раз)

	ПРИМЕРЫ:
1.	Обход различных иерархий объектов/древовидных структур с различными типами нод, например, иерархия объектов в библиотеках
	для работы с графикой
2.	Обработка различных событий в событийно-ориентированных приложениях (клики, ошибки, запросы пользователя)
*/

type bikePartsVisitor interface {
	visitWheel(wheel)
	visitBrakes(brakes)
	visitTransmission(transmission)
}

type bikeChecker struct {
	errors []error
}

func newBikeChecker() *bikeChecker {
	return &bikeChecker{errors: make([]error, 0, 5)}
}

func (b *bikeChecker) Check() error {
	return errors.Join(b.errors...)
}

func (b *bikeChecker) visitWheel(w wheel) {
	if w.CheckPressure() < 2 {
		b.errors = append(b.errors, errors.New("low tire pressure"))
	}
}

func (b *bikeChecker) visitBrakes(br brakes) {
	if br.Creak() {
		b.errors = append(b.errors, errors.New("break creak"))
	}
	if br.TestBreakPath() > 10 {
		b.errors = append(b.errors, errors.New("long brake path"))
	}
}

func (b *bikeChecker) visitTransmission(t transmission) {
	if !t.Greased() {
		b.errors = append(b.errors, errors.New("chain not greased"))
	}
	b.errors = append(b.errors, t.TestSwitch())
}

type visitorAcceptor interface {
	accept(bikePartsVisitor)
}

type wheel struct {
	// some implementation
}

func (w wheel) accept(v bikePartsVisitor) {
	v.visitWheel(w)
}

func (w wheel) CheckPressure() int {
	// some implementation
	return 2
}

type brakes struct {
	// some implementation
}

func (b brakes) accept(v bikePartsVisitor) {
	v.visitBrakes(b)
}

func (b brakes) Creak() bool {
	// some implementation
	return true
}

func (b brakes) TestBreakPath() int {
	// some implementation
	return 4
}

type transmission struct {
	// some implementation
}

func (t transmission) accept(v bikePartsVisitor) {
	v.visitTransmission(t)
}

func (t transmission) Greased() bool {
	// some implementation
	return false
}

func (t transmission) TestSwitch() error {
	// some implementation
	return nil
}

func main() {
	parts := []visitorAcceptor{brakes{}, transmission{}, wheel{}}
	checker := newBikeChecker()

	for _, part := range parts {
		part.accept(checker)
	}

	fmt.Println(checker.Check())
}
