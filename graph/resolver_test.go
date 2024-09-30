//go:build testcoverage
// +build testcoverage

package graph

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"iohk-golang-backend/graph/model"
	internalModel "iohk-golang-backend/internal/domain/model"
)

type MockCustomerService struct {
	mock.Mock
}

func (m *MockCustomerService) GetCustomer(ctx context.Context, id string) (*internalModel.Customer, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*internalModel.Customer), args.Error(1)
}

func (m *MockCustomerService) GetAllCustomers(ctx context.Context) ([]*internalModel.Customer, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*internalModel.Customer), args.Error(1)
}

func (m *MockCustomerService) CreateCustomer(ctx context.Context, customer *internalModel.Customer) (*internalModel.Customer, error) {
	args := m.Called(ctx, customer)
	return args.Get(0).(*internalModel.Customer), args.Error(1)
}

func (m *MockCustomerService) UpdateCustomer(ctx context.Context, id string, customer *internalModel.Customer) (*internalModel.Customer, error) {
	args := m.Called(ctx, id, customer)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*internalModel.Customer), args.Error(1)
}

func (m *MockCustomerService) DeleteCustomer(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func TestCustomer(t *testing.T) {
	testCases := []struct {
		name          string
		id            string
		mockBehavior  func(m *MockCustomerService)
		expected      *model.Customer
		expectedError error
	}{
		{
			name: "Customer exists",
			id:   "1",
			mockBehavior: func(m *MockCustomerService) {
				m.On("GetCustomer", mock.Anything, "1").Return(&internalModel.Customer{ID: 1, Name: "Alice", BirthDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)}, nil)
			},
			expected: &model.Customer{
				ID:        "1",
				Name:      "Alice",
				BirthDate: "2000-01-01",
			},
			expectedError: nil,
		},
		{
			name: "Customer not found",
			id:   "2",
			mockBehavior: func(m *MockCustomerService) {
				m.On("GetCustomer", mock.Anything, "2").Return((*internalModel.Customer)(nil), errors.New("customer not found"))
			},
			expected:      nil,
			expectedError: errors.New("customer not found"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockService := new(MockCustomerService)
			resolver := &Resolver{customerService: mockService}
			tc.mockBehavior(mockService)

			// Act
			result, err := resolver.Query().Customer(context.Background(), tc.id)

			// Assert
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected.ID, result.ID)
				assert.Equal(t, tc.expected.Name, result.Name)
				assert.Equal(t, tc.expected.BirthDate, result.BirthDate)
			}
			mockService.AssertExpectations(t)
		})
	}
}

func TestCustomers(t *testing.T) {
	testCases := []struct {
		name          string
		mockBehavior  func(m *MockCustomerService)
		expected      []*model.Customer
		expectedError error
	}{
		{
			name: "Customers exist",
			mockBehavior: func(m *MockCustomerService) {
				m.On("GetAllCustomers", mock.Anything).Return([]*internalModel.Customer{
					{ID: 1, Name: "Alice", BirthDate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)},
					{ID: 2, Name: "Bob", BirthDate: time.Date(1995, 2, 15, 0, 0, 0, 0, time.UTC)},
				}, nil)
			},
			expected: []*model.Customer{
				{ID: "1", Name: "Alice", BirthDate: "1990-01-01"},
				{ID: "2", Name: "Bob", BirthDate: "1995-02-15"},
			},
			expectedError: nil,
		},
		{
			name: "No customers",
			mockBehavior: func(m *MockCustomerService) {
				m.On("GetAllCustomers", mock.Anything).Return([]*internalModel.Customer{}, nil)
			},
			expected:      []*model.Customer{},
			expectedError: nil,
		},
		{
			name: "Service error",
			mockBehavior: func(m *MockCustomerService) {
				m.On("GetAllCustomers", mock.Anything).Return([]*internalModel.Customer(nil), errors.New("service error"))
			},
			expected:      nil,
			expectedError: errors.New("service error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockService := new(MockCustomerService)
			resolver := &Resolver{customerService: mockService}
			tc.mockBehavior(mockService)

			// Act
			result, err := resolver.Query().Customers(context.Background())

			// Assert
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tc.expected), len(result))
				for i, expectedCustomer := range tc.expected {
					assert.Equal(t, expectedCustomer.ID, result[i].ID)
					assert.Equal(t, expectedCustomer.Name, result[i].Name)
					assert.Equal(t, expectedCustomer.BirthDate, result[i].BirthDate)
				}
			}
			mockService.AssertExpectations(t)
		})
	}
}

