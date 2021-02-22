package stock

import (
	"encoding/json"
	"github.com/go-pg/pg/v10"
	"github.com/segmentio/kafka-go"
	"time"
)

type StockData struct {
	ID     string  `json:"id"`
	Code   string  `json:"code"`
	Date   string  `json:"date"`
	Open   float64 `json:"open"`
	High   float64 `json: "high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume int64   `json:"volume"`
}

type Stock struct {
	Ticker          string    `json:"ticker"`
	LastUpdatedDate time.Time `json:"lastUpdatedAt"`
	Description     string    `json:"description"`
}

func (s *StockData) Save(db *pg.DB) error {
	_, err := db.Model(s).Insert()
	return err
}

func Consume(conn *kafka.Conn, db *pg.DB) (err error) {
	batch := conn.ReadBatch(10e2, 1e6) // fetch 10KB min, 1MB max

	b := make([]byte, 10e3) // 10KB max per message
	for {
		_, err := batch.Read(b)
		if err != nil {
			break
		}
		stock := StockData{}
		err = json.Unmarshal(b, &stock)
		if err != nil {
			return err
		}

		err = stock.Save(db)
		if err != nil {
			return err
		}
	}
	return
}
