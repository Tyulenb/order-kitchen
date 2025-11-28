package app

import (
	"log"
	"net/http"

	"github.com/Tyulenb/order-kitchen/pos/internal/transport"
)

type App struct {
    Port string
}

func NewApp(port string) *App {
    return &App{
        Port: port,
    }
}

func (a *App) Run() error{
    router := http.NewServeMux()

    tr := transport.NewTransport()
    tr.RegisterRoutes(router)
    server := &http.Server{
        Addr: a.Port,
        Handler: router,
    }
    log.Printf("Server started on port: %v\n", a.Port)
    return server.ListenAndServe()
}
