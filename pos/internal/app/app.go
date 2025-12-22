package app

import (
	"log"
	"net/http"

	"github.com/Tyulenb/order-kitchen/pos/internal/service"
	"github.com/Tyulenb/order-kitchen/pos/internal/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "github.com/Tyulenb/order-kitchen/proto"
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
    
    conn, err := grpc.NewClient("server:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        return err
    }
    client := pb.NewRestaurantClient(conn)

    sr := service.NewService(client)

    tr := transport.NewTransport(sr)
    tr.RegisterRoutes(router)
    server := &http.Server{
        Addr: a.Port,
        Handler: router,
    }
    log.Printf("Server started on port: %v\n", a.Port)
    return server.ListenAndServe()
}
