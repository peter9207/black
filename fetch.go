package main

import (
	"fmt"
	// "github.com/peter9207/black/datapoint"
	"github.com/peter9207/black/fetchers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"time"
)

func fetch(fetcher fetchers.Fetcher, ticker string) (err error) {
	err = fetcher.Fetch(ticker)
	return
}

var fetchAllCmd = &cobra.Command{
	Use:   "fetchAll",
	Short: "fetch data from AA for all stored tickers",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := ConnectDB(viper.GetString("database_url"))
		if err != nil {
			panic(err)
		}
		fetcher := &fetchers.AlphaAdvantage{
			ApiKey: viper.GetString("aa_apikey"),
			DB:     db,
		}

		for _, v := range sp500 {
			if err := fetch(fetcher, v); err != nil {
				fmt.Println(err.Error())
			}
			time.Sleep(20 * time.Second)
		}

	},
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
