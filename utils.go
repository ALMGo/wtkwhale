package main

import (
	"fmt"
	"github.com/almaclaine/ethplorer"
	"github.com/almaclaine/whalenamegenerator"
	"github.com/jmoiron/sqlx"
)

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
			Limit: 1000,
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
