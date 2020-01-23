package main

import (
    "fmt"
    "github.com/aeidelos/deliverzes/config"
    "github.com/aeidelos/deliverzes/service/telegram"
    tb "github.com/demget/telebot"
    "github.com/dgraph-io/badger"
    "log"
    "net/http"
    "time"
)

func main() {
    c := config.NewConfig()
    p := &tb.LongPoller{Timeout: 15 * time.Second}
    d, err := badger.Open(badger.DefaultOptions(c.DbPath))
    if err != nil {
        log.Fatalln(err)
    }
    t, err := tb.NewBot(tb.Settings{
        Token:  c.TelegramToken,
        Poller: p,
    })
    if err != nil {
        log.Fatalln(err)
    }
    b := telegram.NewBot(t, c, d)
    b.Run()

    http.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
        _, _ = writer.Write([]byte("pong"))
    })

    http.HandleFunc("/create", b.GenerateSubscriberIdHandler)
    http.HandleFunc("/send", b.SendMessageHandler)
    log.Printf("starting listen to :%v..", c.HttpPort)
    go log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", c.HttpPort), nil))
}
