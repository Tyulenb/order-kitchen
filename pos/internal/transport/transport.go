package transport

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

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
    log.Println("makeOrder")
    err := r.ParseForm()
    if err != nil {
        smthWentWrong(w, err)
        return
    }
    chsBurg, err := strconv.ParseInt(r.FormValue("quantity1"), 10, 32)
    if err != nil {
        smthWentWrong(w, err)
        return
    }
    frFries, err := strconv.ParseInt(r.FormValue("quantity2"), 10, 32)
    if err != nil {
        smthWentWrong(w, err)
        return
    }
    cola, err := strconv.ParseInt(r.FormValue("quantity3"), 10, 32)
    if err != nil {
        smthWentWrong(w, err)
        return
    }
    fmt.Println(chsBurg, frFries, cola) 
    err = t.srv.MakeOrder(int32(chsBurg), int32(frFries), int32(cola))
    if err != nil {
        smthWentWrong(w, err)
        return
    }
}

func (t *Transport) readyOrders(w http.ResponseWriter, r *http.Request) {
    cooked, err := t.srv.ListReadyOrders()
    if err != nil {
        smthWentWrong(w, err)
        return
    }
    for i := range cooked {
        fmt.Fprintln(w, cooked[i])
    }
}

func (t *Transport) cookingOrders(w http.ResponseWriter, r *http.Request) {
    cooking, err := t.srv.ListCookingOrders()
    if err != nil {
        smthWentWrong(w, err)
        return
    }
    for i := range cooking {
        fmt.Fprintln(w, cooking[i])
    }
}

func smthWentWrong(w http.ResponseWriter, err error) {
    http.Error(w, "Something went wrong", 400)
    log.Println(err)
}
