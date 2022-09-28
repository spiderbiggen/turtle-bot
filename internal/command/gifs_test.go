package command

import (
	"math/rand"
	"testing"
)

func TestWeightedArgument(t *testing.T) {
	runs := 100_000
	tolerance := 0.01
	weight := 0.0
	list := Args{
		{Query: "sleep", Weight: 80},
		{Query: "night", Weight: 70},
		{Query: "froggers", Weight: 1, GifCount: 1, IsSearch: true},
	}
	for _, argument := range list {
		weight += float64(argument.Weight)
	}

	var sum int64
	for i := 0; i < runs; i++ {
		rand.Seed(int64(31 * i))
		j := 1
		for ; j <= 10_000; j++ {
			a := list.Pick()
			if a.Query == "froggers" {
				sum += int64(j)
				j = -1
				break
			}
		}
		if j != -1 {
			t.Logf("%05d Didn't occur in 10,000 rounds", i)
		}
	}
	avg := float64(sum) / float64(runs)
	t.Logf("Avg: %.2f", avg)
	if avg < weight-weight*tolerance || avg > weight+weight*tolerance {
		t.Errorf("Average %.2f not within %.0f of %.0f", avg, tolerance*100, weight)
	}
}
