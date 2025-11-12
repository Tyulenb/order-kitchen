package main

import (
	"sync"

	pb "github.com/Tyulenb/order-kitchen/proto"
)

type SyncQueue struct {
    orders []string
    mtx sync.RWMutex
}

type Kitchen struct {
    client pb.RestaurantClient
}

func NewKitchen(client pb.RestaurantClient) *Kitchen {
    return &Kitchen{client: client}
}

func main() {
    
}
