package postgresql

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/romanchechyotkin/car_booking_service/schema"
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


func Migrate(dbUrl string) {
	source, err := iofs.New(schema.DB, "migrations")
	if err != nil {
		log.Println("failed to read migrations source", err)
		return
	}

	instance, err := migrate.NewWithSourceInstance("iofs", source, makeMigrateUrl(dbUrl))
	if err != nil {
		log.Println("failed to initialization the migrate instance", err)
		return
	}

	err = instance.Up()

	switch err {
	case nil:
		log.Println("the migration schema successfully upgraded!")
	case migrate.ErrNoChange:
		log.Println("the migration schema not changed")
	default:
		log.Println("could not apply the migration schema", err)
	}
}

func makeMigrateUrl(dbUrl string) string {
	urlRe := regexp.MustCompile("^[^\\?]+")
	url := urlRe.FindString(dbUrl)

	sslModeRe := regexp.MustCompile("(sslmode=)[a-zA-Z0-9]+")
	sslMode := sslModeRe.FindString(dbUrl)

	return fmt.Sprintf("%s?%s", url, sslMode)
}
