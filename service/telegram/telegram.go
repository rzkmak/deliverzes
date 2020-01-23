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
        m, err := t.B.Send(m.Sender, constant.HelpMessage)
        if err != nil {
            log.Println(err)
        }
    })
}

func (t *Bot) OnUnknown() {
    t.B.Handle(tb.OnText, func(m *tb.Message) {
        m, err := t.B.Send(m.Sender, fmt.Sprintf(constant.UnknownMessageReply, m.Text))
        if err != nil {
            log.Println(err)
        }
    })
}

func (t *Bot) OnSubscribe() {
    t.B.Handle(constant.SubscribeCommand, func(m *tb.Message) {
        split := strings.Split(m.Text, " ")
        if len(split) < 2 {
            _, err := t.B.Send(m.Sender, constant.SubscribeParameterRequired)
            if err != nil {
                log.Println(err)
            }
            return
        }

        topics := split[1]
    })
}

func (t *Bot) OnUnsubscribe() {
}

func (t *Bot) Run() {
    t.OnCreate()
    t.OnHelp()
    t.OnSubscribe()
    t.OnUnsubscribe()

    go t.B.Start()
}


