package main

import (
	pb "github.com/Tyulenb/order-kitchen/proto"
)

type Kitchen struct {
    client pb.RestaurantClient
}

func NewKitchen(client pb.RestaurantClient) *Kitchen {
    return &Kitchen{client: client}
}

func main() {
    
}
