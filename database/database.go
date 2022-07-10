package database

import (
	"context"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/database"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/logrusadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

var Client *database.Database

func ConnectToDatabase() {
	config, err := pgxpool.ParseConfig(config.Conf.Database.Uri)
	if err != nil {
		panic(err)
	}

	// TODO: Sentry
	config.ConnConfig.LogLevel = pgx.LogLevelWarn
	config.ConnConfig.Logger = logrusadapter.NewLogger(logrus.New())

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		panic(err)
	}

	Client = database.NewDatabase(pool)
	Client.CreateTables(pool)
}
