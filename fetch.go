package main

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/peter9207/black/stock"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch <ticker>",
	Short: "fetch data from AA for a ticker",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			err := cmd.Help()
			if err != nil {
				panic(err)
			}
			return
		}

		ticker := args[0]
		dbURL := viper.GetString("database_url")
		apiKey := viper.GetString("aa_apikey")

		conn, err := pgxpool.Connect(context.Background(), dbURL)
		if err != nil {
			panic(err)
		}
		stockService := stock.NewService(conn, apiKey)
		err = stockService.FetchData(ticker)
		if err != nil {
			panic(err)
		}

	},
}
