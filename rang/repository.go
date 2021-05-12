package rang

import (
	"context"
	"cw/models"
)

type Repository interface {
	AddProblem(ctx context.Context, problem *models.ProblemInput) error
	AddAlternative(ctx context.Context, alternative *models.AlternativeInput) error
	Gets(ctx context.Context) ([]*models.Problem, error)
	GetProblemReport(ctx context.Context, id int) (*models.Problem, error)
	UpdateMarks(ctx context.Context, newMarks *models.AlternativeMarksInput) error
	DeleteProblem(ctx context.Context, id int) error
}
