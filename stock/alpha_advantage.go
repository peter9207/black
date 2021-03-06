package stock

import (
	"encoding/json"
	"fmt"
	"github.com/go-pg/pg/v10"
	// "github.com/satori/go.uuid"
	// "github.com/segmentio/kafka-go"
	// "io/ioutil"
	// "net/http"
	"strconv"
	// "time"
)

var requestURLTemplate = "https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=%s&outputsize=full&apikey=%s"

type AlphaAdvantage struct {
	ApiKey string
	DB     *pg.DB
}

type AAMetadata struct {
	Info          string `json:"1. Information"`
	Symbol        string `json:"2. Symbol"`
	LastRefreshed string `json:"3. Last Refreshed"`
	Size          string `json:"4. Output Size"`
	Timezone      string `json:"5. Time Zone"`
}

type AAResponse struct {
	Metadata AAMetadata        `json:"Meta Data"`
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

// func (aa *AlphaAdvantage) ToKafka(ticker string, conn *kafka.Conn) (err error) {

// 	fmt.Println("starting to fetch ", ticker)

// 	stocks, err := aa.readTickerRequest(ticker)
// 	if err != nil {
// 		return
// 	}

// 	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

// 	for _, v := range stocks {

// 		jsonData, err := json.Marshal(v)
// 		if err != nil {
// 			return err
// 		}

// 		msg := kafka.Message{Value: jsonData, Key: []byte(v.ID)}

// 		_, err = conn.WriteMessages(msg)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return
// }

// func (aa *AlphaAdvantage) readTickerRequest(ticker string) (stocks []*StockData, err error) {
// 	var url = fmt.Sprintf(requestURLTemplate, ticker, aa.ApiKey)

// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return
// 	}
// 	defer resp.Body.Close()

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return
// 	}
// 	response := AAResponse{}
// 	err = json.Unmarshal(body, &response)
// 	if err != nil {
// 		return
// 	}
// 	for k, v := range response.Data {

// 		s := &StockData{
// 			ID:     uuid.NewV4().String(),
// 			Date:   k,
// 			Code:   ticker,
// 			Open:   v.Open,
// 			Close:  v.Close,
// 			High:   v.High,
// 			Low:    v.Low,
// 			Volume: v.Volume,
// 		}

// 		stocks = append(stocks, s)
// 	}
// 	return
// }

// func (aa *AlphaAdvantage) Fetch(ticker string) (err error) {

// 	var url = fmt.Sprintf(requestURLTemplate, ticker, aa.ApiKey)
// 	fmt.Println("fetching data from ", url)

// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return
// 	}
// 	defer resp.Body.Close()

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return
// 	}

// 	response := AAResponse{}
// 	err = json.Unmarshal(body, &response)
// 	if err != nil {
// 		return
// 	}

// 	stockName := response.Metadata.Symbol
// 	fmt.Println("name", stockName)

// 	for k, v := range response.Data {

// 		s := &StockData{
// 			ID:     uuid.NewV4().String(),
// 			Date:   k,
// 			Code:   ticker,
// 			Open:   v.Open,
// 			Close:  v.Close,
// 			High:   v.High,
// 			Low:    v.Low,
// 			Volume: v.Volume,
// 		}

// 		err = s.Save(aa.DB)
// 		if err != nil {
// 			fmt.Println(err.Error())
// 		}
// 	}
// 	return
// }
