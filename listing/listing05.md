Что выведет программа? Объяснить вывод программы.

```go
package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	{
		// do something
	}
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
```

Ответ:
```
error
```
Ситуация аналогична той, что была в листинге №3.

Завели переменную err интерфейсного типа error.
Вызов err = test() вернул нулевой тип, реализующий интерфейс Error.
Таким образом проверка на err на nil возвращает true, т.к. лежащий внутри
интерфейса тип определен, хоть и реальное значение nil. 
