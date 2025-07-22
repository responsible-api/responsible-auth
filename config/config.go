package config

import (
	"log"
	"time"

	"github.com/joeshaw/envdecode"
	"github.com/joho/godotenv"
)

type Conf struct {
	Server ConfServer
	DB     ConfDB
}

type ConfServer struct {
	Port         int           `env:"SERVER_PORT,default=8080"`
	TimeoutRead  time.Duration `env:"SERVER_TIMEOUT_READ,default=5s"`
	TimeoutWrite time.Duration `env:"SERVER_TIMEOUT_WRITE,default=5s"`
	TimeoutIdle  time.Duration `env:"SERVER_TIMEOUT_IDLE,default=5s"`
	Debug        bool          `env:"SERVER_DEBUG,default=false"`
}

type ConfDB struct {
	Host     string `env:"DB_HOST,default="`
	Port     int    `env:"DB_PORT,default=0"`
	Username string `env:"DB_USER,default="`
	Password string `env:"DB_PASS,default="`
	DBName   string `env:"DB_NAME,default="`
	Debug    bool   `env:"DB_DEBUG,default="`
}

func Config() *Conf {
	// Load .env file if present
	_ = godotenv.Load()

	var c Conf
	if err := envdecode.StrictDecode(&c); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}

	return &c
}

func ConfigDB() *ConfDB {
	// Load .env file if present
	_ = godotenv.Load()

	var c ConfDB
	if err := envdecode.StrictDecode(&c); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}

	return &c
}
