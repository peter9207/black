package stock

import "github.com/go-pg/pg/v10"

type Stock struct {
	ID     string  `json:"id"`
	Code   string  `json:"code"`
	Date   string  `json:"date"`
	Open   float64 `json:"open"`
	High   float64 `json: "high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume int64   `json:"volume"`
}

func (s *Stock) Save(db *pg.DB) error {
	_, err := db.Model(s).Insert()
	return err
}
