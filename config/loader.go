package config

type (
    Config struct {
        AppUri          string
        HttpPort        int
        DbPath          string
        TelegramToken   string
        TelegramPollInt int
    }
)

func NewConfig() *Config {
    return &Config{
        AppUri:          GetString("APP_URI"),
        HttpPort:        GetInt("HTTP_PORT"),
        DbPath:          GetString("DB_PATH"),
        TelegramToken:   GetString("TELEGRAM_TOKEN"),
        TelegramPollInt: GetInt("TELEGRAM_POLLING_INTERVAL"),
    }
}
