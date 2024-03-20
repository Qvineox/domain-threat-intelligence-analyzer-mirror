package record

import (
	"errors"
)

func countSymbols(name string) (totalCount, vowelsCount, consonantsCount, numericsCount, specialsCount, pointsCount int, vowelsRatio, consonantsRatio, numericsRatio, specialsRatio, pointsRatio float64, err error) {
	totalCount = len(name)

	if totalCount == 0 {
		return 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, errors.New("empty name passed")
	}

	for _, s := range name {
		switch s {
		case 'a', 'e', 'i', 'o', 'u':
			vowelsCount++
		case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
			numericsCount++
		case 'b', 'c', 'd', 'f', 'g', 'h', 'j', 'k', 'l', 'm', 'n', 'p', 'q', 'r', 's', 't', 'v', 'w', 'x', 'y', 'z':
			consonantsCount++
		case '-', '_':
			specialsCount++
		case '.':
			pointsCount++
		default:
			return 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, errors.New("empty name passed")
		}
	}

	vowelsRatio = float64(vowelsCount) / float64(totalCount)
	consonantsRatio = float64(consonantsCount) / float64(totalCount)
	numericsRatio = float64(numericsCount) / float64(totalCount)
	specialsRatio = float64(specialsCount) / float64(totalCount)
	pointsRatio = float64(pointsCount) / float64(totalCount)

	return totalCount, vowelsCount, consonantsCount, numericsCount, specialsCount, pointsCount, vowelsRatio, consonantsRatio, numericsRatio, specialsRatio, pointsRatio, nil
}

func maxRepeatedSymbolCount(name string) int {
	var count = 0

	for i1, v1 := range name {
		var inRow = 0

		for _, v2 := range name[i1:] {
			if v1 == v2 {
				inRow++
			} else {
				if inRow > count {
					count = inRow
				}
				break
			}
		}
	}

	return count
}

func uniqueSymbolsRation(name string) float64 {
	var unique = make(map[int32]bool)

	for _, s := range name {
		_, ok := unique[s]
		if !ok {
			unique[s] = false
		}
	}

	return float64(len(unique)) / float64(len(name))
}
