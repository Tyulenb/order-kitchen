package main

import (
	"context"
	"io"
	"log"
	"net"

	pb "github.com/Tyulenb/order-kitchen/proto"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

type Order struct {
    pb.UnimplementedOrderServer
    rbd *redis.Client
}

func NewOrder(rbd *redis.Client) *Order {
    return &Order{
        rbd: rbd,
    }
}

func (o *Order) CreateOrder(stream pb.Order_CreateOrderServer) error {
    id := uuid.NewString()
    ctx := context.Background()
    dishes := make(map[string]string)
    for {
        req, err := stream.Recv()
        if err == io.EOF {
            err := o.rbd.HSet(ctx, id, dishes).Err()
            if err != nil {
                return err
            }
            return stream.SendAndClose(&pb.OrderResponse{Id: id})
        }
        if err != nil {
            return err 
        }
        dishes[req.DishName] = "Cooking"
    }
}

func (o *Order) GetOrderStatus(ctx context.Context, id *pb.OrderId) (*pb.OrderStatusResponse, error) {
    log.Println("Status request was accepted:", id.Id)
    item, err := o.rbd.HGetAll(ctx, id.Id).Result()
    if err != nil {
        return &pb.OrderStatusResponse{}, err
    }
    log.Println(item)
    return &pb.OrderStatusResponse{Id: id.Id, Status: "status"}, nil
}

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatal("Failed to listen 50051")
    }

    redisdb := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
        Password: "",
        DB: 0,
    })
    defer redisdb.Close()

    if err := PingDB(redisdb); err != nil {
        log.Fatalf("Cannot connect to DB: %v", err)
    }
    order := NewOrder(redisdb)

    grpcServer := grpc.NewServer()
    pb.RegisterOrderServer(grpcServer, order)
    log.Println("Server listening at:", lis.Addr())
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatal(err)
    }
}

func PingDB(rdb *redis.Client) error {
    ctx := context.Background()
    _, err := rdb.Ping(ctx).Result()
    if err == nil {
        log.Println("Database successfully connected")
    }
    return err
}
