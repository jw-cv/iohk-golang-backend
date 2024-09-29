package service

import (
	"context"

	domainmodel "iohk-golang-backend-preprod/internal/domain/model"
	"iohk-golang-backend-preprod/internal/domain/repository"
	"iohk-golang-backend-preprod/internal/infra/mapper"
)

type CustomerService interface {
	CreateCustomer(ctx context.Context, customer *domainmodel.Customer) (*domainmodel.Customer, error)
	GetCustomer(ctx context.Context, id string) (*domainmodel.Customer, error)
	GetAllCustomers(ctx context.Context) ([]*domainmodel.Customer, error)
	UpdateCustomer(ctx context.Context, id string, customer *domainmodel.Customer) (*domainmodel.Customer, error)
	DeleteCustomer(ctx context.Context, id string) (bool, error)
}

type customerService struct {
	repo repository.CustomerRepository
}

func NewCustomerService(repo repository.CustomerRepository) CustomerService {
	return &customerService{repo: repo}
}

func (s *customerService) CreateCustomer(ctx context.Context, customer *domainmodel.Customer) (*domainmodel.Customer, error) {
	return s.repo.Create(ctx, customer)
}

func (s *customerService) GetCustomer(ctx context.Context, id string) (*domainmodel.Customer, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *customerService) GetAllCustomers(ctx context.Context) ([]*domainmodel.Customer, error) {
	return s.repo.GetAll(ctx)
}

func (s *customerService) UpdateCustomer(ctx context.Context, id string, customer *domainmodel.Customer) (*domainmodel.Customer, error) {
	input := mapper.DomainToUpdateInput(customer)
	return s.repo.Update(ctx, id, input)
}

func (s *customerService) DeleteCustomer(ctx context.Context, id string) (bool, error) {
	err := s.repo.Delete(ctx, id)
	return err == nil, err
}
