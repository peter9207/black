package main

import (
	"github.com/go-pg/pg/v10/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v3"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec("create table stocks (ticker varchar primary key, last_updated_date date, description varchar )")
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec("drop table stocks")
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20210222010156_create_stocks", up, down, opts)
}
