package average

import (
	"fmt"
	"github.com/peter9207/black/shapes"
)

func Box(data []float64, size int) (filtered []float64) {

	linkedList := shapes.NewLinkedList()
	bucketSum := float64(0)

	count := float64(0)

	//initialize with first couple elements

	// fmt.Println("adding stuff to initilize")
	// for i := 0; i < int(size/2)+1; i++ {
	// 	fmt.Println(data[i])
	// 	v := data[i]
	// 	linkedList.AddFirst(v)
	// 	bucketSum = bucketSum + data[i]
	// 	count = count + 1
	// }

	for i, v := range data {

		fmt.Println("on: ", v)

		if int(i+size) < len(data) {
			// fmt.Println("adding stuff to LL", v)
			bucketSum = bucketSum + v
			linkedList.AddLast(v)
			count = count + 1
		}

		// fmt.Printf("size: %v current count: %v\n", size, count)
		if count > float64(size) {

			old := linkedList.RemoveFirst()
			// fmt.Println("removing things from list", old)
			bucketSum = bucketSum - old
			count = count - 1
		}

		// fmt.Println("bucket", linkedList)

		average := float64(bucketSum) / float64(count)
		// fmt.Println("adding value", average)
		filtered = append(filtered, average)
	}

	fmt.Println("result", filtered)
	return
}
