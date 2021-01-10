package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/satori/go.uuid"
)

type RollingWindowCrossingEvent struct {
	Date       int64   `json:"timestamp"`
	Value      float64 `json:"value"`
	Window     int     `json:"size"`
	Meta       string  `json:"meta"`
	Increasing bool    `json:"isIncreasing"`
}

type EventProducer struct {
	esClient *elasticsearch.Client
}

func NewEventProducer() (ep *EventProducer, err error) {
	client, err := connectES()
	if err != nil {
		return
	}
	ep = &EventProducer{
		esClient: client,
	}
	return
}

func (ep *EventProducer) produceEvent(e RollingWindowCrossingEvent) (err error) {

	data, err := json.Marshal(e)
	if err != nil {
		return
	}

	id := uuid.NewV4()

	request := esapi.IndexRequest{
		Index:      "events",
		DocumentID: id.String(),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := request.Do(context.Background(), ep.esClient)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		fmt.Printf("[%s] Error indexing document ID=%s", res.Status(), id.String())
	}

	fmt.Println(res.Body)

	return
}

func connectES() (client *elasticsearch.Client, err error) {
	client, err = elasticsearch.NewDefaultClient()
	return

}
