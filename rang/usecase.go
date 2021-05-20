package rang

import (
	"context"
	"cw/models"
)

type AlternativeMarkInput struct {
	ProblemName     string
	AlternativeName string
	ExpertLogin     string
	Mark            float32
}

type Usecase interface {
	Add(ctx context.Context, problem *models.ProblemInput) error
	AddAlternative(ctx context.Context, problemName string, problem *models.AlternativeInput) error
	AddAlternativeMark(ctx context.Context, marks *models.AlternativeMarkInput) error
	AddMarkByNames(ctx context.Context, marks *AlternativeMarkInput) error

	GetProblemReport(ctx context.Context, id int) (*models.Problem, error)
	Gets(ctx context.Context) ([]*models.Problem, error)
}
