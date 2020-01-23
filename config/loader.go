package config

type (
    Config struct {
        HttpPort        int
        TelegramToken   string
        TelegramPollInt int
    }
)

func NewConfig() *Config {
    return &Config{
        HttpPort:        GetInt("HTTP_PORT"),
        TelegramToken:   GetString("TELEGRAM_TOKEN"),
        TelegramPollInt: GetInt("TELEGRAM_POLLING_INTERVAL"),
    }
}
