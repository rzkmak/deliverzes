package telegram

import "C"
import (
    "crypto/rand"
    "encoding/json"
    "fmt"
    "github.com/aeidelos/deliverzes/config"
    "github.com/aeidelos/deliverzes/constant"
    tb "github.com/demget/telebot"
    "github.com/dgraph-io/badger"
    "io/ioutil"
    "log"
    "net/http"
    "strconv"
    "strings"
)

type Bot struct {
    B *tb.Bot
    C *config.Config
    D *badger.DB
}

func NewBot(b *tb.Bot, c *config.Config, d *badger.DB) *Bot {
    return &Bot{B: b, C: c, D: d}
}

func (t *Bot) OnCreate() {
    log.Println("starting telegram bot..")
}

func (t *Bot) OnHelp() {
    t.B.Handle(constant.StartCommand, func(m *tb.Message) {
        t.send(m, constant.HelpMessage)
    })
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
        if m.FromGroup() {
            return
        }
        if m.FromChannel() {
            return
        }
        split := strings.Split(strings.TrimSpace(m.Text), " ")
        if len(split) < 2 {
            t.send(m, constant.SubscribeParameterRequired)
            return
        }
        topics := split[1]
        err := t.D.Update(func(txn *badger.Txn) error {
            bTopics := []byte(topics)
            item, err := txn.Get(bTopics)
            if err != nil {
                return err
            }
            if err := item.Value(func(val []byte) error {
                sender := strconv.Itoa(m.Sender.ID)
                subscribers := string(val)
                if subscribers == "" {
                    return txn.Set(bTopics, []byte(sender))
                }
                arrSubscribers := strings.Split(subscribers, ",")
                remSubscriber := remove(arrSubscribers, sender)
                updSubscribers := strings.Join(append(remSubscriber, sender), ",")
                return txn.Set(bTopics, []byte(updSubscribers))
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
            return
        }
        t.send(m, fmt.Sprintf(constant.SubscribeIdSuccess, topics))
    })
}

func (t *Bot) OnUnsubscribe() {
    t.B.Handle(constant.UnsubscribeCommand, func(m *tb.Message) {
        split := strings.Split(strings.TrimSpace(m.Text), " ")
        if len(split) < 2 {
            t.send(m, constant.UnsubscribeCommand)
            return
        }
        topics := split[1]
        err := t.D.Update(func(txn *badger.Txn) error {
            bTopics := []byte(topics)
            item, err := txn.Get(bTopics)
            if err != nil {
                return err
            }
            if err := item.Value(func(val []byte) error {
                subscribers := string(val)
                if subscribers == "" {
                    return nil
                }
                arrSubscribers := strings.Split(subscribers, ",")
                updSubscribers := strings.Join(remove(arrSubscribers, subscribers), ",")
                return txn.Set(bTopics, []byte(updSubscribers))
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
            return
        }
        t.send(m, fmt.Sprintf(constant.UnsubscribeIdSuccess, topics))
    })
}

func (t *Bot) send(m *tb.Message, message string) {
    m, err := t.B.Send(m.Sender, message)
    if err != nil {
        log.Println(err)
    }
}

func (t *Bot) sendToUser(m *tb.User, message string) {
    _, err := t.B.Send(m, message)
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
    return l
}

func (t *Bot) Run() {
    t.OnCreate()
    t.OnHelp()
    t.OnSubscribe()
    t.OnUnsubscribe()

    go t.B.Start()
}

type (
    ReqSubIdParam struct {
        SubName string `json:"subscriber_name"`
    }
    ReqSubIdResp struct {
        Message string `json:"message"`
        Status  bool   `json:"status"`
        SubId   string `subscriber_id`
        Url     string `webhook_url`
    }
    ReqHookParam struct {
        Title  string                 `json:"title"`
        Origin string                 `json:"origin"`
        Value  map[string]interface{} `json:"value"`
    }
    ResHookResp struct {
        Status  bool   `json:"status"`
        Message string `json:"message"`
    }
)

func (t *Bot) GenerateSubscriberIdHandler(w http.ResponseWriter, r *http.Request) {
    decoder := json.NewDecoder(r.Body)
    var body ReqSubIdParam
    if err := decoder.Decode(&body); err != nil {
        res := ReqSubIdResp{
            Status:  false,
            Message: constant.ProvideCorrectJson,
        }
        resJson, err := json.Marshal(res)
        if err != nil {
            _, _ = w.Write([]byte(constant.InternalServerError))
        }
        _, _ = w.Write(resJson)
        return
    }

    if len(body.SubName) < constant.MinimumSubscriberName {
        res := ReqSubIdResp{
            Status:  false,
            Message: constant.MinimumSubscriberNameMsg,
        }
        resJson, err := json.Marshal(res)
        if err != nil {
            _, _ = w.Write([]byte(constant.InternalServerError))
        }
        _, _ = w.Write(resJson)
        return
    }
    topics := body.SubName + "-" + randString()
    if err := t.D.Update(func(txn *badger.Txn) error {
        return txn.Set([]byte(topics), []byte(""))
    }); err != nil {
        res := ReqSubIdResp{
            Status:  false,
            Message: err.Error(),
        }
        resJson, err := json.Marshal(res)
        if err != nil {
            _, _ = w.Write([]byte(constant.InternalServerError))
        }
        _, _ = w.Write(resJson)
        return
    }
    res := ReqSubIdResp{
        Status:  true,
        Message: "",
        SubId:   topics,
        Url:     fmt.Sprintf("http://%v:%v/send?hook_url=%v" ,t.C.AppUri, t.C.HttpPort, topics),
    }
    resJson, err := json.Marshal(res)
    if err != nil {
        _, _ = w.Write([]byte(constant.InternalServerError))
    }
    _, _ = w.Write(resJson)
}

func (t *Bot) SendMessageHandler(w http.ResponseWriter, r *http.Request) {

    if r.URL.Query().Get("hook_url") == "" {
        res := ReqSubIdResp{
            Status:  false,
            Message: constant.BlankHookUrl,
        }
        resJson, err := json.Marshal(res)
        if err != nil {
            _, _ = w.Write([]byte(constant.InternalServerError))
        }
        _, _ = w.Write(resJson)
        return
    }

    topics := r.URL.Query().Get("hook_url")

    err := t.D.View(func(txn *badger.Txn) error {
        bTopics := []byte(topics)
        item, err := txn.Get(bTopics)
        if err != nil {
            return err
        }
        if err := item.Value(func(val []byte) error {
            subscribers := string(val)
            if subscribers == "" {
                return nil
            }
            body, err := ioutil.ReadAll(r.Body)
            if err != nil {
                return err
            }

            arrSubscribers := strings.Split(subscribers, ",")
            for _, sub := range arrSubscribers {
                user, err := strconv.Atoi(sub)
                if err != nil {
                    continue
                }
                t.sendToUser(&tb.User{
                    ID: user,
                }, fmt.Sprintf("%s", body))
            }
            return nil

        }); err != nil {
            return err
        }
        return nil
    })
    res := ResHookResp{
        Status: true,
    }
    if err == badger.ErrKeyNotFound {
        res.Message = fmt.Sprintf(constant.SubscribeIdNotFound, topics)
        res.Status = false
    }
    if err != nil {
        res.Message = constant.InternalServerError
        res.Status = false
    }
    resJson, err := json.Marshal(res)
    if err != nil {
        _, _ = w.Write([]byte(constant.InternalServerError))
    }
    _, _ = w.Write(resJson)
}

func randString() string {
    n := 5
    b := make([]byte, n)
    if _, err := rand.Read(b); err != nil {
        panic(err)
    }
    return fmt.Sprintf("%X", b)
}
