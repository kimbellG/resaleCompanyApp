package rang

import (
	"context"
	"cw/models"
)

type Repository interface {
	Add(ctx context.Context, problem *models.ProblemInput) error
	AddAlternativeMark(ctx context.Context, problem *models.Problem) error
	Gets(ctx context.Context) ([]*models.Problem, error)
	GetProblemReport(ctx context.Context, id int) (*models.Problem, error)
}
