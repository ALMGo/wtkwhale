package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"math/rand"
	"os"
	"time"
)

var logger *zap.Logger

func main() {
	rand.Seed(time.Now().UnixNano())
	db, err := sqlx.Connect("sqlite3", "test.db")
	logger, _ = zap.NewProduction()

	if err != nil {
		logger.Info("error Connecting To Database")
		os.Exit(1)
	}

	err = CreateWhale(db, "0x0000000000000000000000000000000000000000")
	err = UpdateWhaleLastSent(db, "0x0000000000000000000000000000000000000000", time.Now())
	if err != nil {
		logger.Info(err.Error())
		os.Exit(1)
	}
}
