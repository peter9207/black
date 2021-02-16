package main

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/peter9207/black/datapoint"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strconv"
)

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
