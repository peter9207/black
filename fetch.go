package main

import (
	"fmt"
	"github.com/peter9207/black/datapoint"
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

		fetcher := AlphaAdvantage{
			ApiKey: viper.GetString("aa_apikey"),
		}

		stockData, err := fetcher.fetch(ticker)
		if err != nil {
			panic(err)
		}

		data := []datapoint.Data{}
		for _, v := range stockData.Data {

			d := datapoint.Data{
				Name:  v.Date,
				Value: v.High,
			}
			data = append(data, d)
		}

		db, err := ConnectDB(viper.GetString("DATABASE_URL"))
		if err != nil {
			panic(err)
		}

		results := datapoint.GroupDatapointsByMonth(data)

		fmt.Println("saving ", stockData.Name)
		fmt.Println("saving ", len(stockData.Data))

		for _, v := range results {
			if err := saveDatapoint(v, stockData.Name, db); err != nil {
				fmt.Println(err.Error())
			}
		}
	},
}
