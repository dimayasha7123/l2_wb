package main

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
)

/*
	Реализовать паттерн «фабричный метод».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Factory_method_pattern
*/

/*
	ПРИМЕНИМОСТЬ:
1.	Необходимость выбора реализации (подкласса) на этапе выполнения программы.

	ПЛЮСЫ:
1.	Создание объектов различных типов без больших изменений существующего кода.
2.	Улучшение и упрощение тестирования, т.к. присутствует возможность создания заглушек.
3.	Инкапсуляция логики создания различных объектов в отдельном классе (функции)

	МИНУСЫ:
1.	Необходимость поддержки новых подклассов в классе фабрики

	ПРИМЕРЫ:
1.	Различные конвертеры и парсеры (xml парсер в java), конструкторы объектов в древовидных структурах (DOM model in HTML5)
2.	Создание объектов для работы с различными базами данных в Go
*/

type weapon interface {
	Name() string
	QuickAttack() int
	StrongAttack() int
}

type sword struct {
	bladeLength int
	twoHanded   bool
}

func newSword(bladeLength int, twoHanded bool) sword {
	return sword{
		bladeLength: bladeLength,
		twoHanded:   twoHanded,
	}
}

func (s sword) Name() string {
	hand := "one"
	if s.twoHanded {
		hand = "two"
	}
	return fmt.Sprintf("Sword with blade length = %d sm., %s handed", s.bladeLength, hand)
}

func (s sword) QuickAttack() int {
	ret := 4
	if s.twoHanded {
		ret--
	}
	ret -= int(math.Round(math.Abs(float64(60-s.bladeLength)) / 10))
	ret += dice(6, 1)
	return ret
}

func (s sword) StrongAttack() int {
	ret := 6
	if s.twoHanded {
		ret += 2
	}
	ret -= int(math.Round(math.Abs(float64(70-s.bladeLength)) / 8))
	ret += dice(10, 1)
	return ret
}

func dice(faces int, count int) int {
	ret := 0
	for i := 0; i < count; i++ {
		ret += rand.Intn(faces) + 1
	}
	return ret
}

type rifle struct {
	caliber       int
	rof           int
	accuracy      int
	bulletNumbers int
}

func newRifle(caliber, rof, accuracy, bulletNumbers int) rifle {
	return rifle{
		caliber:       caliber,
		rof:           rof,
		accuracy:      accuracy,
		bulletNumbers: bulletNumbers,
	}
}

func (r rifle) Name() string {
	return fmt.Sprintf(
		"Rifle with %dmm caliber bullets, ROF = %d, ACC = %d, has %d bullets",
		r.caliber,
		r.rof,
		r.accuracy,
		r.bulletNumbers,
	)
}

func (r rifle) QuickAttack() int {
	ret := 0
	for i := 0; i < r.rof; i++ {
		if dice(6, 1)+r.accuracy <= 3 {
			continue
		}
		ret += int(float64(r.accuracy)/10 + float64(dice(4, 1)))
	}
	return ret
}

func (r rifle) StrongAttack() int {
	ret := 0
	for i := 0; i < r.bulletNumbers; i++ {
		if dice(6, 1)+r.accuracy <= 4 {
			continue
		}
		ret += int(float64(r.accuracy)/10 + float64(dice(4, 1)))
	}
	return ret
}

type weaponType int

const (
	swordType weaponType = iota
	rifleType
)

func getWeapon(wType weaponType) (weapon, error) {
	switch wType {
	case swordType:
		return newSword(65, false), nil
	case rifleType:
		return newRifle(40, 20, 2, 30), nil
	}
	return nil, errors.New("no such type")
}

func main() {
	r, _ := getWeapon(rifleType)
	s, _ := getWeapon(swordType)

	testWeapon(r)
	testWeapon(s)
}

func testWeapon(w weapon) {
	fmt.Printf("Testing: %q\n", w.Name())
	fmt.Printf("Quick attack: %d\n", w.QuickAttack())
	fmt.Printf("Strong attack: %d\n", w.StrongAttack())
}
