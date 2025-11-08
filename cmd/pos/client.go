package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/Tyulenb/order-kitchen/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PoS struct {
    client pb.RestaurantClient
}

func NewPos(client pb.RestaurantClient) *PoS {
    return &PoS{client: client}
}

func (p *PoS) makeOrder() {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
    defer cancel()
    order := make([]*pb.OrderRequest, 0)
    order = append(order, &pb.OrderRequest{DishName: "Scrambled eggs", Amount: 1})
    order = append(order, &pb.OrderRequest{DishName: "Orange juice", Amount: 1})
    stream, err := p.client.CreateOrder(ctx)
    if err != nil {
        log.Fatalf("stream, %v", err)
    }
    for _, v := range order {
        stream.Send(v)
    }
    response, err := stream.CloseAndRecv()
    if err != nil {
        log.Fatalf("Close stream, %v", err)
    }
    log.Println(response.Id)
}

func (p *PoS) listOrderStatus() {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
    defer cancel()
    stream, err := p.client.ListOrderStatus(ctx, &pb.Empty{})
    if err != nil {
        log.Fatalf("listOrderStatus, open stream, %v", err)
        return
    }
    for {
        statusId, err := stream.Recv()
        if err == io.EOF {
            break
        }
        log.Println(statusId.Id, statusId.Status)
    }
}

func main(){
    conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatal("did not connect:", err)
    }
    defer conn.Close()
    client := pb.NewRestaurantClient(conn)
    pos := NewPos(client)
    pos.makeOrder()
    pos.listOrderStatus()
}
