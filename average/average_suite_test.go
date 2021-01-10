package average_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAverage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Average Suite")
}
