package graph

import "iohk-golang-backend-preprod/graph/model"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	customers []*model.Customer
	nextID    int
}

func NewResolver() *Resolver {
	return &Resolver{
		customers: make([]*model.Customer, 0),
		nextID:    1,
	}
}
