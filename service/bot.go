package service

import "net/http"

type Bot interface {
    OnCreate()
    OnHelp()
    OnUnknown()
    OnSubscribe()
    OnUnsubscribe()
    OnUnsubscribeAll()
    GenerateSubscriberIdHandler(http.ResponseWriter, *http.Request)
    SendMessageHandler(http.ResponseWriter, *http.Request)
    Run()
}
