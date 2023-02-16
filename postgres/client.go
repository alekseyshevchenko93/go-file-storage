package postgres

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type PostgresClient interface {
	Shutdown() error
	GetClient() *sqlx.DB
}

type postgresClient struct {
	db  *sqlx.DB
	log *logrus.Logger
}

func (c *postgresClient) GetClient() *sqlx.DB {
	return c.db
}

func (c *postgresClient) Shutdown() error {
	err := c.db.Close()

	if err != nil {
		return fmt.Errorf("failed to shutdown database %w", err)
	}

	c.log.Info("db.shutdown")

	return nil
}

func New(log *logrus.Logger) (*postgresClient, error) {
	host := os.Getenv("POSTGRES_HOST")
	database := os.Getenv("POSTGRES_DATABASE")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")

	dsn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=disable", host, user, database, password)
	db, err := sqlx.Connect("postgres", dsn)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to db %w", err)
	}

	log.Info("db.connected")

	return &postgresClient{
		db,
		log,
	}, nil
}
