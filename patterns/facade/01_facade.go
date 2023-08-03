package facade

/*
	Реализовать паттерн «фасад».
	Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Facade_pattern
*/

/*
	ПРИМЕНИМОСТЬ:
1.	Предоставление простого интерфейса к сложной подсистеме (подсистема состоит из кучи различных мелких объектов со сложной схемой
	взаимодействия, а клиентам этой подсистемы необходим некоторый простой и устраивающий большинство вид по умолчания)
2. 	Уменьшение количества зависимостей между клиентами системы и классами ее реализации (отделение подсистемы от клиентов
	способствует независимости подсистем и повышению уровня переносимости)
3.	Разбиение системы на отдельные уровни (использование только фасада в качестве точки входа на каждый из уровней подсистемы)

	ПЛЮСЫ:
1.	Упрощение интерфейса сложной системы
2.	Сокрытие деталей взаимодействия со сложной системой
3.	Увеличение независимости компонентов и простота поддержки

	МИНУСЫ:
1.	Потеря гибкости взаимодействия с системой (отдельным клиентам может потребоваться гибкость выбора между подклассами, или же
	доступ к инкапсулированным фасадом функциям)
2.	Возможно увеличение времени работы из-за недостаточной оптимизации

	ПРИМЕРЫ:
1.	Пакет "database/sql" предоставляет один простой доступ к различным базам данных
2.	Различные веб-фреймворки предоставляют упрощенный способ разрабатывать веб-приложения (маршрутиризация, обработка запросов и т.д.)
*/

type repository struct {
	cache *cache
	db    *db
}

func newRepository(db *db) (*repository, error) {
	values, err := db.getValues()
	if err != nil {
		return nil, err
	}

	return &repository{
		cache: newCache(values),
		db:    db,
	}, nil
}

func (r *repository) getValue(key string) (string, error) {
	value, err := r.cache.getValue(key)
	if err == nil {
		return value, nil
	}

	value, err = r.db.getValue(key)
	if err != nil {
		return "", err
	}

	err = r.cache.addValue(key, value)
	if err != nil {
		// log it
	}

	return value, nil
}

func (r *repository) addValue(data string) (string, error) {
	key, err := r.db.addValue(data)
	if err != nil {
		return "", err
	}

	err = r.cache.addValue(key, data)
	if err != nil {
		// log it
	}

	return key, nil
}

type cache struct {
	// cache fields
}

func (c *cache) getValue(key string) (string, error) {
	// cache getValue realization
	return "some data", nil
}

func (c *cache) addValue(key string, data string) error {
	// cache addValue realization
	return nil
}

func newCache(values [][]string) *cache {
	// newCache realization
	return nil
}

type db struct {
	// db fields
}

func (b *db) getValue(key string) (string, error) {
	// db getValue realization
	return "some data", nil
}

func (b *db) addValue(data string) (string, error) {
	// db addValue realization
	return "some key", nil
}

func (b *db) getValues() ([][]string, error) {
	// db getValues realization
	return [][]string{{"key1", "data1"}, {"key2", "data2"}}, nil
}
