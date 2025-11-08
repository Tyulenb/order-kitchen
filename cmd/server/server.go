package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	pb "github.com/Tyulenb/order-kitchen/proto"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

type Restaurant struct {
    pb.UnimplementedRestaurantServer
    rbd *redis.Client
}

func NewRestaurant(rbd *redis.Client) *Restaurant {
    return &Restaurant{
        rbd: rbd,
    }
}

func (r *Restaurant) CreateOrder(stream pb.Restaurant_CreateOrderServer) error {
    id := uuid.NewString()
    ctx := context.Background()
    dishes := make(map[string]string)
    for {
        req, err := stream.Recv()
        if err == io.EOF {
            err := r.rbd.HSet(ctx, "order:"+id, "status", "Is Cooking").Err()
            if err != nil {
                return err
            }
            err = r.rbd.HSet(ctx, fmt.Sprintf("order:%s:dishes", id), dishes).Err()
            if err != nil {
                return err
            }
            return stream.SendAndClose(&pb.OrderId{Id: id})
        }
        if err != nil {
            return err 
        }
        dishes[req.DishName] = string(req.Amount)
    }
}

func (r *Restaurant) ListOrderStatus(empty *pb.Empty, stream pb.Restaurant_ListOrderStatusServer) error {
    var orders []string
    var err error
    var cursor uint64
    for {
        var keysFromScan []string
        keysFromScan, cursor, err = r.rbd.Scan(context.TODO(), cursor, "order:*", 10).Result()
        if err != nil {
            return err
        }
        orders = append(orders, keysFromScan...)
        if cursor == 0 {
            break
        }
    }
    for i := range orders {
        parts := strings.Split(orders[i], ":")
        if len(parts) > 2 {
            continue
        }
        status, err := r.rbd.HGet(context.TODO(), orders[i], "status").Result()
        if err != nil {
            return err
        }
        stream.Send(&pb.OrderStatusId{Id: parts[1], Status: status})
    }
    return nil
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
    redisdb.FlushDB(context.Background())

    restaurant := NewRestaurant(redisdb)
    grpcServer := grpc.NewServer()
    pb.RegisterRestaurantServer(grpcServer, restaurant)
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
