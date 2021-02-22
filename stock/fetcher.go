package stock

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Fetcher interface {
	Fetch(ticker string) error
}

type Service struct {
	conn   *pgxpool.Pool
	apiKey string
}

func NewService(conn *pgxpool.Pool, apiKey string) (s Service) {
	s = Service{
		conn:   conn,
		apiKey: apiKey,
	}
	return
}

var getLatestQuery = `select date from stock_data where code = $1 order by date desc limit 1`
var insertStockDataQuery = `
INSERT INTO "stock_data" ("id", "code", "date", "open", "close", "high", "low", "volume")
                  VALUES ($1,     $2,     $3,    $4,     $5,     $6,      $7,     $8); `
var dateFormat = "2006-01-02"

func (s Service) FetchData(ticker string) (err error) {

	var lastDate time.Time

	err = s.conn.QueryRow(context.Background(), getLatestQuery, ticker).Scan(&lastDate)
	if err != nil {
		log.Println("error querying existing data", err)
		return
	}

	var url = fmt.Sprintf(requestURLTemplate, ticker, s.apiKey)
	fmt.Println("fetching data from ", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error querying AA", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	response := AAResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println("Error parsing json", err)
		return
	}

	stockName := response.Metadata.Symbol
	fmt.Println("name", stockName)

	for k, v := range response.Data {
		recordDate, err := time.Parse(dateFormat, k)
		if err != nil {
			log.Println("Error parsing date", err)
			return err
		}

		if recordDate.Before(lastDate) {
			continue
		}

		id := uuid.NewV4().String()
		_, err = s.conn.Exec(context.Background(),
			insertStockDataQuery,
			id,
			ticker,
			k,
			v.Open,
			v.Close,
			v.High,
			v.Low,
			v.Volume)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	return

}
