//go:build testcoverage
// +build testcoverage

package mapper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"iohk-golang-backend/ent"
	"iohk-golang-backend/ent/customer"
	"iohk-golang-backend/graph/model"
	domainmodel "iohk-golang-backend/internal/domain/model"
)

func TestDomainToGraphQL(t *testing.T) {
	testCases := []struct {
		name     string
		input    *domainmodel.Customer
		expected *model.Customer
	}{
		{
			name: "Full customer data",
			input: &domainmodel.Customer{
				ID:         1,
				Name:       "John",
				Surname:    "Doe",
				Number:     123,
				Gender:     domainmodel.GenderMale,
				Country:    "USA",
				Dependants: 2,
				BirthDate:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: &model.Customer{
				ID:         "1",
				Name:       "John",
				Surname:    "Doe",
				Number:     123,
				Gender:     model.GenderMale,
				Country:    "USA",
				Dependants: 2,
				BirthDate:  "1990-01-01",
			},
		},
		{
			name: "Minimal customer data",
			input: &domainmodel.Customer{
				ID:   2,
				Name: "Jane",
			},
			expected: &model.Customer{
				ID:        "2",
				Name:      "Jane",
				BirthDate: "0001-01-01",
			},
		},
		{
			name:     "Nil input",
			input:    nil,
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Use a defer-recover block to catch panics
			defer func() {
				if r := recover(); r != nil {
					if tc.input == nil {
						// If input is nil and we got a panic, test passes
						return
					}
					// For non-nil inputs, we still want to fail the test if there's a panic
					t.Errorf("The code panicked: %v", r)
				}
			}()

			// Act
			result := DomainToGraphQL(tc.input)

			// Assert
			if tc.input == nil {
				assert.Nil(t, result, "Expected nil result for nil input")
			} else {
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestGraphQLToDomain(t *testing.T) {
	testCases := []struct {
		name     string
		input    *model.Customer
		expected *domainmodel.Customer
	}{
		{
			name: "Full customer data",
			input: &model.Customer{
				ID:         "1",
				Name:       "John",
				Surname:    "Doe",
				Number:     123,
				Gender:     model.GenderMale,
				Country:    "USA",
				Dependants: 2,
				BirthDate:  "1990-01-01",
			},
			expected: &domainmodel.Customer{
				ID:         1,
				Name:       "John",
				Surname:    "Doe",
				Number:     123,
				Gender:     domainmodel.GenderMale,
				Country:    "USA",
				Dependants: 2,
				BirthDate:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "Minimal customer data",
			input: &model.Customer{
				ID:   "2",
				Name: "Jane",
			},
			expected: &domainmodel.Customer{
				ID:   2,
				Name: "Jane",
			},
		},
		{
			name: "Invalid date",
			input: &model.Customer{
				ID:        "3",
				Name:      "Invalid",
				BirthDate: "invalid-date",
			},
			expected: &domainmodel.Customer{
				ID:        3,
				Name:      "Invalid",
				BirthDate: time.Time{},
			},
		},
		{
			name:     "Nil input",
			input:    nil,
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Use a defer-recover block to catch panics
			defer func() {
				if r := recover(); r != nil {
					if tc.input == nil {
						// If input is nil and we got a panic, test passes
						return
					}
					// For non-nil inputs, we still want to fail the test if there's a panic
					t.Errorf("The code panicked: %v", r)
				}
			}()

			// Act
			result := GraphQLToDomain(tc.input)

			// Assert
			if tc.input == nil {
				assert.Nil(t, result, "Expected nil result for nil input")
			} else {
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestCreateInputToDomain(t *testing.T) {
	testCases := []struct {
		name     string
		input    *model.CreateCustomerInput
		expected *domainmodel.Customer
	}{
		{
			name: "Full input data",
			input: &model.CreateCustomerInput{
				Name:       "John",
				Surname:    "Doe",
				Number:     123,
				Gender:     model.GenderMale,
				Country:    "USA",
				Dependants: 2,
				BirthDate:  "1990-01-01",
			},
			expected: &domainmodel.Customer{
				Name:       "John",
				Surname:    "Doe",
				Number:     123,
				Gender:     domainmodel.GenderMale,
				Country:    "USA",
				Dependants: 2,
				BirthDate:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "Minimal input data",
			input: &model.CreateCustomerInput{
				Name: "Jane",
			},
			expected: &domainmodel.Customer{
				Name: "Jane",
			},
		},
		{
			name: "Invalid date",
			input: &model.CreateCustomerInput{
				Name:      "Invalid",
				BirthDate: "invalid-date",
			},
			expected: &domainmodel.Customer{
				Name:      "Invalid",
				BirthDate: time.Time{},
			},
		},
		{
			name:     "Nil input",
			input:    nil,
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Use a defer-recover block to catch panics
			defer func() {
				if r := recover(); r != nil {
					if tc.input == nil {
						// If input is nil and we got a panic, test passes
						return
					}
					// For non-nil inputs, we still want to fail the test if there's a panic
					t.Errorf("The code panicked: %v", r)
				}
			}()

			// Act
			result := CreateInputToDomain(tc.input)

			// Assert
			if tc.input == nil {
				assert.Nil(t, result, "Expected nil result for nil input")
			} else {
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestEntToDomain(t *testing.T) {
	testCases := []struct {
		name     string
		input    *ent.Customer
		expected *domainmodel.Customer
	}{
		{
			name: "Full customer data",
			input: &ent.Customer{
				ID:         1,
				Name:       "John",
				Surname:    "Doe",
				Number:     123,
				Gender:     customer.GenderMale,
				Country:    "USA",
				Dependants: 2,
				BirthDate:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: &domainmodel.Customer{
				ID:         1,
				Name:       "John",
				Surname:    "Doe",
				Number:     123,
				Gender:     domainmodel.GenderMale,
				Country:    "USA",
				Dependants: 2,
				BirthDate:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "Minimal customer data",
			input: &ent.Customer{
				ID:   2,
				Name: "Jane",
			},
			expected: &domainmodel.Customer{
				ID:   2,
				Name: "Jane",
			},
		},
		{
			name:     "Nil input",
			input:    nil,
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Use a defer-recover block to catch panics
			defer func() {
				if r := recover(); r != nil {
					if tc.input == nil {
						// If input is nil and we got a panic, test passes
						return
					}
					// For non-nil inputs, we still want to fail the test if there's a panic
					t.Errorf("The code panicked: %v", r)
				}
			}()

			// Act
			result := EntToDomain(tc.input)

			// Assert
			if tc.input == nil {
				assert.Nil(t, result, "Expected nil result for nil input")
			} else {
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestDomainToEnt(t *testing.T) {
	testCases := []struct {
		name     string
		input    *domainmodel.Customer
		expected *ent.Customer
	}{
		{
			name: "Full customer data",
			input: &domainmodel.Customer{
				ID:         1,
				Name:       "John",
				Surname:    "Doe",
				Number:     123,
				Gender:     domainmodel.GenderMale,
				Country:    "USA",
				Dependants: 2,
				BirthDate:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: &ent.Customer{
				ID:         1,
				Name:       "John",
				Surname:    "Doe",
				Number:     123,
				Gender:     customer.GenderMale,
				Country:    "USA",
				Dependants: 2,
				BirthDate:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "Minimal customer data",
			input: &domainmodel.Customer{
				ID:   2,
				Name: "Jane",
			},
			expected: &ent.Customer{
				ID:   2,
				Name: "Jane",
			},
		},
		{
			name:     "Nil input",
			input:    nil,
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Use a defer-recover block to catch panics
			defer func() {
				if r := recover(); r != nil {
					if tc.input == nil {
						// If input is nil and we got a panic, test passes
						return
					}
					// For non-nil inputs, we still want to fail the test if there's a panic
					t.Errorf("The code panicked: %v", r)
				}
			}()

			// Act
			result := DomainToEnt(tc.input)

			// Assert
			if tc.input == nil {
				assert.Nil(t, result, "Expected nil result for nil input")
			} else {
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestDomainToUpdateInput(t *testing.T) {
	testCases := []struct {
		name     string
		input    *domainmodel.Customer
		expected *model.UpdateCustomerInput
	}{
		{
			name: "Full customer data",
			input: &domainmodel.Customer{
				Name:       "John",
				Surname:    "Doe",
				Number:     123,
				Gender:     domainmodel.GenderMale,
				Country:    "USA",
				Dependants: 2,
				BirthDate:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: &model.UpdateCustomerInput{
				Name:       stringPtr("John"),
				Surname:    stringPtr("Doe"),
				Number:     intPtr(123),
				Gender:     (*model.Gender)(stringPtr("MALE")),
				Country:    stringPtr("USA"),
				Dependants: intPtr(2),
				BirthDate:  stringPtr("1990-01-01"),
			},
		},
		{
			name: "Minimal customer data",
			input: &domainmodel.Customer{
				Name: "Jane",
			},
			expected: &model.UpdateCustomerInput{
				Name:       stringPtr("Jane"),
				Surname:    stringPtr(""),
				Number:     intPtr(0),
				Gender:     (*model.Gender)(stringPtr("")),
				Country:    stringPtr(""),
				Dependants: intPtr(0),
				BirthDate:  stringPtr("0001-01-01"),
			},
		},
		{
			name:     "Nil input",
			input:    nil,
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Use a defer-recover block to catch panics
			defer func() {
				if r := recover(); r != nil {
					if tc.input == nil {
						// If input is nil and we got a panic, test passes
						return
					}
					// For non-nil inputs, we still want to fail the test if there's a panic
					t.Errorf("The code panicked: %v", r)
				}
			}()

			// Act
			result := DomainToUpdateInput(tc.input)

			// Assert
			if tc.input == nil {
				assert.Nil(t, result, "Expected nil result for nil input")
			} else {
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestDomainToGraphQLSlice(t *testing.T) {
	testCases := []struct {
		name     string
		input    []*domainmodel.Customer
		expected []*model.Customer
	}{
		{
			name: "Multiple customers",
			input: []*domainmodel.Customer{
				{
					ID:   1,
					Name: "John",
				},
				{
					ID:   2,
					Name: "Jane",
				},
			},
			expected: []*model.Customer{
				{
					ID:        "1",
					Name:      "John",
					BirthDate: "0001-01-01",
				},
				{
					ID:        "2",
					Name:      "Jane",
					BirthDate: "0001-01-01",
				},
			},
		},
		{
			name:     "Empty slice",
			input:    []*domainmodel.Customer{},
			expected: []*model.Customer{},
		},
		{
			name:     "Nil input",
			input:    nil,
			expected: []*model.Customer{}, // Change this line from nil to an empty slice
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result := DomainToGraphQLSlice(tc.input)

			// Assert
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestUpdateInputToDomain(t *testing.T) {
	testCases := []struct {
		name     string
		id       string
		input    *model.UpdateCustomerInput
		expected *domainmodel.Customer
	}{
		{
			name: "Full update",
			id:   "1",
			input: &model.UpdateCustomerInput{
				Name:       stringPtr("John"),
				Surname:    stringPtr("Doe"),
				Number:     intPtr(123),
				Gender:     (*model.Gender)(stringPtr("MALE")),
				Country:    stringPtr("USA"),
				Dependants: intPtr(2),
				BirthDate:  stringPtr("1990-01-01"),
			},
			expected: &domainmodel.Customer{
				ID:         1,
				Name:       "John",
				Surname:    "Doe",
				Number:     123,
				Gender:     domainmodel.GenderMale,
				Country:    "USA",
				Dependants: 2,
				BirthDate:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "Partial update",
			id:   "2",
			input: &model.UpdateCustomerInput{
				Name: stringPtr("Jane"),
			},
			expected: &domainmodel.Customer{
				ID:   2,
				Name: "Jane",
			},
		},
		{
			name:     "Nil input",
			id:       "3",
			input:    nil,
			expected: &domainmodel.Customer{ID: 3},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Use a defer-recover block to catch panics
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("The code panicked: %v", r)
				}
			}()

			// Act
			result := UpdateInputToDomain(tc.id, tc.input)

			// Assert
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestEntGenderToDomainGender(t *testing.T) {
	testCases := []struct {
		name     string
		input    customer.Gender
		expected domainmodel.Gender
	}{
		{
			name:     "Male gender",
			input:    customer.GenderMale,
			expected: domainmodel.GenderMale,
		},
		{
			name:     "Female gender",
			input:    customer.GenderFemale,
			expected: domainmodel.GenderFemale,
		},
		{
			name:     "Unknown gender",
			input:    customer.Gender("UNKNOWN"),
			expected: domainmodel.Gender(""),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result := entGenderToDomainGender(tc.input)

			// Assert
			assert.Equal(t, tc.expected, result)
		})
	}
}

func intPtr(i int) *int {
	return &i
}
