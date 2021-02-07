package datapoint

import (
	"context"
	// "fmt"
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

// var groupQuery = `
// SELECT
//     code,
// 	extract(month FROM date)::text AS month,
// 	extract(year FROM date)::text AS year,
//     extract(day FROM date):: text AS day,
// 	max(high) AS max
// FROM
// 	stocks
// GROUP BY
// 	code,
//     date
// `

var insertQuery = `
INSERT INTO datapoints (id, code, key, type, value) values ($1, $2, $3, $4, $5)
`

type DatapointAggregator struct {
	DBURL string
}

// func (da DatapointAggregator) GroupDatapointsByMonth() (err error) {

// 	conn, err := pgxpool.Connect(context.Background(), da.DBURL)
// 	if err != nil {
// 		return
// 	}
// 	defer conn.Close()

// 	rows, err := conn.Query(context.Background(), groupQuery)
// 	if err != nil {
// 		return
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var code string
// 		var month string
// 		var day string
// 		var year string
// 		var max float64
// 		err = rows.Scan(&code, &month, &year, &day, &max)
// 		if err != nil {
// 			return
// 		}
// 		key := fmt.Sprintf("%v-%v-%s", year, month, day)

// 		id := uuid.NewV4().String()

// 		_, err = conn.Exec(context.Background(), insertQuery, id, code, key, "monthly-max-1", max)
// 		if err != nil {
// 			return
// 		}

// 	}
// 	return
// }

var dateRangeQuery = `
SELECT
	code,
	max(date)::text endDate,
	min(date)::text startDate,
	extract(month FROM date)::text AS month,
	extract(year FROM date)::text AS year
FROM
	stocks
GROUP BY
	code,
	month,
	year
ORDER BY
	code,
	startDate,
	endDate;
`

var findMaxQuery = `
SELECT
	date::text,
	high,
	code
FROM
	stocks
WHERE
	code = $1
	AND date BETWEEN $2 AND $3
ORDER BY
	high DESC
LIMIT 1;
`

func findMax(code, start, end string, conn *pgxpool.Pool) (date string, value float64, err error) {

	rows, err := conn.Query(context.Background(), findMaxQuery, code, start, end)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var s string
		err = rows.Scan(&date, &value, &s)
		if err != nil {
			return
		}
	}

	return
}

func (da DatapointAggregator) GroupMonthlyMax() (err error) {

	conn, err := pgxpool.Connect(context.Background(), da.DBURL)
	if err != nil {
		return
	}
	defer conn.Close()

	rows, err := conn.Query(context.Background(), dateRangeQuery)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {

		var code string
		var startDate string
		var endDate string
		var year string
		var month string

		err = rows.Scan(&code, &endDate, &startDate, &month, &year)
		if err != nil {
			return
		}

		d, v, err := findMax(code, startDate, endDate, conn)
		if err != nil {
			return err
		}

		id := uuid.NewV4().String()
		_, err = conn.Exec(context.Background(), insertQuery, id, code, d, "monthly-max-1", v)
		if err != nil {
			return err
		}

	}

	return
}
