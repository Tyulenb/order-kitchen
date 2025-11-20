package main

import (
	"context"
	"io"
	"log"
	"strings"
	"time"

	pb "github.com/Tyulenb/order-kitchen/proto"
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

func cooker(order chan string) {
    for i := range order {   

    }
}

func main() {
}
