package main

import (
	"github.com/go-pg/pg/v10"
	// migrations "github.com/robinjoseph08/go-pg-migrations/v3"
)

func ConnectDB(url string) (db *pg.DB, err error) {

	opt, err := pg.ParseURL(url)
	if err != nil {
		return
	}
	db = pg.Connect(opt)
	return
}
