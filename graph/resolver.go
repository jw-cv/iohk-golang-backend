package graph

import (
	"iohk-golang-backend-preprod/ent"
	"iohk-golang-backend-preprod/internal/domain/repository"
	"iohk-golang-backend-preprod/internal/domain/service"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	customerService service.CustomerService
}

func NewResolver(client *ent.Client) *Resolver {
	customerRepo := repository.NewCustomerRepository(client)
	customerService := service.NewCustomerService(customerRepo)
	return &Resolver{
		customerService: customerService,
	}
}
