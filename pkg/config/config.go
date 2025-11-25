package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type (
	HttpServerConfig struct {
		Host           string
		Port           string
		ReadTimeout    time.Duration
		WriteTimeout   time.Duration
		IdleTimeout    time.Duration
		MaxHeaderBytes int
	}

	PostgresConfig struct {
		Host     string
		Port     string
		Username string
		Password string
		DBName   string
		SSLMode  string
	}

	JwtConfig struct {
		Secret   string
		TokenTTL time.Duration
	}

	HashConfig struct {
		Cost uint
	}

	Config struct {
		Env        string
		HttpServer HttpServerConfig
		Postgres   PostgresConfig
		Jwt        JwtConfig
		Hash       HashConfig
	}
)

func MustLoad() *Config {
	v := viper.New()

	v.SetConfigFile(".env")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	v.SetDefault("ENV", "development")
	v.SetDefault("HTTP_HOST", "localhost")
	v.SetDefault("HTTP_PORT", "8080")
	v.SetDefault("HTTP_READ_TIMEOUT", "10s")
	v.SetDefault("HTTP_WRITE_TIMEOUT", "10s")
	v.SetDefault("HTTP_IDLE_TIMEOUT", "120s")
	v.SetDefault("HTTP_MAX_HEADER_BYTES", 1)
	v.SetDefault("JWT_TOKEN_TTL", "24h")
	v.SetDefault("HASH_COST", 12)

	readDuration := func(key string) time.Duration {
		d, err := time.ParseDuration(v.GetString(key))
		if err != nil {
			log.Fatalf("invalid duration for %s: %v", key, err)
		}
		return d
	}

	cfg := &Config{
		Env: v.GetString("ENV"),
		HttpServer: HttpServerConfig{
			Host:           v.GetString("HTTP_HOST"),
			Port:           v.GetString("HTTP_PORT"),
			ReadTimeout:    readDuration("HTTP_READ_TIMEOUT"),
			WriteTimeout:   readDuration("HTTP_WRITE_TIMEOUT"),
			IdleTimeout:    readDuration("HTTP_IDLE_TIMEOUT"),
			MaxHeaderBytes: v.GetInt("HTTP_MAX_HEADER_BYTES"),
		},
		Postgres: PostgresConfig{
			Host:     v.GetString("POSTGRES_HOST"),
			Port:     v.GetString("POSTGRES_PORT"),
			Username: v.GetString("POSTGRES_USER"),
			Password: v.GetString("POSTGRES_PASSWORD"),
			DBName:   v.GetString("POSTGRES_DBNAME"),
			SSLMode:  v.GetString("POSTGRES_SSLMODE"),
		},
		Jwt: JwtConfig{
			Secret:   v.GetString("JWT_SECRET"),
			TokenTTL: readDuration("JWT_TOKEN_TTL"),
		},
		Hash: HashConfig{
			Cost: v.GetUint("HASH_COST"),
		},
	}

	return cfg
}
