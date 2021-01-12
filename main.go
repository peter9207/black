package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/peter9207/black/predictors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

type Stock struct {
	Date     string
	Open     float64
	High     float64
	Low      float64
	Close    float64
	AdjClose float64
	Volume   int64
	Name     string
}

func readExport(stockName, filename string) (stocks []Stock, err error) {

	f, _ := os.Open(filename)
	reader := csv.NewReader(bufio.NewReader(f))
	first := true
	for {
		var line []string
		line, err = reader.Read()
		if err == io.EOF {
			return stocks, nil
		} else if err != nil {
			log.Fatal(err)
			return
		}
		if first {
			first = false
			continue
		}

		stocks = append(stocks, Stock{
			Date:     line[0],
			Open:     parseFloat(line[1]),
			High:     parseFloat(line[2]),
			Low:      parseFloat(line[3]),
			Close:    parseFloat(line[4]),
			AdjClose: parseFloat(line[5]),
			Volume:   parseInt(line[6]),
			Name:     stockName,
		})
	}

}

func parseFloat(s string) (f float64) {
	var err error
	if f, err = strconv.ParseFloat(s, 64); err != nil {
		panic(err)
	}
	return
}

func parseInt(s string) (i int64) {
	var err error
	if i, err = strconv.ParseInt(s, 10, 64); err != nil {
		panic(err)
	}
	return
}

var simpleCmd = &cobra.Command{
	Use:   "simple <filename>",
	Short: "a simple rolling average calculation",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			err := cmd.Help()
			if err != nil {
				panic(err)
			}
			return
		}

		csvFile := args[1]
		name := args[0]
		stocks, err := readExport(name, csvFile)

		if err != nil {
			log.Fatal(err)
			return
		}

		data := []float64{}
		for _, v := range stocks {
			data = append(data, v.Close)
		}

		producer, err := NewEventProducer()
		if err != nil {
			panic(err)
		}

		p := predictors.SimpleRolling(10)
		result := p.POI(data)

		format := "2006-01-02"

		for _, v := range result {

			date := stocks[v.Index].Date

			t, err := time.Parse(format, date)
			if err != nil {
				fmt.Println("error parsing date", date)
				continue
			}

			event := RollingWindowCrossingEvent{
				Value:      v.Value,
				Window:     10,
				Meta:       name,
				Increasing: v.Increasing,
				Date:       t.Unix(),
			}

			err = producer.produceEvent(event)
			if err != nil {
				fmt.Println("error producing", err)
			}

		}

		log.Printf("result: %v", result)
	},
}

var esCmd = &cobra.Command{
	Use:   "es ",
	Short: "test es is connected",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		client, err := connectES()
		if err != nil {
			panic(err)
		}
		res, err := client.Info()
		if err != nil {
			return
		}

		defer res.Body.Close()
		log.Println(res)

	},
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "test various dependencies",
}

var rootCmd = &cobra.Command{
	Use:   "S",
	Short: "an attempt to try various means of analysing stocks data",
}

var initConfig = &cobra.Command{
	Use:   "init",
	Short: "init a config file",
	Run: func(cmd *cobra.Command, args []string) {
		err := viper.SafeWriteConfig()
		if err != nil {
			panic(err)
		}
	},
}

func main() {

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	viper.SetDefault("DATABASE_URL", "postgres://postgres:password@localhost:5432/postgres?sslmode=disable")
	viper.SetDefault("AA_APIKEY", "no_default")

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Config file not found...")
		} else {
			panic(err)
		}
	}

	testCmd.AddCommand(esCmd)
	rootCmd.AddCommand(simpleCmd)
	rootCmd.AddCommand(testCmd)

	rootCmd.AddCommand(groupCmd)
	groupCmd.AddCommand(groupMaxCmd)

	rootCmd.AddCommand(fetchCmd)
	rootCmd.AddCommand(initConfig)

	rootCmd.Execute()
}
