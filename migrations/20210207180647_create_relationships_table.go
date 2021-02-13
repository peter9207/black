package main

import (
	"github.com/go-pg/pg/v10/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v3"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec(`
CREATE TABLE relationships (
	id uuid PRIMARY KEY,
    "code_a" varchar,
    "code_b" varchar,
    score numeric
)
`)
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec("drop table relationships")
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20210207180647_create_relationships_table", up, down, opts)
}
