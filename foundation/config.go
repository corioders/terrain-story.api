package foundation

import (
	"flag"
)

// LoadConfig loads config from flags provided to the process.
func LoadConfig() (*Config, error) {
	config := &Config{}

	flag.StringVar(&config.Web.Port, "port", "8080", "Server port.")

	flag.StringVar(&config.Qr.GamesCodeJsonPath, "gamesCode", "../data/gamesCode.jsonc", "Games code json path.")

	flag.Parse()

	err := config.validate()
	if err != nil {
		return nil, err
	}

	return config, nil
}

type Config struct {
	Web WebConfig
	Qr  QrConfig
}

func (c *Config) validate() error {
	if err := c.Web.validate(); err != nil {
		return err
	}

	return nil
}

type WebConfig struct {
	Port string
}

func (c *WebConfig) validate() error {
	return nil
}

type QrConfig struct {
	GamesCodeJsonPath string
}
