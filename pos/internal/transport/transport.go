package transport

import (
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
}

func (t *Transport) makeOrder(w http.ResponseWriter, r *http.Request) {
    log.Println("makeOrder")
    err := r.ParseForm()
    if err != nil {
        http.Error(w, "Something went wrong", 400)
        log.Println(err)
        return
    }
    chsBurg_str := r.FormValue("quantity1")
    frFries_str := r.FormValue("quantity2")
    cola_str := r.FormValue("quantity3")
    fmt.Println(chsBurg_str, frFries_str, cola_str) 
    /*
    err = t.srv.MakeOrder(chsBurg, frFries, cola)
    if err != nil {
        http.Error(w, "Something went wrong", 400)
        log.Println(err)
        return
    }
    */
}
