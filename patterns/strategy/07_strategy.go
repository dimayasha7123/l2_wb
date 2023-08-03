package main

import "fmt"

/*
	Реализовать паттерн «стратегия».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Strategy_pattern
*/

/*
	ПРИМЕНИМОСТЬ:
1.	Наличие нескольких разновидностей алгоритма, например, эффективный и "в лоб"
2.	Наличие множества классов, отличающихся только поведением. Стратегия позволяет настроить это поведение.
3.	Инкапсуляция алгоритма внутри стратегии и сокрытие его от классов.
4. 	Куча ифов, определяющих различное поведение. Можно вынести в стратегии.

	ПЛЮСЫ:
1.	Определение семейства схожих алгоритмов, которые можно использовать с различными классами.
2.	Избавление от ифов.
3. 	Возможность выбора реализации.

	МИНУСЫ:
1.	Сложность выбора подходящей стратегии со стороны клиента.
2.	Увеличение числа объектов

	ПРИМЕРЫ:
1.	Выбор алгоритма хэширования (Least Recently Used, First In First Out, First In First Out)
2.	Выбор алгоритма поиска по тексту (паттерн, строгое соответствие, с возможной инверсией выбора и т.д.)
3.	Создание парсера выражений (реализация различных бинарных операторов)
*/

type operator interface {
	apply(int, int) int
}

type operation struct {
	operator operator
}

func (o *operation) operate(leftValue, rightValue int) int {
	return o.operator.apply(leftValue, rightValue)
}

type addition struct{}

func (addition) apply(lval, rval int) int {
	return lval + rval
}

type multiplication struct{}

func (multiplication) apply(lval, rval int) int {
	return lval * rval
}

func main() {
	add := operation{addition{}}
	fmt.Println(add.operate(3, 5))

	mult := operation{multiplication{}}
	fmt.Println(mult.operate(3, 5))
}
