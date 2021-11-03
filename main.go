package main

import (
	"fmt"
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
var config Config

func main() {
	rand.Seed(time.Now().UnixNano())
	db, err := sqlx.Connect("sqlite3", "wtkwhales.db")
	topHolders, _ = updateTopHolders(db)
	logger, _ = zap.NewProduction()

	if err != nil {
		logger.Error("error Connecting To Database")
		os.Exit(1)
	}

	config, err = setupConfig(db)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	configTime := uint64(config.LastUpdated.UnixNano() / 1000000000)
	transactions, err := ethplorer.GetTokenHistory(
		WTK_ADDRESS,
		ethplorer.GetTokenHistoryParams{
			Type: "transfer",
			Limit: 5,
		},
		ETHPLORER_KEY)
	for _, element := range transactions.Operations {
		if element.Timestamp > configTime {
			whaleTo, err := GetWhaleByAddress(db, element.To)
			if err != nil {
				logger.Error(err.Error())
			}
			whaleFrom, err := GetWhaleByAddress(db, element.From)
			if err != nil {
				logger.Error(err.Error())
			}
			fmt.Println("------------------------------")
			fmt.Println("Transaction: ", element.TransactionHash)
			fmt.Println("To Whale: ", whaleTo.Name, element.To)
			fmt.Println("From Whale: ", whaleFrom.Name, element.From)
			fmt.Println("------------------------------")
		}
	}
	config.LastUpdated = time.Unix(int64(transactions.Operations[0].Timestamp), 0)
	err = UpdateConfigLastUpdated(db, config.LastUpdated)
	if err != nil {
		logger.Error(err.Error())
	}

	//fmt.Println(transactions.Operations)
	//err = CreateWhale(db, "0x0000000000000000000000000000000000000000")
	//err = UpdateWhaleLastSent(db, "0x0000000000000000000000000000000000000000", time.Now())
	//if err != nil {
	//	logger.Info(err.Error())
	//	os.Exit(1)
	//}
}
