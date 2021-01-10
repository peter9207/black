package main

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/peter9207/black/datapoint"
	"github.com/satori/go.uuid"
	"github.com/spf13/cobra"
	"time"
)

var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "grouping data by various metrics",
}

var groupMaxCmd = &cobra.Command{
	Use:   "max <filename>",
	Short: "group data by months and find max",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			err := cmd.Help()
			if err != nil {
				panic(err)
			}
			return
		}

		csvFile := args[0]
		stocks, err := readExport(csvFile)
		if err != nil {
			panic(err)
			return
		}

		data := []datapoint.Data{}
		for _, v := range stocks {

			fmt.Println("stuffs")
			d := datapoint.Data{
				Name:  v.Date,
				Value: v.High,
			}
			data = append(data, d)
		}

		db, err := ConnectDB("postgres://postgres:password@localhost:5432/postgres?sslmode=disable")
		if err != nil {
			panic(err)
		}

		results := datapoint.GroupDatapointsByMonth(data)

		for _, v := range results {
			fmt.Println("saving")
			if err := saveDatapoint(v, db); err != nil {
				fmt.Println(err.Error())
			}
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

func saveDatapoint(data datapoint.Data, db *pg.DB) (err error) {

	prop := &Properties{
		ID:        uuid.NewV4().String(),
		Name:      data.Name,
		Value:     data.Value,
		Type:      "MonthlyMax",
		CreatedAt: time.Now(),
	}

	_, err = db.Model(prop).Insert()
	return
}
