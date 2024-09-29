package repository

import (
	"context"
	"strconv"
	"time"

	"iohk-golang-backend-preprod/ent"
	"iohk-golang-backend-preprod/graph/model"
	domainmodel "iohk-golang-backend-preprod/internal/domain/model"
	"iohk-golang-backend-preprod/internal/infra/mapper"
)

type CustomerRepository interface {
	Create(ctx context.Context, customer *domainmodel.Customer) (*domainmodel.Customer, error)
	GetByID(ctx context.Context, id string) (*domainmodel.Customer, error)
	GetAll(ctx context.Context) ([]*domainmodel.Customer, error)
	Update(ctx context.Context, id string, input *model.UpdateCustomerInput) (*domainmodel.Customer, error)
	Delete(ctx context.Context, id string) error
}

type customerRepository struct {
	client *ent.Client
}

func NewCustomerRepository(client *ent.Client) CustomerRepository {
	return &customerRepository{client: client}
}

func (r *customerRepository) Create(ctx context.Context, customer *domainmodel.Customer) (*domainmodel.Customer, error) {
	entCustomer, err := r.client.Customer.
		Create().
		SetName(customer.Name).
		SetSurname(customer.Surname).
		SetNumber(customer.Number).
		SetGender(customer.Gender.ToEntGender()).
		SetCountry(customer.Country).
		SetDependants(customer.Dependants).
		SetBirthDate(customer.BirthDate).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return mapper.EntToDomain(entCustomer), nil
}

func (r *customerRepository) GetByID(ctx context.Context, id string) (*domainmodel.Customer, error) {
	customerID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	c, err := r.client.Customer.Get(ctx, customerID)
	if err != nil {
		return nil, err
	}

	return mapper.EntToDomain(c), nil
}

func (r *customerRepository) GetAll(ctx context.Context) ([]*domainmodel.Customer, error) {
	customers, err := r.client.Customer.Query().All(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*domainmodel.Customer, len(customers))
	for i, c := range customers {
		result[i] = mapper.EntToDomain(c)
	}
	return result, nil
}

func (r *customerRepository) Update(ctx context.Context, id string, input *model.UpdateCustomerInput) (*domainmodel.Customer, error) {
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
		update.SetGender(domainmodel.Gender(*input.Gender).ToEntGender())
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

	return mapper.EntToDomain(c), nil
}

func (r *customerRepository) Delete(ctx context.Context, id string) error {
	customerID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	return r.client.Customer.DeleteOneID(customerID).Exec(ctx)
}
