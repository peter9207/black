package average

import (
	"github.com/peter9207/black/shapes"
)

func Simple(input []float64) (mean float64) {

	sum := Sum(input)
	mean = float64(sum) / float64(len(input))
	return
}

func Sum(input []float64) (sum float64) {

	for _, v := range input {
		sum += v
	}
	return

}

func Rolling(input []float64, size int64) (averages []float64) {

	linkedList := shapes.NewLinkedList()
	bucketSum := float64(0)

	count := float64(0)

	for _, v := range input {
		bucketSum = bucketSum + v

		linkedList.AddFirst(v)

		count = count + 1
		if count > float64(size) {
			old := linkedList.RemoveFirst()
			bucketSum = bucketSum - old
			count = count - 1
		}

		average := float64(bucketSum) / float64(count)
		averages = append(averages, average)

	}

	return
}
