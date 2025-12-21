package transport

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Tyulenb/order-kitchen/pos/internal/service"
)

type Transport struct {
    srv *service.Service
}

func NewTransport(srv *service.Service) *Transport {
    return &Transport{
        srv: srv,
    }
}

func (t *Transport) RegisterRoutes(router *http.ServeMux) {
    router.Handle("/", http.FileServer(http.Dir("pos/web")))
    router.HandleFunc("POST /order", t.makeOrder)
    router.HandleFunc("GET /cookingOrders", t.cookingOrders)
    router.HandleFunc("GET /readyOrders", t.readyOrders)
}

func (t *Transport) makeOrder(w http.ResponseWriter, r *http.Request) {

    if r.Body == nil{
        smthWentWrong(w, fmt.Errorf("The body is nil"))
        return
    }
    defer r.Body.Close()

    type Order struct {
        QuantityBurger int 
        QuantityFries int 
        QuantityCola int 
    }
    var o Order
    err := json.NewDecoder(r.Body).Decode(&o)
    if err != nil {
        smthWentWrong(w, err)
        return 
    }

    err = t.srv.MakeOrder(int32(o.QuantityBurger), int32(o.QuantityFries), int32(o.QuantityCola))
    if err != nil {
        smthWentWrong(w, err)
        return
    }
    log.Printf("Make Order: %d, %d, %d", o.QuantityBurger, o.QuantityFries, o.QuantityCola)
}

func (t *Transport) readyOrders(w http.ResponseWriter, r *http.Request) {
    cooked, err := t.srv.ListReadyOrders()
    if err != nil {
        smthWentWrong(w, err)
        return
    }
    var responseString string
    for i := range cooked {
        responseString += fmt.Sprintf(`<h2>Order№ %s</h2>`, cooked[i])
    }
    fmt.Fprint(w, responseString)
}

func (t *Transport) cookingOrders(w http.ResponseWriter, r *http.Request) {
    cooking, err := t.srv.ListCookingOrders()
    if err != nil {
        smthWentWrong(w, err)
        return
    }
    var responseString string
    for i := range cooking {
        responseString += fmt.Sprintf(`<h2>Order№ %s</h2>`, cooking[i])
    }
    fmt.Fprint(w, responseString)
}

func smthWentWrong(w http.ResponseWriter, err error) {
    http.Error(w, "Something went wrong", 400)
    log.Println(err)
}
