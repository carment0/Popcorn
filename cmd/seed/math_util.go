// Copyright (c) 2018 Popcorn
// Author(s) Calvin Feng

package main

func Average(numbers []float64) float64 {
	if len(numbers) == 0 {
		return 0.0
	}

	sum := 0.0
	for i := 0; i < len(numbers); i += 1 {
		sum += numbers[i]
	}

	return sum / float64(len(numbers))
}
