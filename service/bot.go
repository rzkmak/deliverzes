package service

type Bot interface {
    OnCreate()
    OnHelp()
    OnUnknown()
    OnSubscribe()
    OnUnsubscribe()
    Run()
}
