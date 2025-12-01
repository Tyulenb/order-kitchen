package service

import (
	"context"
	"io"
	"log"
	"strings"
	"time"

	pb "github.com/Tyulenb/order-kitchen/proto"
)

type Service struct {
    client pb.RestaurantClient
}

func NewService(client pb.RestaurantClient) *Service {
    return &Service{client: client}
}

func (s *Service) MakeOrder(cb, ff, cl int32) error{
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
    defer cancel()

    order := make([]*pb.OrderDishes, 0)
    order = append(order, &pb.OrderDishes{DishName: "Cheesebureger", Amount: cb})
    order = append(order, &pb.OrderDishes{DishName: "French Fries", Amount: ff})
    order = append(order, &pb.OrderDishes{DishName: "Cola", Amount: cl})

    stream, err := s.client.CreateOrder(ctx)
    if err != nil {
        return err
    }
    for _, v := range order {
        stream.Send(v)
    }
    response, err := stream.CloseAndRecv()
    if err != nil {
        return err
    }
    log.Println(response.Id)
    return nil
}

//If param == "" returns all orders without filtering
func (s *Service) listOrderStatus(param string) ([]*pb.OrderStatusId, error) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
    defer cancel()
    orders := make([]*pb.OrderStatusId, 0)
    stream, err := s.client.ListOrderStatus(ctx, &pb.Empty{})
    if err != nil {
        return nil, err
    }
    for {
        statusId, err := stream.Recv()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, err
        }
        if statusId != nil { 
            orders = append(orders, statusId)
        }
    }

    if param == "" {
        return orders, nil
    }
    filtred := make([]*pb.OrderStatusId, 0)
    for i := range orders {
        if strings.Compare(orders[i].GetStatus(), param) == 0 {
            filtred = append(filtred, orders[i])
        }
    }
    return filtred, nil
}

func (s *Service) ListCookingOrders() ([]string, error) {
    orders, err := s.listOrderStatus("Cooking")
    if err != nil {
        return nil, err
    }
    cooking := make([]string, 0)
    for i := range orders {
        cooking = append(cooking, orders[i].Id)
    }
    return cooking, nil
}

func (s *Service) ListReadyOrders() ([]string, error) {
    orders, err := s.listOrderStatus("Cooked")
    if err != nil {
        return nil, err
    }
    cooked := make([]string, 0)
    for i := range orders {
        cooked = append(cooked, orders[i].Id)
    }
    return cooked, nil
}

func (s *Service) UpdateOrderStatus () {
    ctx := context.Background()
    r, err := s.client.UpdateOrderStatus(ctx, &pb.OrderStatusId{Id: "1", Status: "Canceled"})
    if err != nil {
        log.Fatalf("updateOrderStatus %v", err)
    }
    log.Println(r)
}

func (s *Service) GetLastOrder() {
    ctx := context.Background()
    r, err := s.client.GetLastOrder(ctx, &pb.Empty{})
    if err != nil {
        log.Fatalf("getLastOrder %v", err)
    }
    log.Println(r)
}

func (s *Service) GetOrderDishes() {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
    defer cancel()
    stream, err := s.client.GetOrderDishes(ctx, &pb.OrderId{Id: "1"})
    if err != nil {
        log.Fatal(err)
    }
    for {
        orderDish, err := stream.Recv()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Fatalf("getLastOrder %v", err)
        }
        log.Println(orderDish)
    }
}
