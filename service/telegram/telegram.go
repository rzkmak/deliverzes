package telegram

import (
    "fmt"
    "github.com/aeidelos/deliverzes/config"
    "github.com/aeidelos/deliverzes/constant"
    tb "github.com/demget/telebot"
    "log"
    "strings"
)

type Bot struct {
    C config.Config
    B *tb.Bot
}

func NewTelegramBot(b *tb.Bot) *Bot {
    return &Bot{B: b}
}

func (t *Bot) OnCreate() {
    log.Println("starting telegram bot..")
}

func (t *Bot) OnHelp() {
    t.B.Handle(constant.HelpCommand, func(m *tb.Message) {
        t.send(m, constant.HelpMessage)
    })
}

func (t *Bot) OnUnknown() {
    t.B.Handle(tb.OnText, func(m *tb.Message) {
        t.send(m, fmt.Sprintf(constant.UnknownMessageReply, m.Text))
    })
}

func (t *Bot) OnSubscribe() {
    t.B.Handle(constant.SubscribeCommand, func(m *tb.Message) {
        split := strings.Split(m.Text, " ")
        if len(split) < 2 {
            t.send(m, constant.SubscribeParameterRequired)
            return
        }

        topics := split[1]
    })
}

func (t *Bot) OnUnsubscribe() {
}

func (t *Bot) send(m *tb.Message, message string) {
    m, err := t.B.Send(m.Sender, message)
    if err != nil {
        log.Println(err)
    }
}

func (t *Bot) Run() {
    t.OnCreate()
    t.OnHelp()
    t.OnSubscribe()
    t.OnUnsubscribe()

    go t.B.Start()
}


