package command

import (
	"math/rand"
	"testing"
	"time"
)

func TestWeightedArgument(t *testing.T) {
	runs := 100_000
	tolerance := 0.01
	weight := 0.0
	list := Args{
		{Url: "https://tenor.com/view/frog-dance-animation-cute-funny-gif-17184624"},
		{Query: "sleep", Weight: 20},
		{Query: "dogsleep", Weight: 20},
		{Query: "catsleep", Weight: 20},
		{Query: "rabbitsleep", Weight: 20},
		{Query: "ratsleep", Weight: 20},
		{Query: "ducksleep", Weight: 20},
		{Query: "animalsleep", Weight: 20},
	}
	for _, argument := range list {
		weight += float64(argument.NormalizedWeight())
	}

	offset := time.Now().UnixMilli()
	var sum int64
	for i := 0; i < runs; i++ {
		rand.Seed(int64(31*i) + offset)
		j := 1
		for ; j <= 10_000; j++ {
			a := list.Pick()
			if a.Url == "https://tenor.com/view/frog-dance-animation-cute-funny-gif-17184624" {
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
	t.Logf("Sum: %.2f, Avg: %.2f", weight, avg)
	if avg < weight-weight*tolerance || avg > weight+weight*tolerance {
		t.Errorf("Average %.2f not within %.0f of %.0f", avg, tolerance*100, weight)
	}
}
