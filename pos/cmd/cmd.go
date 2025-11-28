package main

import "github.com/Tyulenb/order-kitchen/pos/internal/app"

func main() {
    ap := app.NewApp(":9999")
    ap.Run()
}
