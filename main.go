package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/peter9207/black/datapoint"
	"github.com/peter9207/black/fetchers"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"strconv"
	"time"
)

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
	Use:   "printEnv",
	Short: "print a list of current config parameters",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		printEnv()
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
	Use:   "black",
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

var downloadCmd = &cobra.Command{
	Use:   "download <location>",
	Short: "download the data to a location",
	Long:  "location currently could be kafka",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			err := cmd.Help()
			if err != nil {
				panic(err)
			}
			return
		}

		location := args[0]

		db, err := ConnectDB(viper.GetString("database_url"))
		if err != nil {
			panic(err)
		}
		fetcher := &fetchers.AlphaAdvantage{
			ApiKey: viper.GetString("aa_apikey"),
			DB:     db,
		}

		switch location {
		case "kafka":
			conn, err := kafka.DialLeader(context.Background(), "tcp", viper.GetString("kafka_broker"), "stock-events", 0)
			if err != nil {
				log.Fatal("failed to dial leader:", err)
			}
			for _, v := range sp500 {
				err = fetcher.ToKafka(v, conn)
				if err != nil {
					panic(err)
				}
				time.Sleep(20 * time.Second)

			}

		default:
			fmt.Println("invalid type")
		}

	},
}

func printEnv() {

	knownEnvs := []string{
		"database_url", "AA_APIKEY", "kafka_broker",
	}

	for _, v := range knownEnvs {
		fmt.Printf("%s : %s\n", v, viper.GetString(v))
	}

}

var aggCmd = &cobra.Command{
	Use:   "agg",
	Short: "process various types of data",
}
var windowMaxCmd = &cobra.Command{
	Use:   "window <days>",
	Short: "group maxes by their related max points",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Help()
			return
		}

		i, err := strconv.Atoi(args[0])
		if err != nil {
			panic(err)
		}

		conn, err := pgxpool.Connect(context.Background(), viper.GetString("database_url"))
		if err != nil {
			panic(err)
		}
		service := datapoint.NewService(conn)

		err = service.GroupBySimilarHighWithinWindow(i)
		if err != nil {
			panic(err)
		}

	},
}

var monthlyMaxCmd = &cobra.Command{
	Use:   "monthlyMax",
	Short: "process various types of data",
	Run: func(cmd *cobra.Command, args []string) {

		agg := datapoint.DatapointAggregator{
			DBURL: viper.GetString("database_url"),
		}

		err := agg.GroupMonthlyMax()
		if err != nil {
			panic(err)
		}

	},
}

var relationshipsCmd = &cobra.Command{
	Use:   "relationships",
	Short: "process various types of data",
	Run: func(cmd *cobra.Command, args []string) {

		err := datapoint.CountRelationships(viper.GetString("database_url"))
		if err != nil {
			panic(err)
		}

	},
}

func main() {

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	viper.SetDefault("DATABASE_URL", "postgres://postgres:password@localhost:5432/postgres?sslmode=disable")
	viper.SetDefault("AA_APIKEY", "no_default")
	viper.SetDefault("kafka_broker", "localhost:9092")

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

	rootCmd.AddCommand(aggCmd)

	aggCmd.AddCommand(windowMaxCmd)
	aggCmd.AddCommand(relationshipsCmd)
	aggCmd.AddCommand(monthlyMaxCmd)

	rootCmd.AddCommand(loadCmd)

	rootCmd.AddCommand(fetchCmd)
	rootCmd.AddCommand(initConfig)
	rootCmd.AddCommand(fetchAllCmd)
	rootCmd.AddCommand(downloadCmd)

	rootCmd.Execute()
}
