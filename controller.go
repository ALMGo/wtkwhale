package main

import (
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"regexp"
	"time"
)

func ExistsWhaleName(conn *sqlx.DB, name string) (bool, error) {
	sql, args, err := squirrel.Select("1").
		Prefix("SELECT EXISTS (").
		From("whale").
		Where(squirrel.Eq{"name": name}).
		Suffix(")").
		ToSql()

	if err != nil {
		return false, err
	}

	var exists bool
	err = conn.QueryRow(sql, args...).Scan(&exists)

	if err != nil {
		return false, err
	}
	return exists, nil
}

func ExistsWhaleAddress(conn *sqlx.DB, address string) (bool, error) {
	sql, args, err := squirrel.Select("1").
		Prefix("SELECT EXISTS (").
		From("whale").
		Where(squirrel.Eq{"address": address}).
		Suffix(")").
		ToSql()

	if err != nil {
		return false, err
	}

	var exists bool
	err = conn.QueryRow(sql, args...).Scan(&exists)

	if err != nil {
		return false, err
	}
	return exists, nil
}


func CreateWhale(conn *sqlx.DB, address string) error {
	match, err := regexp.MatchString("^0x[a-fA-F0-9]{40}$", address)
	if err != nil {
		return err
	}

	if !match {
		return errors.New("invalid ethereum address")
	}

	name := ""
	if val, ok := wellKnownWhales[address]; ok {
		name = val
	} else {
		name, err = getWhaleName(conn)
		if err != nil {
			return err
		}
	}

	sql, args, err := squirrel.Insert("whale").
		Columns("address", "name", "added", "last_sent", "last_received").
		Values(address, name, time.Now(), time.Now(), time.Now()).
		ToSql()

	if err != nil {
		return err
	}

	_, err = conn.Exec(sql, args...)
	return err
}

func GetWhaleByName(conn *sqlx.DB, name string) (Whale, error) {
	var whale []Whale
	sql, args, err := squirrel.Select("*").
		From("whale").
		Where(squirrel.Eq{"name": name}).
		ToSql()

	if err != nil {
		return Whale{}, err
	}

	err = conn.Select(&whale, sql, args[0])
	if err != nil {
		return Whale{}, err
	}

	if len(whale) == 0 {
		return Whale{}, err
	}

	return whale[0], nil
}

func GetWhaleByAddress(conn *sqlx.DB, address string) (Whale, error) {
	var whale []Whale
	sql, args, err := squirrel.Select("*").
		From("whale").
		Where(squirrel.Eq{"address": address}).
		ToSql()

	if err != nil {
		return Whale{}, err
	}

	err = conn.Select(&whale, sql, args[0])
	if err != nil {
		return Whale{}, err
	}

	if len(whale) == 0 {
		return Whale{}, err
	}

	return whale[0], nil
}

func UpdateWhaleLastSent(conn *sqlx.DB, address string, time time.Time) error {
	sql, args, err := squirrel.Update("whale").
		Set("last_sent", time).
		Where(squirrel.Eq{"address": address}).
		ToSql()

	_, err = conn.Exec(sql, args...)
	return err
}

func UpdateWhaleLastReceived(conn *sqlx.DB, address string, time time.Time) error {
	sql, args, err := squirrel.Update("whale").
		Set("last_received", time).
		Where(squirrel.Eq{"address": address}).
		ToSql()

	_, err = conn.Exec(sql, args...)
	return err
}

func ExistsConfig(conn *sqlx.DB) (bool, error) {
	sql, args, err := squirrel.Select("1").
		Prefix("SELECT EXISTS (").
		From("config").
		Where(squirrel.Eq{"id": 1}).
		Suffix(")").
		ToSql()

	if err != nil {
		return false, err
	}

	var exists bool
	err = conn.QueryRow(sql, args...).Scan(&exists)

	if err != nil {
		return false, err
	}
	return exists, nil
}

func CreateConfig(conn *sqlx.DB) error {
	sql, args, err := squirrel.Insert("config").
		Columns("last_updated").
		Values(time.Now()).
		ToSql()

	if err != nil {
		return err
	}

	_, err = conn.Exec(sql, args...)
	return err
}

func GetConfig(conn *sqlx.DB) (Config, error) {
	var config []Config
	sql, args, err := squirrel.Select("*").
		From("config").
		Where(squirrel.Eq{"id": 1}).
		ToSql()

	if err != nil {
		return Config{}, err
	}

	err = conn.Select(&config, sql, args[0])
	if err != nil {
		return Config{}, err
	}

	if len(config) == 0 {
		return Config{}, err
	}

	return config[0], nil
}

func UpdateConfigLastUpdated(conn *sqlx.DB, time time.Time) error {
	sql, args, err := squirrel.Update("config").
		Set("last_updated", time).
		Where(squirrel.Eq{"id": 1}).
		ToSql()

	_, err = conn.Exec(sql, args...)
	return err
}
