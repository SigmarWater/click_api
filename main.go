package main

import (
	"context"
	"github.com/ClickHouse/clickhouse-go/v2"
	"log"
	"net"
	"os"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/SigmarWater/click_api/scripts"
	"github.com/joho/godotenv"
)

type Config struct {
	Host     string
	Port     string
	DBName   string
	User     string
	Password string
}

func mustHaveEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("missing env: %s", key)
	}
	return value
}

func loadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &Config{
		Host:     mustHaveEnv("DB_HOST"),
		Port:     mustHaveEnv("DB_PORT"),
		DBName:   mustHaveEnv("DB_NAME"),
		User:     mustHaveEnv("DB_USER"),
		Password: mustHaveEnv("DB_PASSWORD"),
	}, nil
}

func connect(ctx context.Context, config *Config) (driver.Conn, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{net.JoinHostPort(config.Host, config.Port)},
		Auth: clickhouse.Auth{
			Database: config.DBName,
			Username: config.User,
			Password: config.Password,
		},
	})

	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			log.Printf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		return nil, err
	}
	return conn, nil
}

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("failed loading config: %v", err)
	}

	connCtx, connCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer connCancel()

	conn, err := connect(connCtx, config)
	if err != nil {
		log.Fatalf("failed connection: %v", err)
	}

	defer func() {
		_ = conn.Close()
	}()

	queryCtx, queryCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer queryCancel()

	err = scripts.InsertBatchToDb(queryCtx, conn)
	if err != nil {
		log.Println("failed insert dates into db")
		return
	}
	log.Println("success insert dates into db")
}
