package config

import (
	"log"
	"time"

	"github.com/joeshaw/envdecode"
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
	Host     string `env:"DB_HOST,default=responsible-api-db"`
	Port     int    `env:"DB_PORT,default=3306"`
	Username string `env:"DB_USER,default=responsible_api_user"`
	Password string `env:"DB_PASS,default=responsible_api_pass"`
	DBName   string `env:"DB_NAME,default=responsible_api"`
	Debug    bool   `env:"DB_DEBUG,default=false"`
}

func New() *Conf {
	var c Conf
	if err := envdecode.StrictDecode(&c); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}

	return &c
}

func NewDB() *ConfDB {
	var c ConfDB
	if err := envdecode.StrictDecode(&c); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}

	return &c
}
