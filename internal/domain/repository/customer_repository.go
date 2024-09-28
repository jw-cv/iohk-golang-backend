package repository

import (
	"context"
	"strconv"
	"time"

	"iohk-golang-backend-preprod/ent"
	"iohk-golang-backend-preprod/ent/customer"
	"iohk-golang-backend-preprod/graph/model"
)

type CustomerRepository interface {
	Create(ctx context.Context, input *model.CreateCustomerInput) (*model.Customer, error)
	GetByID(ctx context.Context, id string) (*model.Customer, error)
	GetAll(ctx context.Context) ([]*model.Customer, error)
	Update(ctx context.Context, id string, input *model.UpdateCustomerInput) (*model.Customer, error)
	Delete(ctx context.Context, id string) error
}

type customerRepository struct {
	client *ent.Client
}

func NewCustomerRepository(client *ent.Client) CustomerRepository {
	return &customerRepository{client: client}
}

func (r *customerRepository) Create(ctx context.Context, input *model.CreateCustomerInput) (*model.Customer, error) {
	birthDate, err := time.Parse("2006-01-02", input.BirthDate)
	if err != nil {
		return nil, err
	}

	c, err := r.client.Customer.
		Create().
		SetName(input.Name).
		SetSurname(input.Surname).
		SetNumber(input.Number).
		SetGender(customer.Gender(input.Gender)).
		SetCountry(input.Country).
		SetDependants(input.Dependants).
		SetBirthDate(birthDate).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return entToGraphQL(c), nil
}

func (r *customerRepository) GetByID(ctx context.Context, id string) (*model.Customer, error) {
	customerID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	c, err := r.client.Customer.Get(ctx, customerID)
	if err != nil {
		return nil, err
	}

	return entToGraphQL(c), nil
}

func (r *customerRepository) GetAll(ctx context.Context) ([]*model.Customer, error) {
	customers, err := r.client.Customer.Query().All(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*model.Customer, len(customers))
	for i, c := range customers {
		result[i] = entToGraphQL(c)
	}
	return result, nil
}

func (r *customerRepository) Update(ctx context.Context, id string, input *model.UpdateCustomerInput) (*model.Customer, error) {
	customerID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	update := r.client.Customer.UpdateOneID(customerID)

	if input.Name != nil {
		update.SetName(*input.Name)
	}
	if input.Surname != nil {
		update.SetSurname(*input.Surname)
	}
	if input.Number != nil {
		update.SetNumber(*input.Number)
	}
	if input.Gender != nil {
		update.SetGender(customer.Gender(*input.Gender))
	}
	if input.Country != nil {
		update.SetCountry(*input.Country)
	}
	if input.Dependants != nil {
		update.SetDependants(*input.Dependants)
	}
	if input.BirthDate != nil {
		birthDate, err := time.Parse("2006-01-02", *input.BirthDate)
		if err != nil {
			return nil, err
		}
		update.SetBirthDate(birthDate)
	}

	c, err := update.Save(ctx)
	if err != nil {
		return nil, err
	}

	return entToGraphQL(c), nil
}

func (r *customerRepository) Delete(ctx context.Context, id string) error {
	customerID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	return r.client.Customer.DeleteOneID(customerID).Exec(ctx)
}

func entToGraphQL(c *ent.Customer) *model.Customer {
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
