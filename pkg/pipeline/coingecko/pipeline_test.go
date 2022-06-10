package coingecko

import (
	"fmt"
	"testing"
)

func TestPrepareURL(t *testing.T) {
	date := "20-04-2022"
	cases := []struct {
		name     string
		pipeline *Pipeline
		expected string
	}{
		{
			name:     "Happy path",
			pipeline: NewPipeline("ethereum"),
			expected: fmt.Sprintf("https://api.coingecko.com/api/v3/coins/ethereum/history?date=%s&localization=false", date),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := c.pipeline.prepareURL(date)
			if actual != c.expected {
				t.Fail()
			}
		})

	}
}

func TestProcess(t *testing.T) {

}
