package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	pb "github.com/Tyulenb/order-kitchen/proto"
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
    ctx := context.Background()
    id, err := r.rbd.Incr(ctx, "OrderId").Result()
    if err != nil {
        return err
    }

    dishes := make(map[string]string)
    for {
        req, err := stream.Recv()
        if err == io.EOF {
            err := r.rbd.HSet(ctx, fmt.Sprintf("order:%d", id), "status", "Is Cooking").Err()
            if err != nil {
                return err
            }
            err = r.rbd.HSet(ctx, fmt.Sprintf("order:%d:dishes", id), dishes, ).Err()
            if err != nil {
                return err
            }
            return stream.SendAndClose(&pb.OrderId{Id: strconv.FormatInt(id, 10)})
        }
        if err != nil {
            return err 
        }
        dishes[req.DishName] = strconv.Itoa(int(req.Amount))
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
    //Remove "order:id:dishes
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

func (r *Restaurant) UpdateOrderStatus(ctx context.Context, osi *pb.OrderStatusId) (*pb.OrderId, error) {
    err := r.rbd.HSet(ctx, fmt.Sprintf("order:%s", osi.Id), "status", osi.Status).Err()
    if err != nil {
        return nil, err
    }
    return &pb.OrderId{Id: osi.Id}, nil
}

func (r *Restaurant) GetLastOrder(ctx context.Context, empty *pb.Empty) (*pb.OrderId, error) {
    lastOrder, err := r.rbd.Get(ctx, "lastOrder").Result()
    if err != nil {
        if err != redis.Nil {
            return nil, err
        }
        _, err := r.rbd.Incr(ctx, "lastOrder").Result()
        if err != nil {
            return nil, err
        }
        lastOrder, err = r.rbd.Get(ctx, "lastOrder").Result()
        if err != nil {
            return nil, err
        }
    }

    orderId, err := r.rbd.Get(ctx, "OrderId").Result()
    if err != nil {
        return nil, err
    }

    if lastOrder >= orderId {
        return &pb.OrderId{Id: fmt.Sprint(lastOrder)}, nil
    }
    order, err := r.rbd.Incr(ctx, "lastOrder").Result()
    if err != nil {
        return nil, err
    }
    
    return &pb.OrderId{Id: fmt.Sprint(order)}, nil
}

func (r *Restaurant) GetOrderDishes(order *pb.OrderId, stream pb.Restaurant_GetOrderDishesServer) error {
    id := order.GetId()
    dishes, err := r.rbd.HGetAll(context.Background(), fmt.Sprintf("order:%s:dishes", id)).Result()
    if err != nil {
        return err
    }
    for k, v := range dishes {
        amount, err := strconv.ParseInt(v, 10, 32)
        if err != nil {
            return err
        }
        stream.Send(&pb.OrderDishes{DishName: k, Amount: int32(amount)})
    }
    return nil 
}

func (r *Restaurant) GetOrderStatus(ctx context.Context, orderId *pb.OrderId) (*pb.OrderStatus, error) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
    defer cancel()

    status, err := r.rbd.HGetAll(ctx, fmt.Sprintf("order:%s", orderId.Id)).Result()
    if err != nil {
        return nil, err
    } 

    return &pb.OrderStatus{Status: status["status"]}, nil
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
