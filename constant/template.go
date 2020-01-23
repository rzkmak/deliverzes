package constant

const (
    HelpCommand        = "/help"
    SubscribeCommand   = "/subscribe"
    UnsubscribeCommand = "/unsubscribe"

    HelpMessage = "available command, \n" +
        "1. /subscribe <your-subscribe-id> to subscribe specific topic.\n" +
        "2. /unsubscribe <topic> to unsubscribe topic.\n" +
        "3./help to display all command"

    UnknownMessageReply = "%v is unknown command, please see /help to view available one"

    SubscribeParameterRequired = "subscribe need subscribe_id to used, example /subscribe tcp-123-ack"
    UnsubscribeParameterRequired = "unsubscribe need subscribe_id to used, example /subscribe tcp-123-ack"

    SubscribeIdNotFound = "subscribe id %v not found"
    SubscribeIdFailed = "subscribe id %v failed, internal server error"
)
