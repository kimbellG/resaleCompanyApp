package rang

import (
	"context"
	"cw/models"
)

type Usecase interface {
	Add(ctx context.Context, problem *models.ProblemInput) error
	AddAlternativeMark(ctx context.Context, marks *models.AlternativeMarkInput) error
	GetProblemReport(ctx context.Context, id int) (*models.Problem, error)
	Gets(ctx context.Context) ([]*models.Problem, error)
}