func TestCreateCustomer(t *testing.T) {
	testCases := []struct {
		name          string
		input         model.CreateCustomerInput
		mockBehavior  func(m *MockCustomerService)
		expected      *model.Customer
		expectedError error
	}{
		{
			name: "Successful creation",
			input: model.CreateCustomerInput{
				Name:      "Alice",
				BirthDate: "1990-01-01",
			},
			mockBehavior: func(m *MockCustomerService) {
				m.On("CreateCustomer", mock.Anything, mock.AnythingOfType("*model.Customer")).Return(&internalModel.Customer{
					ID:        1,
					Name:      "Alice",
					BirthDate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				}, nil)
			},
			expected: &model.Customer{
				ID:        "1",
				Name:      "Alice",
				BirthDate: "1990-01-01",
			},
			expectedError: nil,
		},
		{
			name: "Service error",
			input: model.CreateCustomerInput{
				Name: "Bob",
			},
			mockBehavior: func(m *MockCustomerService) {
				m.On("CreateCustomer", mock.Anything, mock.AnythingOfType("*model.Customer")).Return((*internalModel.Customer)(nil), errors.New("service error"))
			},
			expected:      nil,
			expectedError: errors.New("service error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockService := new(MockCustomerService)
			resolver := &Resolver{customerService: mockService}
			tc.mockBehavior(mockService)

			// Act
			result, err := resolver.Mutation().CreateCustomer(context.Background(), tc.input)

			// Assert
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected.ID, result.ID)
				assert.Equal(t, tc.expected.Name, result.Name)
				assert.Equal(t, tc.expected.BirthDate, result.BirthDate)
			}
			mockService.AssertExpectations(t)
		})
	}
}

func TestUpdateCustomer(t *testing.T) {
	testCases := []struct {
		name          string
		id            string
		input         model.UpdateCustomerInput
		mockBehavior  func(m *MockCustomerService)
		expected      *model.Customer
		expectedError error
	}{
		{
			name: "Successful update",
			id:   "1",
			input: model.UpdateCustomerInput{
				Name:      stringPtr("Alice Updated"),
				BirthDate: stringPtr("1990-01-01"),
			},
			mockBehavior: func(m *MockCustomerService) {
				m.On("UpdateCustomer", mock.Anything, "1", mock.AnythingOfType("*model.Customer")).Return(&internalModel.Customer{
					ID:        1,
					Name:      "Alice Updated",
					BirthDate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				}, nil)
			},
			expected: &model.Customer{
				ID:        "1",
				Name:      "Alice Updated",
				BirthDate: "1990-01-01",
			},
			expectedError: nil,
		},
		{
			name: "Customer not found",
			id:   "2",
			input: model.UpdateCustomerInput{
				Name: stringPtr("Bob"),
			},
			mockBehavior: func(m *MockCustomerService) {
				m.On("UpdateCustomer", mock.Anything, "2", mock.AnythingOfType("*model.Customer")).Return((*internalModel.Customer)(nil), errors.New("customer not found"))
			},
			expected:      nil,
			expectedError: errors.New("customer not found"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockService := new(MockCustomerService)
			resolver := &Resolver{customerService: mockService}
			tc.mockBehavior(mockService)

			// Act
			result, err := resolver.Mutation().UpdateCustomer(context.Background(), tc.id, tc.input)

			// Assert
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected.ID, result.ID)
				assert.Equal(t, tc.expected.Name, result.Name)
				assert.Equal(t, tc.expected.BirthDate, result.BirthDate)
			}
			mockService.AssertExpectations(t)
		})
	}
}

func TestDeleteCustomer(t *testing.T) {
	testCases := []struct {
		name          string
		id            string
		mockBehavior  func(m *MockCustomerService)
		expected      bool
		expectedError error
	}{
		{
			name: "Successful deletion",
			id:   "1",
			mockBehavior: func(m *MockCustomerService) {
				m.On("DeleteCustomer", mock.Anything, "1").Return(true, nil)
			},
			expected:      true,
			expectedError: nil,
		},
		{
			name: "Customer not found",
			id:   "2",
			mockBehavior: func(m *MockCustomerService) {
				m.On("DeleteCustomer", mock.Anything, "2").Return(false, errors.New("customer not found"))
			},
			expected:      false,
			expectedError: errors.New("customer not found"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockService := new(MockCustomerService)
			resolver := &Resolver{customerService: mockService}
			tc.mockBehavior(mockService)

			// Act
			result, err := resolver.Mutation().DeleteCustomer(context.Background(), tc.id)

			// Assert
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.False(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
			mockService.AssertExpectations(t)
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
