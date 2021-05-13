package usecase

import (
	"cw/models"
	"math"
	"testing"
)

func TestRankMethod(t *testing.T) {
	alter := []*models.Alternative{{Marks: map[int]float32{1: 10, 2: 8}, SpecificWeights: map[int]float32{}}, {Marks: map[int]float32{1: 7, 2: 6}, SpecificWeights: map[int]float32{}},
		{Marks: map[int]float32{1: 9, 2: 10}, SpecificWeights: map[int]float32{}}, {Marks: map[int]float32{1: 3, 2: 4}, SpecificWeights: map[int]float32{}},
		{Marks: map[int]float32{1: 4, 2: 2}, SpecificWeights: map[int]float32{}}, {Marks: map[int]float32{1: 5, 2: 7}, SpecificWeights: map[int]float32{}},
	}

	testProblem := &models.Problem{
		Id:           1,
		Name:         "test",
		Description:  "test test",
		Alternatives: alter,
	}

	currentWeight := []float64{0.239, 0.173, 0.254, 0.093, 0.079, 0.162}

	rangMethod(testProblem)

	for i, alternative := range testProblem.Alternatives {
		wfl := float64(alternative.Weight)
		w := math.Floor(wfl*1000) / 1000
		if w != currentWeight[i] {
			t.Errorf("incorrect weight: %v != %v", w, currentWeight[i])
		}
	}
}
