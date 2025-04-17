package config

import "time"

const (
	DefaultPort            = "8080"
	DefaultEnv             = "development"
	DefaultDBName          = "fintech"
	DefaultMongoURI        = "mongodb://localhost:27017"
	DefaultJWTExpiry       = 24 * time.Hour
	DefaultBcryptCost      = 10
	DefaultRateLimit       = 100
	DefaultKYCVerifyURL    = "https://kyc-service.example.com"
	DefaultConnectTimeout  = 5 * time.Second
	DefaultSocketTimeout   = 30 * time.Second
	DefaultMaxPoolSize     = 50
	DefaultMinPoolSize     = 10
)