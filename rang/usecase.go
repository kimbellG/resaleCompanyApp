package rang

import (
	"context"
	"cw/models"
)

type Usecase interface {
	AddProblem(ctx context.Context, problem *models.ProblemInput) error
	AddAlternativeMarks(ctx context.Context, marks *models.AlternativeInput) error
	GetProblemReport(ctx context.Context, id int) (*models.Problem, error)
	Gets(ctx context.Context) ([]*models.Problem, error)
	UpdateMarks(ctx context.Context, newMarks *models.AlternativeInput) error
	DeleteProblem(ctx context.Context, problemId int) error
}
