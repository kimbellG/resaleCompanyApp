package models

type Alternative struct {
	Id              int
	Name            string
	Description     string
	Marks           map[int]float32
	SpecificWeights map[int]float32
	Weight          float32
}

type AlternativeInput struct {
	ProblemId   int
	Name        string
	Description string
}

type AlternativeMarksInput struct {
	AlternativeId int
	ExpertLogin   string
	Marks         []float32
}
