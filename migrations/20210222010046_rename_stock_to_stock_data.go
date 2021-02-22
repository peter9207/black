package main

import (
	"github.com/go-pg/pg/v10/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v3"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec("alter table stocks rename to stock_data ")
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec("alter table stocks_data rename to stocks")
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20210222010046_rename_stock_to_stock_data", up, down, opts)
}
