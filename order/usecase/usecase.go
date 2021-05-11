package usecase

import (
	"context"
	"cw/models"
	"cw/order"
	"fmt"
	"time"
)

type OrderUseCase struct {
	repo    order.Repository
	client  order.ClientController
	offer   order.OfferController
	manager order.ManagerController
}

func NewOrderUseCase(rep order.Repository, cl order.ClientController, off order.OfferController, man order.ManagerController) *OrderUseCase {
	return &OrderUseCase{
		repo:    rep,
		client:  cl,
		offer:   off,
		manager: man,
	}
}

func (o *OrderUseCase) Add(ctx context.Context, order *models.Order) error {

	return fmt.Errorf("usecase add: %v", o.repo.Add(ctx, order))
}

func (o *OrderUseCase) Gets(ctx context.Context) ([]order.OrderOutput, error) {
	result, err := o.repo.Gets(ctx)
	if err != nil {
		return nil, fmt.Errorf("usecase: %v", err)
	}

	return o.arrayModdelsToOrderOutput(result)
}

func (o *OrderUseCase) arrayModdelsToOrderOutput(models []models.Order) ([]order.OrderOutput, error) {
	result := make([]order.OrderOutput, 0)

	for _, val := range models {
		tmp, err := o.modToOrderOutput(&val)
		if err != nil {
			return nil, fmt.Errorf("translate array: %v", err)
		}

		result = append(result, *tmp)
	}

	return result, nil
}

func (o *OrderUseCase) modToOrderOutput(mod *models.Order) (*order.OrderOutput, error) {
	result := &order.OrderOutput{}
	result.Id = mod.Id

	for _, id := range mod.Offers {
		offer, err := o.offer.GetById(id)
		if err != nil {
			return nil, fmt.Errorf("translate offers: %v", err)
		}

		result.Offers = append(result.Offers, *offer)
	}

	client, err := o.client.GetById(mod.ClientId)
	if err != nil {
		return nil, fmt.Errorf("client translate: %v", err)
	}

	managerName, err := o.manager.GetNameByLogin(mod.ManagerLogin)
	if err != nil {
		return nil, fmt.Errorf("manager translate: %v", err)
	}

	result.Client = *client
	result.Manager = managerName
	result.OrderDate = mod.OrderDate
	result.Quantity = mod.Quantity
	result.Status = mod.Status

	return result, nil
}

func (o *OrderUseCase) GetInInterval(ctx context.Context, start, end time.Time) ([]order.OrderOutput, error) {
	result, err := o.repo.GetInInterval(ctx, start.Format("2004-10-19 10:23:54"), end.Format("2004-10-19 10:23:54"))
	if err != nil {
		return nil, fmt.Errorf("repo: %v", err)
	}

	return o.arrayModdelsToOrderOutput(result)
}

func (o *OrderUseCase) UpdateStatus(ctx context.Context, id int, newStatus string) error {
	if err := o.repo.UpdateStatus(ctx, id, newStatus); err != nil {
		return fmt.Errorf("repo: %v", err)
	}

	return nil
}

func (o *OrderUseCase) Filter(ctx context.Context, key string, value interface{}) ([]order.OrderOutput, error) {
	result, err := o.repo.Filter(ctx, key, value)
	if err != nil {
		return nil, fmt.Errorf("repo: %v", err)
	}

	return o.arrayModdelsToOrderOutput(result)
}
