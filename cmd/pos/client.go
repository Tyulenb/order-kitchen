package main

import (
	"context"
	"log"
	"time"

	pb "github.com/Tyulenb/order-kitchen/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PoS struct {
    client pb.OrderClient
}

func NewPos(client pb.OrderClient) *PoS {
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
}

func (p *PoS) getOrderStatus(id string) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
    defer cancel()
    resp, err := p.client.GetOrderStatus(ctx, &pb.OrderId{Id: id})
    if err != nil {
        log.Fatal("could not get order status 2")
    }
    log.Println(resp)
}

func main(){
    conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatal("did not connect:", err)
    }
    defer conn.Close()
    client := pb.NewOrderClient(conn)
    pos := NewPos(client)
    pos.makeOrder()
}
