package model

import (
	"time"

	"iohk-golang-backend-preprod/ent"
	"iohk-golang-backend-preprod/ent/customer"
	"iohk-golang-backend-preprod/graph/model"
)

type Customer struct {
	ID         int
	Name       string
	Surname    string
	Number     int
	Gender     customer.Gender
	Country    string
	Dependants int
	BirthDate  time.Time
}

func CustomerFromEnt(c *ent.Customer) *Customer {
	return &Customer{
		ID:         c.ID,
		Name:       c.Name,
		Surname:    c.Surname,
		Number:     c.Number,
		Gender:     c.Gender,
		Country:    c.Country,
		Dependants: c.Dependants,
		BirthDate:  c.BirthDate,
	}
}

func (c *Customer) ToEnt() *ent.Customer {
	return &ent.Customer{
		ID:         c.ID,
		Name:       c.Name,
		Surname:    c.Surname,
		Number:     c.Number,
		Gender:     c.Gender,
		Country:    c.Country,
		Dependants: c.Dependants,
		BirthDate:  c.BirthDate,
	}
}

// Add these helper functions
func GenderFromGraphQL(g model.Gender) customer.Gender {
	return customer.Gender(g)
}

func GenderToGraphQL(g customer.Gender) model.Gender {
	return model.Gender(g)
}
