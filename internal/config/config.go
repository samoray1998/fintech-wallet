package config

import "time"

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
	KYC      KYCConfig
	Rates    RatesConfig
}

type ServerConfig struct {
	Port      string
	Env       string
	TimeOut   time.Duration
	RateLimit int
	Debug     bool
}

type DatabaseConfig struct {
	Uri            string
	Name           string
	MaxPoolSize    uint64
	MinPoolSize    uint64
	ConnectTimeout time.Duration
	SocketTimeout  time.Duration
}

type AuthConfig struct {
	JWTSecret        string
	JWTAccessExpiry  time.Duration
	JWTRefreshExpiry time.Duration
	BcryptCost       int
}

type KYCConfig struct {
	VerifyURL     string
	APITimeout    time.Duration
	MaxRetries    int
	WebhookSecret string
}

type RatesConfig struct {
	BaseCurrency   string
	ExchangeAPIURL string
	CacheDuration  time.Duration
	APIKey         string
}
