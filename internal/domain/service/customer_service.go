package service

import (
	"context"

	"iohk-golang-backend-preprod/graph/model"
	"iohk-golang-backend-preprod/internal/domain/repository"
)

type CustomerService interface {
	CreateCustomer(ctx context.Context, input model.CreateCustomerInput) (*model.Customer, error)
	GetCustomer(ctx context.Context, id string) (*model.Customer, error)
	GetAllCustomers(ctx context.Context) ([]*model.Customer, error)
	UpdateCustomer(ctx context.Context, id string, input model.UpdateCustomerInput) (*model.Customer, error)
	DeleteCustomer(ctx context.Context, id string) (bool, error)
}

type customerService struct {
	repo repository.CustomerRepository
}

func NewCustomerService(repo repository.CustomerRepository) CustomerService {
	return &customerService{repo: repo}
}

func (s *customerService) CreateCustomer(ctx context.Context, input model.CreateCustomerInput) (*model.Customer, error) {
	return s.repo.Create(ctx, &input)
}

func (s *customerService) GetCustomer(ctx context.Context, id string) (*model.Customer, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *customerService) GetAllCustomers(ctx context.Context) ([]*model.Customer, error) {
	return s.repo.GetAll(ctx)
}

func (s *customerService) UpdateCustomer(ctx context.Context, id string, input model.UpdateCustomerInput) (*model.Customer, error) {
	return s.repo.Update(ctx, id, &input)
}

func (s *customerService) DeleteCustomer(ctx context.Context, id string) (bool, error) {
	err := s.repo.Delete(ctx, id)
	return err == nil, err
}
