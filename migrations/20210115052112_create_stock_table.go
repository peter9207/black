package main

import (
	"github.com/go-pg/pg/v10/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v3"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec(`
CREATE TABLE stocks (
	id uuid PRIMARY KEY,
    code varchar,
	date date,
	OPEN numeric,
	CLOSE numeric,
	high numeric,
	low numeric,
	volume numeric
)
`)
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec(" drop table stocks")
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20210115052112_create_stock_table", up, down, opts)
}
