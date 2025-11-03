package main

import (
	"context"
	"log"
	"time"

	pb "github.com/Tyulenb/order-kitchen/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main(){
    conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatal("did not connect:", err)
    }

    defer conn.Close()
    client := pb.NewOrderClient(conn)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
    defer cancel()
    response, err := client.CreateOrder(ctx, &pb.OrderRequest{Item: "Pork"})
    if err != nil {
        log.Fatal("could not get: Pork")
    }
    log.Println(response)
    resp, err := client.GetOrderStatus(ctx, &pb.OrderId{Id: "2"})
    if err != nil {
        log.Fatal("could not get order status 2")
    }
    log.Println(resp)
}
