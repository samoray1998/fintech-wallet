package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func LoadConfig() *Config {

	_ = godotenv.Load()

	return &Config{
		Server: ServerConfig{
			Port:      getEnv("PORT", DefaultPort),
			Env:       getEnv("Env", DefaultEnv),
			TimeOut:   parseDuration(getEnv("SERVER_TIMEOUT", "30s")),
			RateLimit: getEnvAsInt("RATE_LIMIT", DefaultRateLimit),
			Debug:     getEnvAsBool("DEBUG", false),
		},
		Database: DatabaseConfig{
			Uri:            getMongoURI(),
			Name:           getEnv("DB_NAME", DefaultDBName),
			MaxPoolSize:    getEnvAsUint64("DB_MAX_POOL", DefaultMaxPoolSize),
			MinPoolSize:    getEnvAsUint64("DB_MIN_POOL", DefaultMinPoolSize),
			ConnectTimeout: parseDuration(getEnv("DB_CONNECT_TIMEOUT", DefaultConnectTimeout.String())),
			SocketTimeout:  parseDuration(getEnv("DB_SOCKET_TIMEOUT", DefaultSocketTimeout.String())),
		},
		Auth: AuthConfig{
			JWTSecret:        getEnv("JWT_SECRET", ""),
			JWTAccessExpiry:  parseDuration(getEnv("JWT_ACCESS_EXPIRY", "15m")),
			JWTRefreshExpiry: parseDuration(getEnv("JWT_REFRESH_EXPIRY", DefaultJWTExpiry.String())),
			BcryptCost:       getEnvAsInt("BCRYPT_COST", DefaultBcryptCost),
		},
		KYC: KYCConfig{
			VerifyURL:     getEnv("KYC_VERIFY_URL", DefaultKYCVerifyURL),
			APITimeout:    parseDuration(getEnv("KYC_TIMEOUT", "10s")),
			MaxRetries:    getEnvAsInt("KYC_MAX_RETRIES", 3),
			WebhookSecret: getEnv("KYC_WEBHOOK_SECRET", ""),
		},
		Rates: RatesConfig{
			BaseCurrency:   getEnv("BASE_CURRENCY", "USD"),
			ExchangeAPIURL: getEnv("EXCHANGE_API_URL", ""),
			CacheDuration:  parseDuration(getEnv("RATES_CACHE_DURATION", "1h")),
			APIKey:        getEnv("EXCHANGE_API_KEY", ""),
		},
	}
}

func getMongoURI() string {
	uri := getEnv("MONGO_URI", DefaultMongoURI)

	// Ensure the URI has the minimum required parameters
	if !strings.Contains(uri, "?") {
		uri += "?"
	} else if !strings.HasSuffix(uri, "&") {
		uri += "&"
	}

	// Add critical parameters if not present
	if !strings.Contains(uri, "directConnection=") {
		uri += "directConnection=true&"
	}
	if !strings.Contains(uri, "serverSelectionTimeoutMS=") {
		uri += "serverSelectionTimeoutMS=2000&"
	}

	return strings.TrimSuffix(uri, "&")
}

func getEnv(key string, Defaultval string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	if Defaultval == "" {
		log.Printf("WARNING: Required environment variable %s not set", key)
	}
	return Defaultval
}

func getEnvAsBool(key string, defaultValue bool) bool {
	strValue := getEnv(key, "")
	if strValue == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(strValue)
	if err != nil {
		log.Printf("Invalid boolean value for %s, using default: %v", key, defaultValue)
		return defaultValue
	}
	return value
}
func getEnvAsInt(key string, defaultValue int) int {
	strValue := getEnv(key, "")
	if strValue == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(strValue)
	if err != nil {
		log.Printf("Invalid value for %s, using default: %v", key, defaultValue)
		return defaultValue
	}
	return value
}

func getEnvAsUint64(key string, defaultValue uint64) uint64 {
	strValue := getEnv(key, "")
	if strValue == "" {
		return defaultValue
	}
	value, err := strconv.ParseUint(strValue, 10, 64)
	if err != nil {
		log.Printf("Invalid value for %s, using default: %v", key, defaultValue)
		return defaultValue
	}
	return value
}

func parseDuration(durationStr string) time.Duration {
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		log.Printf("Invalid duration format for %s, defaulting to 0", durationStr)
		return 0
	}
	return duration
}
