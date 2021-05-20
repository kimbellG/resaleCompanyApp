package usecase

import (
	"context"
	"cw/app"
	"cw/models"
	"cw/rang"
	"fmt"
)

type UserController interface {
	GetIdByLogin(login string) (int, error)
}

type RangUseCase struct {
	repo rang.Repository
	user UserController
}

func NewRangUseCase(rep rang.Repository, usr UserController) *RangUseCase {
	return &RangUseCase{
		repo: rep,
		user: usr,
	}
}

func (r *RangUseCase) Add(ctx context.Context, problem *models.ProblemInput) error {
	if err := r.repo.Add(ctx, problem); err != nil {
		return fmt.Errorf("repo: %v", err)
	}

	return nil
}

func (r *RangUseCase) AddAlternative(ctx context.Context, problemName string, alternative *models.AlternativeInput) error {

	id, err := r.getProblemByName(ctx, problemName)
	if err != nil {
		return fmt.Errorf("get problem's id: %v", err)
	}

	if err := r.repo.AddAlternative(ctx, id, alternative); err != nil {
		return fmt.Errorf("repo: %v", err)
	}

	return nil
}

func (r *RangUseCase) getProblemByName(ctx context.Context, name string) (int, error) {
	problems, err := r.Gets(ctx)
	if err != nil {
		return -1, fmt.Errorf("get all problems: %v", err)
	}

	for _, problem := range problems {
		if problem.Name == name {
			return problem.Id, nil
		}
	}

	return -1, fmt.Errorf("not found")
}

func (r *RangUseCase) AddMarkByNames(ctx context.Context, mark *rang.AlternativeMarkInput) error {
	problemID, alternativeID, err := r.getProblemAlternativeIDsByName(ctx, mark)
	if err != nil {
		return fmt.Errorf("get problem and alternative IDs: %v", err)
	}

	err = r.AddAlternativeMark(ctx,
		&models.AlternativeMarkInput{
			ProblemId:     problemID,
			AlternativeId: alternativeID,
			ExpertLogin:   mark.ExpertLogin,
			Mark:          mark.Mark,
		},
	)

	if err != nil {
		return fmt.Errorf("add alternative mark with ids: %v", err)
	}

	return nil
}

func (r *RangUseCase) getProblemAlternativeIDsByName(ctx context.Context, mark *rang.AlternativeMarkInput) (int, int, error) {
	problems, err := r.Gets(ctx)
	if err != nil {
		return -1, -1, fmt.Errorf("get all problem: %v", err)
	}

	for _, problem := range problems {
		if problem.Name == mark.ProblemName {
			for _, alternative := range problem.Alternatives {
				if alternative.Name == mark.AlternativeName {
					return problem.Id, alternative.Id, nil
				}
			}
		}
	}

	return -1, -1, fmt.Errorf("not found")
}

func (r *RangUseCase) AddAlternativeMark(ctx context.Context, alternativeMarks *models.AlternativeMarkInput) error {
	problem, err := r.repo.GetProblemReport(ctx, alternativeMarks.ProblemId)
	if err != nil {
		return fmt.Errorf("get problem repo: %v", err)
	}

	alternative, err := findAlternative(problem, alternativeMarks.AlternativeId)
	if err != nil {
		return fmt.Errorf("find alternative in problem: %v", err)
	}

	alternativeMarks.ExpertLogin = ctx.Value(app.UserInfo).(string)

	expertId, err := r.user.GetIdByLogin(alternativeMarks.ExpertLogin)
	if err != nil {
		return fmt.Errorf("get user id: %v", err)
	}

	alternative.Marks[expertId] = alternativeMarks.Mark
	if isAllMarksUserInput(problem.Alternatives, expertId) {
		rangMethod(problem)
	}

	if err := r.repo.AddAlternativeMark(ctx, problem); err != nil {
		return fmt.Errorf("repo: %v", err)
	}

	return nil
}

func findAlternative(problem *models.Problem, alternativeId int) (*models.Alternative, error) {
	for _, alternative := range problem.Alternatives {
		if alternative.Id == alternativeId {
			return alternative, nil
		}
	}

	return nil, fmt.Errorf("invalid alternative id")
}

func isAllMarksUserInput(alternatives []*models.Alternative, userId int) bool {
	for _, alternative := range alternatives {
		if _, ok := alternative.Marks[userId]; !ok {
			return false
		}
	}

	return true
}

func rangMethod(problem *models.Problem) {
	proccesingSpetificWeights(problem.Alternatives)
	proccesingWeight(problem.Alternatives)
}

func proccesingSpetificWeights(alternatives []*models.Alternative) {
	for _, alternative := range alternatives {
		for id, mark := range alternative.Marks {
			if !isAllMarksUserInput(alternatives, id) {
				continue
			}
			alternative.SpecificWeights[id] = mark / SumAlternativeMarksFromExpert(alternatives, id)
		}
	}
}

func SumAlternativeMarksFromExpert(alternatives []*models.Alternative, expertId int) float32 {
	result := float32(0)
	for _, alternative := range alternatives {
		result += alternative.Marks[expertId]
	}

	return result
}

func proccesingWeight(alternatives []*models.Alternative) {
	for _, alternative := range alternatives {
		sum := float32(0)
		for _, mark := range alternative.SpecificWeights {
			sum += mark
		}
		alternative.Weight = sum / float32(getExpertCount(alternative.SpecificWeights))
	}
}

func getExpertCount(specificWeight map[int]float32) int {
	result := 0
	for _, value := range specificWeight {
		if value != 0 {
			result++
		}
	}

	return result
}

func getResult(alternatives []*models.Alternative) *models.Alternative {
	result := &models.Alternative{}
	maxWeight := alternatives[0].Weight
	for _, alternative := range alternatives {
		if maxWeight < alternative.Weight {
			result = alternative
		}
	}

	return result
}

func (r *RangUseCase) GetProblemReport(ctx context.Context, id int) (*models.Problem, error) {
	result, err := r.repo.GetProblemReport(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("repo: %v", err)
	}

	return result, nil
}

func (r *RangUseCase) Gets(ctx context.Context) ([]*models.Problem, error) {
	result, err := r.repo.Gets(ctx)
	if err != nil {
		return nil, fmt.Errorf("repo: %v", err)
	}

	return result, nil
}
