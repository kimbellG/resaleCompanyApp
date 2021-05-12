package models

type Problem struct {
	Id           int
	Name         string
	Description  string
	Alternatives []*Alternative
	Result       *Alternative
}

type ProblemInput struct {
	Name         string
	Description  string
	Alternatives []*AlternativeInput
}
