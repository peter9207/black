package average_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/peter9207/black/average"
)

var _ = Describe("Average", func() {

	data := []float64{1, 2, 1, 2, 1, 2}

	Describe("Sum", func() {

		result := Sum(data)
		It("Should compute the sum", func() {
			Ω(result).Should(Equal(float64(9)))
		})

	})

	Describe("Rolling average", func() {

		averages := Rolling(data, 2)

		It("should have a length equal to data", func() {
			Ω(len(averages)).Should(Equal(6))
		})

		It("should have correct values", func() {
			Ω(averages[0]).Should(Equal(float64(1)))
			Ω(averages[1]).Should(Equal(1.5))
			Ω(averages[2]).Should(Equal(1.5))
			Ω(averages[3]).Should(Equal(1.5))
			Ω(averages[4]).Should(Equal(1.5))
		})

	})

})
