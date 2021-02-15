package datapoint

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/satori/go.uuid"
	"log"
	"strings"
	"sync"
)

type DataService struct {
	conn *pgxpool.Pool
}

func NewService(c *pgxpool.Pool) (ds DataService) {
	ds.conn = c
	return
}

func (ds DataService) GroupBySimilarHighWithinWindow(days int) (err error) {

	var count int64

	groupByHighQuery := `select id from datapoints order by id`

	rows, err := ds.conn.Query(context.Background(), groupByHighQuery)
	if err != nil {
		return
	}
	defer rows.Close()

	idCh := make(chan string, 30)
	errCh := make(chan error, 10)
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		go func(w *sync.WaitGroup) {
			log.Println("Starting worker", i)
			w.Add(1)
			saveRelation2(ds.conn, idCh, errCh)
			w.Done()
		}(&wg)

	}

	for rows.Next() {
		count++
		if (count % 1000) == 0 {
			log.Println("current count ", count)
		}

		var id string
		err = rows.Scan(&id)
		if err != nil {
			return
		}

		idCh <- id

	}

	close(idCh)
	close(errCh)
	wg.Wait()

	return
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
 WHERE d.id = $1
GROUP BY d.id

`

	for startID := range idCh {
		var a string
		var bs string

		err := conn.QueryRow(context.Background(), groupQuery, startID).Scan(&a, &bs)
		if err != nil {
			if err == pgx.ErrNoRows {
				log.Println("fetch group query returned nothing, skipping")
				continue
			}
			log.Println("Error fetching groups", err.Error())
			return
		}

		err = doSomething(conn, a, bs)
		if err != nil {
			log.Println("Error in worker", err.Error())
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
		if err != nil {
			if err == pgx.ErrNoRows {
				_, err = conn.Exec(context.Background(), insertRelationshipQuery, uuid.NewV4().String(), a, b, 1, dataType)
				if err != nil {
					log.Println("Error creating new relationship record", err.Error())
					return
				}
				continue
			}

			log.Println("Error fetching existing row", err)
			return
		}

		_, err := conn.Exec(context.Background(), updateQuery, score+1, id)
		if err != nil {
			log.Println("Error update new relationship record", err.Error())
			return err
		}
	}
	return
}
