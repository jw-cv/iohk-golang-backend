package graph

import (
	"iohk-golang-backend-preprod/internal/domain/service"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	customerService service.CustomerService
}

func NewResolver(customerService service.CustomerService) *Resolver {
	return &Resolver{
		customerService: customerService,
	}
}
