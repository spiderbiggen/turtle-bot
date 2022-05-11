package command

import (
	"math/rand"
	"testing"
)

func TestCommand(t *testing.T) {
	runs := 100_000
	tolerance := 0.01
	weight := 0.0
	list := Args{
		&WeightedArgument{Query: "sleep", Weight: 80},
		&WeightedArgument{Query: "night", Weight: 70},
		&WeightedArgument{Query: "froggers", Weight: 1, GifCount: 1, IsSearch: true},
	}
	for _, argument := range list {
		weight += float64(argument.Weight)
	}

	results := make([]int, runs)
	var sum int64
	for i := 0; i < runs; i++ {
		rand.Seed(int64(i))
		j := 1
		for ; j <= 10_000; j++ {
			a := list.Random()
			if a.Query == "froggers" {
				results[i] = j
				sum += int64(j)
				j = -1
				break
			}
		}
		if j != -1 {
			results[i] = -1
		}
	}
	//t.Logf("%+v", results)
	avg := float64(sum) / float64(runs)
	t.Logf("Avg: %.2f", avg)
	if avg < weight-weight*tolerance || avg > weight+weight*tolerance {
		t.Errorf("Average %.2f not within %.0f of %.0f", avg, tolerance*100, weight)
	}
}
