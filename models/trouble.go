package models

type Problem struct {
	Id           int            `json:"id"`
	Name         string         `json:"name"`
	Description  string         `json:"decriprion"`
	Alternatives []*Alternative `json:"alternatives"`
}

type ProblemInput struct {
	Name         string
	Description  string
	Alternatives []*AlternativeInput
}
