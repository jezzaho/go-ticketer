package application

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jezzaho/go-ticketer/handler"
	"github.com/jezzaho/go-ticketer/repository/board"
	"github.com/jezzaho/go-ticketer/repository/ticket"
)

func (app *App) loadRoutes() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/tickets", app.loadTicketsRoutes)
	router.Route("/boards", app.loadBoardRoutes)
	app.router = router
}

func (app *App) loadTicketsRoutes(router chi.Router) {
	ticketHandler := &handler.Ticket{
		Repo: &ticket.RedisRepo{
			Client: app.rdb,
		},
	}

	router.Post("/", ticketHandler.Create)
	router.Get("/{id}", ticketHandler.GetByID)
	router.Get("/", ticketHandler.List)
	router.Put("/{id}", ticketHandler.UpdateByID)
	router.Delete("/{id}", ticketHandler.DeleteByID)
}

func (app *App) loadBoardRoutes(router chi.Router) {
	boardHandler := &handler.Board{
		Repo: &board.RedisRepo{
			Client: app.rdb,
		},
	}

	router.Post("/", boardHandler.Create)
	router.Get("/{id}", boardHandler.GetByID)
	router.Get("/", boardHandler.List)
	router.Put("/{id}", boardHandler.Update)
	router.Delete("/{id}", boardHandler.Delete)
}
