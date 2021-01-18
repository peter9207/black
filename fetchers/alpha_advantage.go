package fetchers

import (
	"encoding/json"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/peter9207/black/stock"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"strconv"
)

var requestURLTemplate = "https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=%s&outputsize=full&apikey=%s"

type AlphaAdvantage struct {
	ApiKey string
	DB     *pg.DB
}

// type Stock struct {
// 	ID     string
// 	Code   string
// 	Date   string
// 	Open   float64
// 	High   float64
// 	Low    float64
// 	Close  float64
// 	Volume int64
// }

type AAResponse struct {
	Metadata map[string]string `json:"Meta Data"`
	Data     map[string]AAData `json:"Time Series (Daily)"`
}

type AAData struct {
	Open   float64
	Close  float64
	Low    float64
	High   float64
	Volume int64
}

func getData(source map[string]string, key string) (f float64, err error) {
	openString, ok := source[key]
	if !ok {
		err = fmt.Errorf("value %s not found in response")
		return
	}
	f, err = strconv.ParseFloat(openString, 64)
	return
}

func (a *AAData) UnmarshalJSON(data []byte) error {

	var v map[string]string

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	open, err := getData(v, "1. open")
	if err != nil {
		return err
	}
	a.Open = open

	high, err := getData(v, "2. high")
	if err != nil {
		return err
	}
	a.High = high

	low, err := getData(v, "3. low")
	if err != nil {
		return err
	}
	a.Low = low

	close, err := getData(v, "4. close")
	if err != nil {
		return err
	}
	a.Close = close

	volumeString, ok := v["5. volume"]
	if !ok {
		return fmt.Errorf("value %s not found in response", "5. volume")
	}
	volume, err := strconv.ParseInt(volumeString, 10, 64)
	if err != nil {
		return err
	}
	a.Volume = volume
	return nil
}

func (aa *AlphaAdvantage) Fetch(ticker string) (err error) {
	var url = fmt.Sprintf(requestURLTemplate, ticker, aa.ApiKey)
	fmt.Println("fetching data from ", url)

	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	fmt.Println("got back bytes ", len(body))

	response := AAResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return
	}

	stockName := response.Metadata["2. Symbol"]
	fmt.Println("name", stockName)

	for k, v := range response.Data {

		s := &stock.Stock{
			ID:     uuid.NewV4().String(),
			Date:   k,
			Code:   ticker,
			Open:   v.Open,
			Close:  v.Close,
			High:   v.High,
			Low:    v.Low,
			Volume: v.Volume,
		}
		// _, err = aa.DB.Model(s).Insert()

		err = s.Save(aa.DB)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	return
}
