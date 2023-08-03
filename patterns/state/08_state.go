package main

import (
	"errors"
	"fmt"
	"os"
)

/*
	Реализовать паттерн «состояние».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/State_pattern
*/

/*
	ПРИМЕНИМОСТЬ:
1.	Поведение объекта зависит от его состояния и должно изменяться в рантайме
2.	Большое количество условных операторов и свич-кейсов, которые обрабатывают все состояния, паттерн предлагает
	поместить каждый ифчик в отдельный класс

	ПЛЮСЫ:
1.	Инкапсуляция всего поведения, связанного с конкретным состоянием в отдельный объект
2.	Явно выраженные переходы между состояниями

	МИНУСЫ:
1.	Может привести к созданию большого количества классов состояний
2.	Оверхед для простых случаев

	ПРИМЕРЫ:
1.	Реализация сетевых протоколов, поддерживающих соединение (TCP)
2.	Интерактивные программы (часто с графическим интерфейсом)
3.	Сетевая сессия пользователя в некоторых случаях будет поддерживать состояние (например, ожидание оплаты товара)
*/

func main() {
	cm := newCoffeeMachine(3)
	fmt.Println(cm)

	var err error
	for {
		err = cm.DoCoffee()
		if err != nil {
			fmt.Println("can't do coffee:", err)
			break
		}
		fmt.Println(cm)
	}

	err = cm.AddIngredient(1)
	if err != nil {
		fmt.Println("can't add ingredients:", err)
		os.Exit(1)
	}
	fmt.Println(cm)

	err = cm.DoCoffee()
	if err != nil {
		fmt.Println("can't do coffee:", err)
		os.Exit(1)
	}
	fmt.Println(cm)

	err = cm.DoCoffee()
	if err != nil {
		fmt.Println("can't do coffee:", err)
		os.Exit(1)
	}
	fmt.Println(cm)

}

type coffeeMachineState interface {
	AddIngredient(ingredients int) error
	DoCoffee() error
}

type coffeeMachine struct {
	needIngredients coffeeMachineState
	waiting         coffeeMachineState

	state coffeeMachineState

	ingredientsCount int
	maxIngredients   int
	doneCoffees      int
}

func (m *coffeeMachine) String() string {
	return fmt.Sprintf("%d/%d ingr., done %d coffees", m.ingredientsCount, m.maxIngredients, m.doneCoffees)
}

func newCoffeeMachine(maxIngredients int) *coffeeMachine {
	m := &coffeeMachine{
		ingredientsCount: maxIngredients,
		maxIngredients:   maxIngredients,
	}
	nis := needIngredientsState{m}
	m.needIngredients = &nis

	ws := waitingState{m}
	m.waiting = &ws

	m.setState(m.waiting)
	return m
}

func (m *coffeeMachine) AddIngredient(ingredients int) error {
	return m.state.AddIngredient(ingredients)
}

func (m *coffeeMachine) DoCoffee() error {
	return m.state.DoCoffee()
}

func (m *coffeeMachine) setState(state coffeeMachineState) {
	m.state = state
}

type needIngredientsState struct {
	cm *coffeeMachine
}

func (n needIngredientsState) AddIngredient(ingredients int) error {
	if ingredients <= 0 {
		return errBadIngredientCount
	}
	n.cm.ingredientsCount = min(n.cm.maxIngredients, n.cm.ingredientsCount+ingredients)
	n.cm.setState(n.cm.waiting)

	return nil
}

func (n needIngredientsState) DoCoffee() error {
	return errNotEnoughIngredients
}

type waitingState struct {
	cm *coffeeMachine
}

func (w waitingState) AddIngredient(ingredients int) error {
	if ingredients <= 0 {
		return errBadIngredientCount
	}
	w.cm.ingredientsCount = min(w.cm.maxIngredients, w.cm.ingredientsCount+ingredients)

	return nil
}

func (w waitingState) DoCoffee() error {
	fmt.Println("Done coffee!")
	w.cm.doneCoffees++
	w.cm.ingredientsCount--

	if w.cm.ingredientsCount <= 0 {
		w.cm.setState(w.cm.needIngredients)
	}
	return nil
}

var (
	errBadIngredientCount   = errors.New("zero or less ingredient count")
	errNotEnoughIngredients = errors.New("not enough ingredients")
)

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
