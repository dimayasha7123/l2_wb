package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"l2_wb/develop/dev11/internal/adapters/memrepo"
	"l2_wb/develop/dev11/internal/app"
	"l2_wb/develop/dev11/internal/inputs/httpserver"
)

/*
=== HTTP server ===

Реализовать HTTP сервер для работы с календарем.
В рамках задания необходимо работать строго со стандартной HTTP библиотекой.

В рамках задания необходимо:
	1. Реализовать вспомогательные функции для сериализации объектов доменной области в JSON.
	2. Реализовать вспомогательные функции для парсинга и валидации параметров методов /create_event и /update_event.
	3. Реализовать HTTP обработчики для каждого из методов API, используя вспомогательные функции и объекты доменной области.
	4. Реализовать middleware для логирования запросов

Методы API:
	POST	/create_event
	POST 	/update_event
	POST 	/delete_event
	GET 	/events_for_day
	GET 	/events_for_week
	GET 	/events_for_month

Добавил еще два вспомогательных:
	POST	/create_user (создает нового пользователя с ником nickname)
	GET		/events (все события)

Параметры передаются в виде www-url-form-encoded (т.е. обычные user_id=3&date=2019-09-09).
В GET методах параметры передаются через queryString, в POST через тело запроса.
В результате каждого запроса должен возвращаться JSON документ содержащий либо {"result": "..."}
в случае успешного выполнения метода, либо {"error": "..."} в случае ошибки бизнес-логики.

В рамках задачи необходимо:
	1. Реализовать все методы.
	2. Бизнес логика НЕ должна зависеть от кода HTTP сервера.
	3. В случае ошибки бизнес-логики сервер должен возвращать HTTP 503. В случае ошибки входных данных
	   (невалидный int например) сервер должен возвращать HTTP 400. В случае остальных ошибок сервер
	   должен возвращать HTTP 500. Web-сервер должен запускаться на порту указанном в конфиге и выводить
	   в лог каждый обработанный запрос.
	4. Код должен проходить проверки go vet и golint.
*/

func main() {
	var port = "8001"
	flag.StringVar(&port, "p", port, "port to start HTTP server")
	flag.Parse()

	addr := fmt.Sprintf("localhost:%s", port)

	repo := memrepo.New()
	service := app.New(repo)

	createUserResp, err := service.CreateUser(app.CreateUserReq{Nickname: "admin"})
	if err != nil {
		log.Fatalf("Can't create admin user: %v", err)
	}
	log.Printf("Admin user has id = %d", createUserResp.UserID)
	log.Println("Its just default user to manual testing")

	server := httpserver.New(service, addr)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Println("Server was closed")
			} else {
				log.Fatalf("Get error while serving: %v", err)
			}
		}
	}()
	log.Printf("Server started at http://%s", addr)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Println()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Fatalf("Get error while shutting down server: %v", err)
	}

	log.Println("Server was shutdown successfully!")
}
