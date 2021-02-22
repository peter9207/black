package main

// import (
// 	"context"
// 	"fmt"
// 	"github.com/peter9207/black/stock"
// 	"github.com/segmentio/kafka-go"
// 	"github.com/spf13/cobra"
// 	"github.com/spf13/viper"
// 	"log"
// )

// var loadCmd = &cobra.Command{
// 	Use:   "load <type>",
// 	Short: "load data from sources",
// 	Run: func(cmd *cobra.Command, args []string) {

// 		if len(args) < 1 {
// 			err := cmd.Help()
// 			if err != nil {
// 				panic(err)
// 			}
// 			return
// 		}

// 		source := args[0]
// 		switch source {
// 		case "kafka":
// 			err := loadFromKafka()
// 			if err != nil {
// 				panic(err)
// 			}
// 		default:
// 			fmt.Println("invalid type")
// 		}

// 	},
// }

// func loadFromKafka() error {

// 	db, err := ConnectDB(viper.GetString("database_url"))
// 	if err != nil {
// 		panic(err)
// 	}
// 	conn, err := kafka.DialLeader(context.Background(), "tcp", viper.GetString("kafka_broker"), "stock-events", 0)
// 	if err != nil {
// 		log.Fatal("failed to dial leader:", err)
// 	}

// 	err = stock.Consume(conn, db)
// 	return err
// }
