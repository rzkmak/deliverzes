package config

type (
    Config struct {
        HttpPort int
    }
)

func NewConfig() *Config {
    return &Config{HttpPort: GetInt("HTTP_PORT")}
}
