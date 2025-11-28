package transport

import (
	"net/http"

)

type Transport struct {
    //service
}

func NewTransport() *Transport {
    return &Transport{
    }
}

func (t *Transport) RegisterRoutes(router *http.ServeMux) {
    router.HandleFunc("/", t.homePage)
}

func (t *Transport) homePage(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type","text/html; charset=utf-8")
    http.ServeFile(w, r, "pos/web/index.html")
}
