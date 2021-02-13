package main

import (
	"log"
	"os"

	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10"
)

const directory = "."

func main() {
	db := pg.Connect(&pg.Options{
		Addr:     "localhost:4432",
		User:     "postgres",
		Database: "postgres",
		Password: "password",
	})

	err := migrations.Run(db, directory, os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}
