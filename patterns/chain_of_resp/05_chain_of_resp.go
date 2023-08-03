package main

import "fmt"

/*
	Реализовать паттерн «цепочка вызовов».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern
*/

/*
	ПРИМЕНИМОСТЬ:
1.	Необходимость обработать запросы в последовательности нескольких объектов, каждый из которых
	может либо обработать запрос, либо передать его дальше по цепочке.

	ПЛЮСЫ:
1.	Уменьшение связанности между объектами, так как каждый объект знает только о следующем объекте в цепочке.
2.	Гибкость и расширяемость системы, так как можно добавлять новые объекты в цепочку без изменения существующего кода.
3.	Возможность динамической настройки цепочки вызовов.

	МИНУСЫ:
1.	Усложнение отладки и тестирования системы из-за необходимости отслеживать прохождение запроса через
	несколько объектов.
2.	Возможность возникновения зацикливания при неправильной реализации цепочки вызовов.

	ПРИМЕРЫ:
1.	Цепочки интерцепторов запросов к http/grpc серверу.
2.	Обработка ошибок (возможно ошибку передадут дальше, или отработают ее)
*/

type handler interface {
	setNext(handler handler) handler
	handle(request string) string
}

type abstractHandler struct {
	nextHandler handler
}

func (h *abstractHandler) setNext(handler handler) handler {
	h.nextHandler = handler
	return handler
}

func (h *abstractHandler) handle(request string) string {
	if h.nextHandler != nil {
		return h.nextHandler.handle(request)
	}
	return ""
}

type authenticationHandler struct {
	abstractHandler
}

func (h *authenticationHandler) handle(request string) string {
	if request == "authenticated" {
		return h.abstractHandler.handle(request)
	}
	return "Access denied"
}

type authorizationHandler struct {
	abstractHandler
}

func (h *authorizationHandler) handle(request string) string {
	if request == "authorized" {
		return h.abstractHandler.handle(request)
	}
	return "Not authorized"
}

type requestHandler struct {
	abstractHandler
}

func (h *requestHandler) handle(request string) string {
	return "Request processed"
}

func main() {
	request := "authenticated, authorized, request"

	authenticationHandler := &authenticationHandler{}
	authorizationHandler := &authorizationHandler{}
	requestHandler := &requestHandler{}

	authenticationHandler.setNext(authorizationHandler).setNext(requestHandler)

	result := authenticationHandler.handle(request)

	fmt.Println(result)
}
