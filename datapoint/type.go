package datapoint

import "github.com/nleeper/goment"
import "fmt"

type DataPoint struct {
	Type  string  `json:"type"`
	Value float64 `json:"value"`
}

type Data struct {
	Name  string
	Value float64
}

func GroupDatapointsByMonth(input []Data) (result []Data) {

	bucket := make(map[string][]Data)

	for _, v := range input {
		ment, err := goment.New(v.Name, "YYYY-MM-DD")
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		month := ment.Format("YYYY-MM")
		bucket[month] = append(bucket[month], v)
	}

	for _, v := range bucket {
		max := getMax(v)
		result = append(result, max)
	}
	return
}

func getMax(data []Data) (max Data) {

	maxIndex := 0
	var maxVal float64

	for i, v := range data {
		if maxVal < v.Value {
			maxVal = v.Value
			maxIndex = i
		}
	}

	max = data[maxIndex]
	return
}
