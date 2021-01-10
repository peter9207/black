package predictors

import (
	"github.com/peter9207/black/average"
)

type Predictor interface {
	Predict(data []float64) float64
	POI(data []float64) []int
}

type SimplePredictor struct {
	Score      int64
	windowSize int64
}

func SimpleRolling(bucketSize int64) (p *SimplePredictor) {
	return &SimplePredictor{windowSize: bucketSize}
}

func (p *SimplePredictor) Predict(data []float64) (b float64) {

	averages := average.Rolling(data, p.windowSize)
	b = averages[len(averages)-1] - averages[0]
	return
}

type RollingWindowCrossingPOI struct {
	Index      int
	Value      float64
	Increasing bool
}

// a POI for this version is currently loosely defined as when rolling overage crosses
// data and vice versa
func (p *SimplePredictor) POI(data []float64) (output []RollingWindowCrossingPOI) {

	var isAbove bool
	averages := average.Rolling(data, p.windowSize)

	for i, v := range data {
		var newIsAbove bool
		if data[i] > averages[i] {
			newIsAbove = true
		} else {
			newIsAbove = false
		}

		if isAbove != newIsAbove {
			poi := RollingWindowCrossingPOI{
				Index:      i,
				Value:      v,
				Increasing: newIsAbove,
			}
			output = append(output, poi)
		}
		isAbove = newIsAbove
	}

	return
}
