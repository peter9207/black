package main

import (
	"github.com/go-pg/pg/v10/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v3"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS properties (
  id UUID PRIMARY KEY,
  type varchar,
  value numeric,
  name varchar,
  meta varchar,
  created_at timestamp
  )
`)
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec(" drop table properties ")
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20210108064627_create_properties", up, down, opts)
}
