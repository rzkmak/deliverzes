package main

import (
    "fmt"
    "github.com/aeidelos/deliverzes/config"
    "log"
    "net/http"
)

func main() {
    c := config.NewConfig()
    http.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
        _, _ = writer.Write([]byte("pong"))
    })
    log.Printf("starting listen to :%v..", c.HttpPort)
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", c.HttpPort), nil))
}
