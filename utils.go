package main

import (
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