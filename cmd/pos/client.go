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
    /////////////////////////////////////////////////////////////////////
    order := make([]*pb.OrderRequest, 0)
    order = append(order, &pb.OrderRequest{DishName: "Scrambled eggs"})
    order = append(order, &pb.OrderRequest{DishName: "Orange juice"})
    stream, err := client.CreateOrder(ctx)
    if err != nil {
        log.Fatal(err)
    }
    for _, v := range order {
        stream.Send(v)
    }
    response, err := stream.CloseAndRecv()
    if err != nil {
        log.Fatal(err)
    }
    log.Println(response.Id)
    ///////////////////////////////////////////////////////////////////

    resp, err := client.GetOrderStatus(ctx, &pb.OrderId{Id: response.Id})
    if err != nil {
        log.Fatal("could not get order status 2")
    }
    log.Println(resp)
}
