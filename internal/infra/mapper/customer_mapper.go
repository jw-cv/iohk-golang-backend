package mapper

import (
	"strconv"
	"time"

	"iohk-golang-backend-preprod/ent"
	"iohk-golang-backend-preprod/ent/customer"
	"iohk-golang-backend-preprod/graph/model"
	domainmodel "iohk-golang-backend-preprod/internal/domain/model"
)

func DomainToGraphQL(c *domainmodel.Customer) *model.Customer {
	return &model.Customer{
		ID:         strconv.Itoa(c.ID),
		Name:       c.Name,
		Surname:    c.Surname,
		Number:     c.Number,
		Gender:     model.Gender(c.Gender),
		Country:    c.Country,
		Dependants: c.Dependants,
		BirthDate:  c.BirthDate.Format("2006-01-02"),
	}
}

func GraphQLToDomain(gc *model.Customer) *domainmodel.Customer {
	id, _ := strconv.Atoi(gc.ID)
	birthDate, _ := time.Parse("2006-01-02", gc.BirthDate)
	return &domainmodel.Customer{
		ID:         id,
		Name:       gc.Name,
		Surname:    gc.Surname,
		Number:     gc.Number,
		Gender:     domainmodel.Gender(gc.Gender),
		Country:    gc.Country,
		Dependants: gc.Dependants,
		BirthDate:  birthDate,
	}
}

func CreateInputToDomain(input *model.CreateCustomerInput) *domainmodel.Customer {
	birthDate, _ := time.Parse("2006-01-02", input.BirthDate)
	return &domainmodel.Customer{
		Name:       input.Name,
		Surname:    input.Surname,
		Number:     input.Number,
		Gender:     domainmodel.Gender(input.Gender),
		Country:    input.Country,
		Dependants: input.Dependants,
		BirthDate:  birthDate,
	}
}

func EntToDomain(c *ent.Customer) *domainmodel.Customer {
	return &domainmodel.Customer{
		ID:         c.ID,
		Name:       c.Name,
		Surname:    c.Surname,
		Number:     c.Number,
		Gender:     domainmodel.Gender(c.Gender),
		Country:    c.Country,
		Dependants: c.Dependants,
		BirthDate:  c.BirthDate,
	}
}

func DomainToEnt(c *domainmodel.Customer) *ent.Customer {
	return &ent.Customer{
		ID:         c.ID,
		Name:       c.Name,
		Surname:    c.Surname,
		Number:     c.Number,
		Gender:     customer.Gender(c.Gender),
		Country:    c.Country,
		Dependants: c.Dependants,
		BirthDate:  c.BirthDate,
	}
}

func DomainToUpdateInput(c *domainmodel.Customer) *model.UpdateCustomerInput {
	return &model.UpdateCustomerInput{
		Name:       &c.Name,
		Surname:    &c.Surname,
		Number:     &c.Number,
		Gender:     (*model.Gender)(&c.Gender),
		Country:    &c.Country,
		Dependants: &c.Dependants,
		BirthDate:  stringPtr(c.BirthDate.Format("2006-01-02")),
	}
}

// Helper function for string pointers
func stringPtr(s string) *string {
	return &s
}

func DomainToGraphQLSlice(customers []*domainmodel.Customer) []*model.Customer {
	result := make([]*model.Customer, len(customers))
	for i, c := range customers {
		result[i] = DomainToGraphQL(c)
	}
	return result
}

// Add this new function
func UpdateInputToDomain(id string, input *model.UpdateCustomerInput) *domainmodel.Customer {
	customer := &domainmodel.Customer{}
	customer.ID, _ = strconv.Atoi(id)

	if input.Name != nil {
		customer.Name = *input.Name
	}
	if input.Surname != nil {
		customer.Surname = *input.Surname
	}
	if input.Number != nil {
		customer.Number = *input.Number
	}
	if input.Gender != nil {
		customer.Gender = domainmodel.Gender(*input.Gender)
	}
	if input.Country != nil {
		customer.Country = *input.Country
	}
	if input.Dependants != nil {
		customer.Dependants = *input.Dependants
	}
	if input.BirthDate != nil {
		customer.BirthDate, _ = time.Parse("2006-01-02", *input.BirthDate)
	}

	return customer
}
