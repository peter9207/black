package datapoint

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/satori/go.uuid"
	"sort"
	"strings"
	"sync"
)

type Relationship struct {
	ID    string `json:"id"`
	CodeA string `json:"code_a"`
	CodeB string `json:""code_b`
	Score int64  `json:"score"`
}

var groupQuery = `
select string_agg(code, ' ') from datapoints
group by type, key
order by key
`

func CountRelationships(url string) (err error) {
	conn, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		return
	}

	rows, err := conn.Query(context.Background(), groupQuery)
	if err != nil {
		return
	}

	for rows.Next() {

		var tickers string

		err = rows.Scan(&tickers)
		if err != nil {
			return
		}

		err = saveRelationsAsync(tickers, conn)
		if err != nil {
			return
		}

	}
	return
}

var selectQuery = "SELECT id, score from relationships where code_a = $1 and code_b = $2"
var insertRelationshipQuery = "insert into relationships (id, code_a, code_b, score) VALUES ($1, $2, $3, $4)"
var updateQuery = "update relationships SET score = $1 where id =$2"

func saveRelation(tickers []string, conn *pgxpool.Pool) (err error) {
	rows, err := conn.Query(context.Background(), selectQuery, tickers[0], tickers[1])
	if err != nil {
		return err
	}
	defer rows.Close()

	var id string
	var score int64
	var found bool

	for rows.Next() {
		err = rows.Scan(&id, &score)
		if err != nil {
			return err
		}
		found = true
	}

	if found {
		_, err = conn.Exec(context.Background(), updateQuery, score+1, id)
		if err != nil {
			return err
		}
	} else {
		_, err = conn.Exec(context.Background(), insertRelationshipQuery, uuid.NewV4().String(), tickers[0], tickers[1], 1)
		if err != nil {
			return err
		}
	}
	return
}

func saveRelationAsync(tickersCh chan []string, conn *pgxpool.Pool, errCh chan error) {
	for tickers := range tickersCh {
		err := saveRelation(tickers, conn)
		if err != nil {
			errCh <- err
		}

	}
	return
}

func saveRelationsAsync(tickerString string, conn *pgxpool.Pool) (err error) {
	tickers := strings.Split(tickerString, " ")

	tickersCh := make(chan []string, 20)
	errCh := make(chan error)

	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		fmt.Println("starting worker", i)
		go func(waiting sync.WaitGroup) {
			saveRelationAsync(tickersCh, conn, errCh)
			wg.Done()
		}(wg)

	}

	for _, v1 := range tickers {
		for _, v2 := range tickers {

			if v1 == v2 {
				continue
			}

			a := []string{v1, v2}
			sort.Strings(a)

			tickersCh <- a

		}
	}
	close(tickersCh)

	wg.Wait()
	return
}

func saveRelations(tickerString string, conn *pgxpool.Pool) (err error) {
	tickers := strings.Split(tickerString, " ")
	for _, v1 := range tickers {
		for _, v2 := range tickers {

			if v1 == v2 {
				continue
			}

			a := []string{v1, v2}
			sort.Strings(a)

			err = saveRelation(a, conn)

		}
	}
	return
}
