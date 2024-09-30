//go:build testcoverage
// +build testcoverage

package repository

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"iohk-golang-backend-preprod/ent"
	"iohk-golang-backend-preprod/ent/customer"
	"iohk-golang-backend-preprod/ent/enttest"
	graphModel "iohk-golang-backend-preprod/graph/model" // Add this import
	"iohk-golang-backend-preprod/internal/domain/model"

	_ "github.com/mattn/go-sqlite3"
)

func TestCreate(t *testing.T) {
	testCases := []struct {
		name          string
		input         *model.Customer
		expectedName  string
		expectPanic   bool
		expectedError string
	}{
		{
			name: "Successful creation",
			input: &model.Customer{
				Name:       "Test",
				Surname:    "User",
				Number:     12345,
				Gender:     model.GenderMale,
				Country:    "Testland",
				Dependants: 0,
				BirthDate:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedName: "Test",
			expectPanic:  false,
		},
		{
			name:          "Nil customer input",
			input:         nil,
			expectPanic:   true,
			expectedError: "runtime error: invalid memory address or nil pointer dereference",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer client.Close()
			repo := NewCustomerRepository(client)

			// Act & Assert
			if tc.expectPanic {
				assert.Panics(t, func() {
					repo.Create(context.Background(), tc.input)
				}, "Expected panic for nil input")
			} else {
				createdCustomer, err := repo.Create(context.Background(), tc.input)
				assert.NoError(t, err)
				assert.NotNil(t, createdCustomer)
				assert.NotZero(t, createdCustomer.ID)
				assert.Equal(t, tc.expectedName, createdCustomer.Name)
			}
		})
	}
}

func TestGetByID(t *testing.T) {
	testCases := []struct {
		name          string
		setupFunc     func(*ent.Client) string
		expectedName  string
		expectedError string
	}{
		{
			name: "Existing customer",
			setupFunc: func(client *ent.Client) string {
				customer, err := client.Customer.Create().
					SetName("Existing").
					SetSurname("User").
					SetNumber(12345).
					SetGender(customer.GenderMale).
					SetCountry("TestCountry").
					SetDependants(0).
					SetBirthDate(time.Now()).
					Save(context.Background())
				if err != nil {
					t.Fatalf("Failed to create test customer: %v", err)
				}
				return strconv.Itoa(customer.ID)
			},
			expectedName:  "Existing",
			expectedError: "",
		},
		{
			name: "Non-existing customer",
			setupFunc: func(client *ent.Client) string {
				return "non-existing-id"
			},
			expectedName:  "",
			expectedError: "strconv.Atoi: parsing \"non-existing-id\": invalid syntax",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer client.Close()
			repo := NewCustomerRepository(client)
			id := tc.setupFunc(client)

			// Act
			customer, err := repo.GetByID(context.Background(), id)

			// Assert
			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
				assert.Nil(t, customer)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, customer)
				assert.Equal(t, tc.expectedName, customer.Name)
			}
		})
	}
}

func TestGetAll(t *testing.T) {
	testCases := []struct {
		name          string
		setupFunc     func(*ent.Client) error
		expectedCount int
		expectedError string
	}{
		{
			name: "Multiple customers",
			setupFunc: func(client *ent.Client) error {
				_, err := client.Customer.Create().
					SetName("User1").
					SetSurname("Surname1").
					SetNumber(12345).
					SetGender(customer.GenderMale).
					SetCountry("Country1").
					SetDependants(0).
					SetBirthDate(time.Now()).
					Save(context.Background())
				if err != nil {
					return err
				}
				_, err = client.Customer.Create().
					SetName("User2").
					SetSurname("Surname2").
					SetNumber(67890).
					SetGender(customer.GenderFemale).
					SetCountry("Country2").
					SetDependants(1).
					SetBirthDate(time.Now()).
					Save(context.Background())
				return err
			},
			expectedCount: 2,
			expectedError: "",
		},
		{
			name:          "No customers",
			setupFunc:     func(client *ent.Client) error { return nil },
			expectedCount: 0,
			expectedError: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer client.Close()
			repo := NewCustomerRepository(client)
			err := tc.setupFunc(client)
			assert.NoError(t, err, "Setup should not fail")

			// Act
			customers, err := repo.GetAll(context.Background())

			// Assert
			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
				assert.Nil(t, customers)
			} else {
				assert.NoError(t, err)
				assert.Len(t, customers, tc.expectedCount)
				if tc.expectedCount > 0 {
					assert.NotEmpty(t, customers[0].Name)
					assert.NotEmpty(t, customers[0].Surname)
					assert.NotZero(t, customers[0].Number)
				}
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	testCases := []struct {
		name          string
		setupFunc     func(*ent.Client) string
		updateFunc    func() *graphModel.UpdateCustomerInput
		expectedName  string
		expectedError string
	}{
		{
			name: "Successful update",
			setupFunc: func(client *ent.Client) string {
				customer, err := client.Customer.Create().
					SetName("Original").
					SetSurname("User").
					SetNumber(12345).
					SetGender(customer.GenderMale).
					SetCountry("TestCountry").
					SetDependants(0).
					SetBirthDate(time.Now()).
					Save(context.Background())
				if err != nil {
					t.Fatalf("Failed to create test customer: %v", err)
				}
				return strconv.Itoa(customer.ID)
			},
			updateFunc: func() *graphModel.UpdateCustomerInput {
				return &graphModel.UpdateCustomerInput{
					Name: stringPtr("Updated"),
				}
			},
			expectedName:  "Updated",
			expectedError: "",
		},
		{
			name: "Non-existing customer",
			setupFunc: func(client *ent.Client) string {
				return "non-existing-id"
			},
			updateFunc: func() *graphModel.UpdateCustomerInput {
				return &graphModel.UpdateCustomerInput{
					Name: stringPtr("Updated"),
				}
			},
			expectedName:  "",
			expectedError: "strconv.Atoi: parsing \"non-existing-id\": invalid syntax",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer client.Close()
			repo := NewCustomerRepository(client)
			id := tc.setupFunc(client)

			// Act
			updatedCustomer, err := repo.Update(context.Background(), id, tc.updateFunc())

			// Assert
			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
				assert.Nil(t, updatedCustomer)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, updatedCustomer)
				assert.Equal(t, tc.expectedName, updatedCustomer.Name)
			}
		})
	}
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}

func TestDelete(t *testing.T) {
	testCases := []struct {
		name          string
		setupFunc     func(*ent.Client) string
		expectedError string
	}{
		{
			name: "Successful deletion",
			setupFunc: func(client *ent.Client) string {
				customer, err := client.Customer.Create().
					SetName("ToDelete").
					SetSurname("User").
					SetNumber(12345).
					SetGender(customer.GenderMale).
					SetCountry("TestCountry").
					SetDependants(0).
					SetBirthDate(time.Now()).
					Save(context.Background())
				if err != nil {
					t.Fatalf("Failed to create test customer: %v", err)
				}
				return strconv.Itoa(customer.ID)
			},
			expectedError: "",
		},
		{
			name: "Non-existing customer",
			setupFunc: func(client *ent.Client) string {
				return "non-existing-id"
			},
			expectedError: "strconv.Atoi: parsing \"non-existing-id\": invalid syntax",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer client.Close()
			repo := NewCustomerRepository(client)
			id := tc.setupFunc(client)

			// Act
			err := repo.Delete(context.Background(), id)

			// Assert
			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
			} else {
				assert.NoError(t, err)
				// Verify deletion
				_, err := repo.GetByID(context.Background(), id)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "not found")
			}
		})
	}
}
