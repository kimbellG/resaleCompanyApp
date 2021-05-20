package models

type Alternative struct {
	Id              int             `json:"id"`
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	Marks           map[int]float32 `json:"marks"`
	SpecificWeights map[int]float32 `json:"specific_weights"`
	Weight          float32         `json:"weight"`
}

type AlternativeInput struct {
	Name        string
	Description string
}

type AlternativeMarkInput struct {
	ProblemId     int
	AlternativeId int
	ExpertLogin   string
	Mark          float32
}
