package main

import (
	"github.com/go-pg/pg/v10/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v3"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec(`
CREATE TABLE datapoints (
id UUID primary key,
type varchar,
value numeric,
code varchar,
key varchar
)
`)
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec(" drop table datapoints")
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20210118081455_create_datapoint_table", up, down, opts)
}
