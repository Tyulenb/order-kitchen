package main

import (
	"context"
	"log"
	"net"

	pb "github.com/Tyulenb/order-kitchen/proto"
	"google.golang.org/grpc"
)

type Order struct {
    pb.UnimplementedOrderServer
}

func (o *Order) CreateOrder(ctx context.Context, req *pb.OrderRequest) (*pb.OrderResponse, error){
    log.Println("Your order was accepted:", req.Item)
    return &pb.OrderResponse{Id: "1"}, nil
}

func (o *Order) GetOrderStatus(ctx context.Context, id *pb.OrderId) (*pb.OrderStatusResponse, error) {
    log.Println("Status request was accepted:", id.Id)
    return &pb.OrderStatusResponse{Id: id.Id, Status: "Is cooking"}, nil
}

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatal("Failed to listen 50051")
    }
    grpcServer := grpc.NewServer()
    pb.RegisterOrderServer(grpcServer, &Order{})
    log.Println("Server listening at:", lis.Addr())
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatal(err)
    }
}
