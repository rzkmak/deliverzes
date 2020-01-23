package main

import (
    "fmt"
    "github.com/aeidelos/deliverzes/config"
    "github.com/aeidelos/deliverzes/service/telegram"
    tb "github.com/demget/telebot"
    "log"
    "net/http"
    "time"
)

func main() {
    c := config.NewConfig()

    p := &tb.LongPoller{Timeout: 15 * time.Second}


    t, err := tb.NewBot(tb.Settings{
        Token:  c.TelegramToken,
        Poller: p,
    })

    if err != nil {
        log.Fatalln(err)
    }

    b := telegram.NewTelegramBot(t)

    b.Run()

    http.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
        _, _ = writer.Write([]byte("pong"))
    })
    log.Printf("starting listen to :%v..", c.HttpPort)
    go log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", c.HttpPort), nil))
}
