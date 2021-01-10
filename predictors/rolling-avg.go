package predictors

import (
	"github.com/peter9207/black/average"
)

type RollingAvgPredictor struct {
	Score      int64
	windowSize int64
}

func RollingAvg(windowSize int64) (r *RollingAvgPredictor) {
	r = &RollingAvgPredictor{windowSize: windowSize}
	return
}

func (r *RollingAvgPredictor) Predict(data []float64) (b float64) {

	var result []float64
	var sum float64

	averages := average.Rolling(data, r.windowSize)

	for i := range data {
		d := averages[i] - data[i]
		result = append(result, d)
	}

	for _, v := range result {
		sum = sum + v
	}
	b = sum / float64(len(data))
	return
}
