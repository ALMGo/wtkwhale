package main

import "time"

type Whale struct {
	Id int64 `db:"id"`
	Address string `db:"address"`
	Name string `db:"name"`
	Added time.Time `db:"added"`
	LastSent time.Time `db:"last_sent"`
	LastReceived time.Time `db:"last_received"`
}

type Config struct {
	Id int64 `db:"id"`
	LastUpdated time.Time `db:"last_updated"`
}
