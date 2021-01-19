package main

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/peter9207/black/datapoint"
	"github.com/satori/go.uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"time"
)

var groupCmd = &cobra.Command{
	Use:   "group <type>",
	Short: "grouping data by various metrics",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			err := cmd.Help()
			if err != nil {
				panic(err)
			}
			return
		}

		t := args[0]

		switch t {

		case "max":
			agg := datapoint.DatapointAggregator{DBURL: viper.GetString("database_url")}
			e := agg.GroupDatapointsByMonth()
			if e != nil {
				panic(e)
			}

		default:
			fmt.Println("unknown command ", t)
		}

	},
}

type Properties struct {
	ID        string
	Type      string
	Value     float64
	Name      string
	Meta      string
	CreatedAt time.Time
}

func saveDatapoint(data datapoint.Data, name string, db *pg.DB) (err error) {

	prop := &Properties{
		ID:        uuid.NewV4().String(),
		Name:      data.Name,
		Value:     data.Value,
		Type:      "MonthlyMax",
		CreatedAt: time.Now(),
		Meta:      name,
	}

	_, err = db.Model(prop).Insert()
	return
}
