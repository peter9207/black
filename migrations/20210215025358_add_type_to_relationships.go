package main

import (
	"github.com/go-pg/pg/v10/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v3"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec("alter table relationships add column type varchar")
		if err != nil {
			return err
		}
		_, err = db.Exec("update relationships set type = 'max-same-day'")
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec("alter table relationships drop column type")
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20210215025358_add_type_to_relationships", up, down, opts)
}
