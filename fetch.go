package main

import (
	"github.com/peter9207/black/fetchers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func fetch(fetcher fetchers.Fetcher, ticker string) (err error) {
	err = fetcher.Fetch(ticker)
	return
}

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
		db, err := ConnectDB(viper.GetString("database_url"))
		if err != nil {
			panic(err)
		}

		fetcher := &fetchers.AlphaAdvantage{
			ApiKey: viper.GetString("aa_apikey"),
			DB:     db,
		}

		err = fetch(fetcher, ticker)
		if err != nil {
			panic(err)
		}

	},
}
