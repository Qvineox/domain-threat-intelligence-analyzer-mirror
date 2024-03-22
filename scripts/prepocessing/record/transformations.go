package record

import (
	"github.com/montanaflynn/stats"
	"slices"
)

func levelsMAD(parts []string) (float64, error) {
	var lengths []float64

	for _, p := range parts {
		lengths = append(lengths, float64(len(p)))
	}

	slices.Sort(lengths)

	return stats.MedianAbsoluteDeviation(lengths)
}
