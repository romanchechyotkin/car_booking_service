package postgresql

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"

	"context"
	"log"
	"strings"
)

type pgConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

func NewPgConfig(username, password, host, port, database string) *pgConfig {
	return &pgConfig{
		Username: username,
		Password: password,
		Host:     host,
		Port:     port,
		Database: database,
	}
}

func NewClient(ctx context.Context, cfg *pgConfig) *pgxpool.Pool {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	//dockerConn := "postgresql://postgres:5432@database:5432/car_booking_service?sslmode=disable"
	log.Println(connString)

	log.Println("postgresql client init")
	pool, err := pgxpool.New(ctx, connString)
	err = pool.Ping(ctx)
	if err != nil {
		log.Println(err)
		log.Fatal("cannot to connect to postgres")
	}

	return pool
}

func FormatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", "")
}
