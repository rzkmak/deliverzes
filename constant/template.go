package constant

const (
    StartCommand       = "/start"
    HelpCommand        = "/help"
    SubscribeCommand   = "/subscribe"
    UnsubscribeCommand = "/unsubscribe"

    HelpMessage = "available command: \n" +
        "1. /subscribe <your-subscribe-id> to subscribe specific topic.\n" +
        "2. /unsubscribe <topic> to unsubscribe topic.\n" +
        "3./help to display all command\n" +
        "Remember that I didn't want to reply your group :p, just reach me through personal chat okay"

    UnknownMessageReply = "%v is unknown command, please see /help to view available one"
    SubscribeIdSuccess  = "%v subscribe id success"
    UnsubscribeIdSuccess  = "%v unsubscribe id success"

    SubscribeParameterRequired   = "subscribe need subscribe_id to used, example /subscribe tcp-123-ack"
    UnsubscribeParameterRequired = "unsubscribe need subscribe_id to used, example /subscribe tcp-123-ack"

    SubscribeIdNotFound = "subscribe id %v not found"
    SubscribeIdFailed   = "subscribe id %v failed, internal server error"

    InternalServerError = "internal server error"

    ProvideCorrectJson           = "please provide correct json request based on docs"
    MinimumSubscriberName    int = 6
    MinimumSubscriberNameMsg     = "minimum subscriber name is 6"

    BlankHookUrl = "please provide valid hook url"
)
