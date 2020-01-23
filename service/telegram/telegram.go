package telegram

import (
    "fmt"
    "github.com/aeidelos/deliverzes/constant"
    tb "github.com/demget/telebot"
    "github.com/dgraph-io/badger"
    "log"
    "strconv"
    "strings"
)

const ErrSubscriberEmpty = "subscriber is empty"

type Bot struct {
    B *tb.Bot
    D *badger.DB
}

func NewBot(b *tb.Bot, d *badger.DB) *Bot {
    return &Bot{B: b, D: d}
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
        split := strings.Split(strings.TrimSpace(m.Text), " ")
        if len(split) < 2 {
            t.send(m, constant.SubscribeParameterRequired)
            return
        }
        topics := split[1]
        err := t.D.View(func(txn *badger.Txn) error {
            btopics := []byte(topics)
            item, err := txn.Get(btopics)
            if err != nil {
                return err
            }
            if err := item.Value(func(val []byte) error {
                sender := strconv.Itoa(m.Sender.ID)
                subscribers := string(val)
                if subscribers == "" {
                    return txn.Set(btopics, []byte(sender))
                }
                arrSubscribers := strings.Split(subscribers, ",")
                updSubscribers := strings.Join(append(arrSubscribers, sender), ",")
                return txn.Set(btopics, []byte(updSubscribers))
            }); err != nil {
                return err
            }
            return nil
        })
        if err == badger.ErrKeyNotFound {
            t.send(m, fmt.Sprintf(constant.SubscribeIdNotFound, topics))
            return
        }
        if err != nil {
            log.Println(err)
            t.send(m, fmt.Sprintf(constant.SubscribeIdFailed, topics))
        }
    })
}

func (t *Bot) OnUnsubscribe() {
    t.B.Handle(constant.SubscribeCommand, func(m *tb.Message) {
        split := strings.Split(strings.TrimSpace(m.Text), " ")
        if len(split) < 2 {
            t.send(m, constant.UnsubscribeCommand)
            return
        }
        topics := split[1]
        err := t.D.View(func(txn *badger.Txn) error {
            btopics := []byte(topics)
            item, err := txn.Get(btopics)
            if err != nil {
                return err
            }
            if err := item.Value(func(val []byte) error {
                sender := strconv.Itoa(m.Sender.ID)
                subscribers := string(val)
                if subscribers == "" {
                    return txn.Set(btopics, []byte(sender))
                }
                arrSubscribers := strings.Split(subscribers, ",")
                updSubscribers := strings.Join(remove(arrSubscribers, subscribers), ",")
                return txn.Set(btopics, []byte(updSubscribers))
            }); err != nil {
                return err
            }
            return nil
        })
        if err == badger.ErrKeyNotFound {
            t.send(m, fmt.Sprintf(constant.SubscribeIdNotFound, topics))
            return
        }
        if err != nil {
            log.Println(err)
            t.send(m, fmt.Sprintf(constant.SubscribeIdFailed, topics))
        }

    })
}

func (t *Bot) send(m *tb.Message, message string) {
    m, err := t.B.Send(m.Sender, message)
    if err != nil {
        log.Println(err)
    }
}

func remove(l []string, item string) []string {
    for i, other := range l {
        if other == item {
            return append(l[:i], l[i+1:]...)
        }
    }
}

func (t *Bot) Run() {
    t.OnCreate()
    t.OnHelp()
    t.OnSubscribe()
    t.OnUnsubscribe()

    go t.B.Start()
}
