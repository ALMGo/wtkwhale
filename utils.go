package main

import (
	"fmt"
	"github.com/almaclaine/ethplorer"
	"github.com/almaclaine/whalenamegenerator"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/jmoiron/sqlx"
	"strconv"
)

var wellKnownWhales = map[string]string {
	"0x68b22215ff74e3606bd5e6c1de8c2d68180c85f7": "Bitmart",
	"0x6cc8dcbca746a6e4fdefb98e1d0df903b107fd21": "Bittrue",
}

func getWhaleName(conn *sqlx.DB) (string, error) {
	name := whalenamegenerator.GenerateWhaleName()
	for {
		exists, err := ExistsWhaleName(conn, name)
		if err != nil {
			return "", err
		}
		if !exists {
			return name, nil
		}
		name = whalenamegenerator.RandomWhale()
	}
}

func updateTopHolders(conn *sqlx.DB) (*ethplorer.TopTokenHolders, error) {
	holders, err := ethplorer.GetTopTokenHolders(
		WTK_ADDRESS,
		ethplorer.GetTopTokenHoldersParams{
			Limit: 50,
		},
		ETHPLORER_KEY,
	)
	if err != nil {
		return &ethplorer.TopTokenHolders{}, err
	}
	addTopHoldersToDB(conn, holders)
	return holders, nil
}

func addTopHoldersToDB(conn *sqlx.DB, holders *ethplorer.TopTokenHolders) error {
	for _, holder := range holders.Holders {
		address := holder.Address
		exists, err := ExistsWhaleAddress(conn, address)
		if err != nil {
			return err
		}
		if !exists {
			err := CreateWhale(conn, address)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func findHolder(holders *ethplorer.TopTokenHolders, address string) int {
	for index, element := range holders.Holders {
		if element.Address == address {
			return index
		}
	}
	return -1
}

func setupConfig(conn *sqlx.DB) (Config, error) {
	exists, err := ExistsConfig(conn)
	if err != nil {
		fmt.Println(err)
	}

	if !exists {
		err = CreateConfig(conn)
		if err != nil {
			return Config{}, err
		}
	}

	config, err := GetConfig(conn)

	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func sendTweet(client *twitter.Client, tweet string) error {
	_, _, err := client.Statuses.Update(tweet, nil)
	return err
}

func getTransactions() (*ethplorer.TokenHistory, error) {
	return ethplorer.GetTokenHistory(
		WTK_ADDRESS,
		ethplorer.GetTokenHistoryParams{
			Type:  "transfer",
			Limit: 1000,
		},
		ETHPLORER_KEY)
}

func watchWhales(conn *sqlx.DB) {
	transactions, err := getTransactions()
	if err != nil {
		logger.Error(err.Error())
	}

	for _, element := range transactions.Operations {
		if element.Timestamp.UnixNano() > config.LastUpdated.UnixNano() {

			existsFrom, err := ExistsWhaleAddress(conn, element.From)
			if err != nil {
				logger.Error(err.Error())
			}
			tweet := ""
			if existsFrom {
				whaleFrom, err := GetWhaleByAddress(conn, element.From)
				if err != nil {
					logger.Error(err.Error())
				}
				tweet = whaleFrom.Name + " Whale"
			} else if element.Value/100 > 500000 {
				whaleFrom, err := GetWhaleByAddress(conn, element.From)
				if err != nil {
					logger.Error(err.Error())
				}
				tweet = whaleFrom.Name + "Address: " + element.From
			} else {
				continue
			}

			tweet += " sent "
			tweet += strconv.FormatFloat(float64(element.Value)/100, 'f', 2, 32)
			tweet += " WTK To "
			if err != nil {
				logger.Error(err.Error())
			}

			existsTo, err := ExistsWhaleAddress(conn, element.To)
			if err != nil {
				logger.Error(err.Error())
			}
			if existsTo {
				whaleFrom, err := GetWhaleByAddress(conn, element.To)
				if err != nil {
					logger.Error(err.Error())
				}
				tweet += whaleFrom.Name + " Whale"
			} else {
				tweet += "Address: " + element.To
			}
			tweet += "\nhttps://etherscan.io/tx/" + element.TransactionHash
			tweet += "\n$WTK #WTK $XDC #XDC"
			tweets = append(tweets, tweet)
		}
	}
	config.LastUpdated = transactions.Operations[0].Timestamp.Time
	err = UpdateConfigLastUpdated(conn, config.LastUpdated)
	if err != nil {
		logger.Error(err.Error())
	}
}
