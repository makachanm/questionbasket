package config

import "os"

// Config holds the application's configuration.
type Config struct {
	DatabaseURL         string
	MastodonInstanceURL string
	MastodonToken       string
	ServiceURL          string
}

// Cfg is a global variable holding the loaded configuration.
var Cfg *Config

// LoadConfig loads configuration from environment variables into the Cfg global variable.
func LoadConfig() {
	Cfg = &Config{
		DatabaseURL:         os.Getenv("DATABASE_URL"),
		MastodonInstanceURL: os.Getenv("MASTODON_INSTANCE_URL"),
		MastodonToken:       os.Getenv("MASTODON_TOKEN"),
		ServiceURL:          os.Getenv("SERVICE_URL"),
	}
}