package main

import (
	"fmt"
	"strings"
)

/*
	Реализовать паттерн «строитель».
	Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Builder_pattern
*/

/*
	ПРИМЕНИМОСТЬ:
0.	Разделение процесса конструирования сложного объекта от его представления (один и тот же процесс конструирования может
	создать различные объекта)
1.	В результате конструирования получаются различные представления необходимого объекта
2.	Алгоритм конструирования должен быть простым и не зависеть от типа конструируемых объектов и способа стыковки между ними

	ПЛЮСЫ:
1.	Создание различных представлений в результате одного и того же процесса конструирования
2.	Предоставляет контроль над этапами конструирования объектов
3.	Инкапсуляция логики создания сложного объекта (создание подобъектов и их стыковка)

	МИНУСЫ:
1.	Усложняет использование объекта (нужен руководитель и билдер)
2.	Под каждое представление необходимо реализовать свой билдер (много кода)
3.	Не эффективен для создания простых объектов

	ПРИМЕРЫ:
1.	Различные конвертеры (из RTF в простой текст или виджет на рабочем столе, или же TeX)
2.	Использование различных реализаций интерфейсов ввода-вывода в Go. Например, можно передавать
	как strings.Builder, так и bytes.Buffer в функцию, принимающую io.Writer
*/

type item struct {
	Title string
	Cost  int
}

type boss struct {
	Name string
	// another fields
	Difficulty int
}

type mapCreatorI interface {
	CreateDungeon(size, enemies int)
	CreateCheckpoint(healPoints int)
	CreateTreasureRoom(items []item)
	CreateBossArena(boss boss)
}

type director struct {
	creator mapCreatorI
}

func newDirector(creator mapCreatorI) *director {
	return &director{creator: creator}
}

func (d *director) CreateFirstLevel() {
	d.creator.CreateDungeon(3, 2)
	d.creator.CreateDungeon(4, 5)

	d.creator.CreateCheckpoint(2)

	d.creator.CreateDungeon(4, 3)
	d.creator.CreateDungeon(5, 6)
	d.creator.CreateDungeon(4, 4)

	d.creator.CreateTreasureRoom([]item{{Title: "big sword", Cost: 4}, {Title: "elf arch", Cost: 5}})

	d.creator.CreateDungeon(3, 2)
	d.creator.CreateDungeon(4, 5)

	d.creator.CreateCheckpoint(4)

	d.creator.CreateBossArena(boss{Name: "viverna", Difficulty: 4})
}

type mapBuilder struct {
	rooms []string
}

func newMapBuilder() *mapBuilder {
	return &mapBuilder{rooms: make([]string, 0)}
}

func (m *mapBuilder) CreateDungeon(size, enemies int) {
	m.rooms = append(m.rooms, fmt.Sprintf("dungeon %dx%d meters with %d enemies", size, size, enemies))
}

func (m *mapBuilder) CreateCheckpoint(healPoints int) {
	m.rooms = append(m.rooms, fmt.Sprintf("checkpoint with %d heal points", healPoints))

}

func (m *mapBuilder) CreateTreasureRoom(items []item) {
	m.rooms = append(m.rooms, fmt.Sprintf("treasure room with items %v", items))
}

func (m *mapBuilder) CreateBossArena(boss boss) {
	m.rooms = append(m.rooms, fmt.Sprintf("boss arena with boss %v", boss))
}

func (m *mapBuilder) GetRoomsDescription() string {
	return strings.Join(m.rooms, " --> ")
}

type mapDifficultyCounter struct {
	counter int
}

func newMapDifficultyCounter() *mapDifficultyCounter {
	return &mapDifficultyCounter{counter: 0}
}

func (m *mapDifficultyCounter) CreateDungeon(size, enemies int) {
	m.counter += enemies*2 - size
}

func (m *mapDifficultyCounter) CreateCheckpoint(healPoints int) {
	m.counter -= healPoints * 3
}

func (m *mapDifficultyCounter) CreateTreasureRoom(items []item) {
	for _, item := range items {
		m.counter -= item.Cost
	}
}

func (m *mapDifficultyCounter) CreateBossArena(boss boss) {
	m.counter += boss.Difficulty * 4
}

func (m *mapDifficultyCounter) GetDifficulty() int {
	return m.counter
}

// uncomment this to test

func main() {
	mapCreator := newMapBuilder()
	diffCreator := newMapDifficultyCounter()

	dir := newDirector(mapCreator)
	dir.CreateFirstLevel()
	fmt.Println("Map builder:")
	fmt.Println(mapCreator.GetRoomsDescription())

	dir = newDirector(diffCreator)
	dir.CreateFirstLevel()
	fmt.Println("Difficulty builder:", diffCreator.GetDifficulty())
}
