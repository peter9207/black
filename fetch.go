package main

import (
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
		conn, err := pgxpool.Connect(context.Background(), viper.GetString("database_url"))
		if err != nil {
			panic(err)
		}
		stockService := stock.NewService(conn, viper.GetString("aa_apikey"))
		err = stockService.FetchData(ticker)
		if err != nil {
			panic(err)
		}

	},
}
