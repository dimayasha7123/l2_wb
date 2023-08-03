package httpserver

import (
	"net/http"
	"time"

	"l2_wb/develop/dev11/internal/app"
	"l2_wb/develop/dev11/internal/inputs/httpserver/handlers"
	"l2_wb/develop/dev11/internal/inputs/httpserver/middlewares"
)

//POST 	/create_user
//POST	/create_event
//POST 	/update_event
//POST 	/delete_event
//GET 	/events_for_day
//GET 	/events_for_week
//GET 	/events_for_month

// New constructor
func New(service *app.App, addr string) *http.Server {
	router := http.NewServeMux()
	router.Handle("/create_user", middlewares.NewResultsWrapper(handlers.NewCreateUserHandler(service)))
	router.Handle("/create_event", middlewares.NewResultsWrapper(handlers.NewCreateEventHandler(service)))
	router.Handle("/update_event", middlewares.NewResultsWrapper(handlers.NewUpdateEventHandler(service)))
	router.Handle("/delete_event", middlewares.NewResultsWrapper(handlers.NewDeleteEventHandler(service)))
	router.Handle("/events_for_day", middlewares.NewResultsWrapper(handlers.NewEventsForDayHandler(service)))
	router.Handle("/events_for_week", middlewares.NewResultsWrapper(handlers.NewEventsForWeekHandler(service)))
	router.Handle("/events_for_month", middlewares.NewResultsWrapper(handlers.NewEventsForMonthHandler(service)))
	router.Handle("/events", middlewares.NewResultsWrapper(handlers.NewEventsAllHandler(service)))
	router.Handle("/", middlewares.NewResultsWrapper(handlers.NewNotFoundHandler()))

	srv := &http.Server{
		Addr:         addr,
		Handler:      middlewares.PanicMW(middlewares.LoggingMW(router)),
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
		IdleTimeout:  time.Second * 10,
	}

	return srv
}
