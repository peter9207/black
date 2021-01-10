package average_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"github.com/peter9207/black/average"
)

var _ = Describe("LinearFilter", func() {

	Describe("when window is 1", func() {
		data := []float64{1, 2, 1, 2, 1, 2}

		result := average.Box(data, 1)

		It("should have a length equal to data", func() {
			Ω(len(result)).Should(Equal(len(data)))
		})

		It("should return the same values as original", func() {

			fmt.Println("result", result)
			for i := range data {
				Ω(result[i]).Should(Equal(data[i]))

			}

		})

	})

})
