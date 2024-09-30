package graph

import (
	"context"
	"iohk-golang-backend/graph/model"
	"iohk-golang-backend/internal/domain/service"
	"iohk-golang-backend/internal/infra/mapper"
)

type Resolver struct {
	customerService service.CustomerService
}

func NewResolver(customerService service.CustomerService) *Resolver {
	return &Resolver{
		customerService: customerService,
	}
}

// Query Resolvers
func (r *queryResolver) Customer(ctx context.Context, id string) (*model.Customer, error) {
	domainCustomer, err := r.customerService.GetCustomer(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapper.DomainToGraphQL(domainCustomer), nil
}

func (r *queryResolver) Customers(ctx context.Context) ([]*model.Customer, error) {
	domainCustomers, err := r.customerService.GetAllCustomers(ctx)
	if err != nil {
		return nil, err
	}
	return mapper.DomainToGraphQLSlice(domainCustomers), nil
}

// Mutation Resolvers
func (r *mutationResolver) CreateCustomer(ctx context.Context, input model.CreateCustomerInput) (*model.Customer, error) {
	domainCustomer := mapper.CreateInputToDomain(&input)
	createdCustomer, err := r.customerService.CreateCustomer(ctx, domainCustomer)
	if err != nil {
		return nil, err
	}
	return mapper.DomainToGraphQL(createdCustomer), nil
}

func (r *mutationResolver) UpdateCustomer(ctx context.Context, id string, input model.UpdateCustomerInput) (*model.Customer, error) {
	domainCustomer := mapper.UpdateInputToDomain(id, &input)
	updatedCustomer, err := r.customerService.UpdateCustomer(ctx, id, domainCustomer)
	if err != nil {
		return nil, err
	}
	return mapper.DomainToGraphQL(updatedCustomer), nil
}

func (r *mutationResolver) DeleteCustomer(ctx context.Context, id string) (bool, error) {
	return r.customerService.DeleteCustomer(ctx, id)
}

// Resolver type assertions
func (r *Resolver) Query() QueryResolver       { return &queryResolver{r} }
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

type queryResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
