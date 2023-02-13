package postgres

import (
	"fmt"
	"os"

	log "github.com/alexshv/file-storage/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type DatabaseClient struct {
	db *sqlx.DB
}

var client *sqlx.DB

func Init() {
	host := os.Getenv("POSTGRES_HOST")
	database := os.Getenv("POSTGRES_DATABASE")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")

	dsn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=disable", host, user, database, password)
	db, err := sqlx.Connect("postgres", dsn)

	if err != nil {
		log.GetLogger().WithFields(logrus.Fields{
			"message": err,
		}).Error("postgres.connection.error")

		panic(err)
	}

	log.GetLogger().Info("postgres.connected")

	client = db
}

func Shutdown() {
	err := client.Close()

	if err != nil {
		log.GetLogger().WithFields(logrus.Fields{
			"message": err,
		}).Error("postgres.shutdown.error")
		return
	}

	log.GetLogger().Info("postgres.shutdown.sucess")
}

func GetClient() *sqlx.DB {
	return client
}
