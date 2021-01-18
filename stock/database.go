package stock

import "github.com/go-pg/pg/v10"

type Stock struct {
	ID     string
	Code   string
	Date   string
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int64
}

func (s *Stock) Save(db *pg.DB) error {
	_, err := db.Model(s).Insert()
	return err
}
