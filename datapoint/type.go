package datapoint

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/satori/go.uuid"
)

type DataPoint struct {
	Type  string  `json:"type"`
	Value float64 `json:"value"`
	Code  string
	Key   string
}

type Data struct {
	Name  string
	Value float64
}

var groupQuery = `
SELECT
    code,
	extract(month FROM date)::text AS month,
	extract(year FROM date)::text AS year,
	max(high) AS max
FROM
	stocks
GROUP BY
	code,
	year,
	month;
`

var insertQuery = `
INSERT INTO datapoints (id, code, key, type, value) values ($1, $2, $3, $4, $5)
`

type DatapointAggregator struct {
	DBURL string
}

func (da DatapointAggregator) GroupDatapointsByMonth() (err error) {

	conn, err := pgxpool.Connect(context.Background(), da.DBURL)
	if err != nil {
		return
	}
	defer conn.Close()

	rows, err := conn.Query(context.Background(), groupQuery)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var code string
		var month string
		var year string
		var max float64
		err = rows.Scan(&code, &month, &year, &max)
		if err != nil {
			return
		}
		key := fmt.Sprintf("%v-%v", year, month)

		id := uuid.NewV4().String()

		_, err = conn.Exec(context.Background(), insertQuery, id, code, key, "monthly-max", max)
		if err != nil {
			return
		}

	}
	return
}
