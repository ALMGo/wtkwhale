package main

import (
	"github.com/almaclaine/ethplorer"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"math/rand"
	"os"
	"time"
)

var logger *zap.Logger

var WTK_ADDRESS = "0xdf9d4674a430bdcc096a3a403128357ab36844ba"
var ETHPLORER_KEY = os.Getenv("ETHPLORER_KEY")

var topHolders *ethplorer.TopTokenHolders

func main() {
	rand.Seed(time.Now().UnixNano())
	db, err := sqlx.Connect("sqlite3", "wtkwhales.db")
	topHolders, _ = updateTopHolders(db)
	logger, _ = zap.NewProduction()

	if err != nil {
		logger.Info("error Connecting To Database")
		os.Exit(1)
	}

	//err = CreateWhale(db, "0x0000000000000000000000000000000000000000")
	//err = UpdateWhaleLastSent(db, "0x0000000000000000000000000000000000000000", time.Now())
	//if err != nil {
	//	logger.Info(err.Error())
	//	os.Exit(1)
	//}
}
