package main

import (
	"context"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	pb "github.com/Tyulenb/order-kitchen/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Kitchen struct {
    client pb.RestaurantClient
}

func NewKitchen(client pb.RestaurantClient) *Kitchen {
    return &Kitchen{client: client}
}

func (k *Kitchen) SimulateCooking(id string) {
    timeToCook := 0
    timeList := make(map[string]int, 0)

    ctx := context.Background() 
    orderId := &pb.OrderId{Id: id}
    stream, err := k.client.GetOrderDishes(ctx, orderId)
    if err != nil {
        log.Fatalf("SimulateCooking, opening stream: %v\n", err)
    }
    for {
        dishes, err := stream.Recv()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Fatalf("SimulateCooking, receiving dishes %v\n", err)
        }
        timeToCook += timeList[dishes.DishName] * int(dishes.Amount)
    }

    timeout := time.After(time.Duration(timeToCook))
    statusChecker := time.NewTicker(time.Second*5)
    for {
        select {
        case <- timeout:
            _, err := k.client.UpdateOrderStatus(ctx, &pb.OrderStatusId{Id: id, Status: "Cooked"})
            if err != nil {
                log.Fatalf("SimulateCooking, UpdateOrderStatus error: %v\n", err)
            }
            return
        case <- statusChecker.C:
            status, err := k.client.GetOrderStatus(ctx, &pb.OrderId{Id: id})
            if err != nil {
                log.Fatalf("SimulateCooking, GetOrderStatus error: %v\n", err)
            }
            if strings.Compare(status.Status, "Canceled") == 0 {
                return
            }
        }
    }
    
}

func cooker(order chan string, k *Kitchen) {
    for i := range order {   
        k.SimulateCooking(i) 
    }
}

func main() {
    const numCookers = 5

    conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatal("did not connect:", err)
    }
    defer conn.Close()
    client := pb.NewRestaurantClient(conn)
    kitchen := NewKitchen(client)

    order := make(chan string, 5)
    for range numCookers {
        go cooker(order, kitchen)
    }

    prev_order := "" 
    ctx := context.Background()
    for {
        lorder, err := kitchen.client.GetLastOrder(ctx, &pb.Empty{})
        if err != nil {
            continue
        }
        if lorder.Id != prev_order {
            order <- lorder.Id
            prev_order = lorder.Id
        }
    }
}

