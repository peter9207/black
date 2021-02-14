package main

import (
	"github.com/go-pg/pg/v10/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v3"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec("CREATE INDEX code_b_idx on relationships ( code_b )")
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec(`drop index code_b_idx`)
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20210213195100_add_idx_to_relationships_code_b", up, down, opts)
}
