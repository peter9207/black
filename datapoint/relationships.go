package datapoint

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/satori/go.uuid"
	"strings"
)

type DataService struct {
	conn *pgxpool.Pool
}

func NewService(c *pgxpool.Pool) (ds DataService) {
	ds.conn = c
	return
}

func getDatapoints(conn *pgxpool.Pool) {

}

func (ds DataService) GroupBySimilarHighWithinWindow(days int) (err error) {

	// select string_agg(code, ' '), key from datapoints
	// group by type, key
	// order by key
	groupByHighQuery := `select * from datapoints order by id`

	rows, err := ds.conn.Query(context.Background(), groupByHighQuery)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {

	}

	return
}

type SaveRelationRequest struct {
	a       string
	tickers string
}

func saveRelation2(conn *pgxpool.Pool, idCh chan string, errCh chan error) {

	groupQuery = `
 SELECT
	d.code,
	string_agg(d2.code, ' ')
FROM
	datapoints d
	JOIN datapoints d2 ON d.code != d2.code
		AND d.type = d2.type
		AND d2.key::TIMESTAMP >= d.key::TIMESTAMP
		AND d2.key::TIMESTAMP < d.key::TIMESTAMP + interval '5' day
 WHERE d.id = ?
GROUP BY
	d.id;

`

	for startId := range idCh {

		var a string
		var bs string

		err := conn.QueryRow(context.Background(), groupQuery, startId).Scan(&a, &bs)
		if err != nil {
			return

		}

		err = doSomething(conn, a, bs)
		if err != nil {
			return
		}

	}

	return
}

func doSomething(conn *pgxpool.Pool, a, bs string) (err error) {

	tickers := strings.Split(bs, " ")

	selectQuery := "SELECT id, score from relationships where code_a = $1 and code_b = $2 and type = $3"
	var insertRelationshipQuery = "insert into relationships (id, code_a, code_b, score, type) VALUES ($1, $2, $3, $4, $5)"
	var updateQuery = "update relationships SET score = $1 where id =$2"

	dataType := "window-5-days"

	for _, b := range tickers {

		var id string
		var score int64

		err = conn.QueryRow(context.Background(), selectQuery, a, b, dataType).Scan(&id, &score)
		if err == pgx.ErrNoRows {
			_, err = conn.Exec(context.Background(), insertRelationshipQuery, uuid.NewV4().String(), a, b, dataType)
			if err != nil {
				return
			}
		}
		if err != nil {
			return
		}

		_, err := conn.Exec(context.Background(), updateQuery, score+1, id)
		if err != nil {
			return err
		}
	}
	return
}
