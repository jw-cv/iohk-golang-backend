package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	graphModel "iohk-golang-backend-preprod/graph/model"
	"iohk-golang-backend-preprod/internal/domain/model"
)

type MockCustomerRepository struct {
	mock.Mock
}

func (m *MockCustomerRepository) Create(ctx context.Context, customer *model.Customer) (*model.Customer, error) {
	args := m.Called(ctx, customer)
	return args.Get(0).(*model.Customer), args.Error(1)
}

func (m *MockCustomerRepository) GetByID(ctx context.Context, id string) (*model.Customer, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Customer), args.Error(1)
}

func (m *MockCustomerRepository) GetAll(ctx context.Context) ([]*model.Customer, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*model.Customer), args.Error(1)
}

func (m *MockCustomerRepository) Update(ctx context.Context, id string, input *graphModel.UpdateCustomerInput) (*model.Customer, error) {
	args := m.Called(ctx, id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Customer), args.Error(1)
}

func (m *MockCustomerRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateCustomer(t *testing.T) {
	testCases := []struct {
		name          string
		input         *model.Customer
		mockBehavior  func(m *MockCustomerRepository, input *model.Customer)
		expected      *model.Customer
		expectedError error
	}{
		{
			name: "Successful creation",
			input: &model.Customer{
				Name:       "Alice",
				Surname:    "Smith",
				Number:     123,
				Gender:     model.GenderFemale,
				Country:    "UK",
				Dependants: 0,
				BirthDate:  time.Now(),
			},
			mockBehavior: func(m *MockCustomerRepository, input *model.Customer) {
				m.On("Create", mock.Anything, input).Return(input, nil)
			},
			expected:      &model.Customer{Name: "Alice", Surname: "Smith"},
			expectedError: nil,
		},
		{
			name: "Repository error",
			input: &model.Customer{
				Name: "Bob",
			},
			mockBehavior: func(m *MockCustomerRepository, input *model.Customer) {
				m.On("Create", mock.Anything, input).Return((*model.Customer)(nil), errors.New("repository error"))
			},
			expected:      nil,
			expectedError: errors.New("repository error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockRepo := new(MockCustomerRepository)
			service := NewCustomerService(mockRepo)
			tc.mockBehavior(mockRepo, tc.input)

			// Act
			result, err := service.CreateCustomer(context.Background(), tc.input)

			// Assert
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected.Name, result.Name)
				assert.Equal(t, tc.expected.Surname, result.Surname)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetCustomer(t *testing.T) {
	testCases := []struct {
		name          string
		id            string
		mockBehavior  func(m *MockCustomerRepository)
		expected      *model.Customer
		expectedError error
	}{
		{
			name: "Customer exists",
			id:   "1",
			mockBehavior: func(m *MockCustomerRepository) {
				m.On("GetByID", mock.Anything, "1").Return(&model.Customer{ID: 1, Name: "Alice"}, nil)
			},
			expected:      &model.Customer{ID: 1, Name: "Alice"},
			expectedError: nil,
		},
		{
			name: "Customer not found",
			id:   "2",
			mockBehavior: func(m *MockCustomerRepository) {
				m.On("GetByID", mock.Anything, "2").Return((*model.Customer)(nil), errors.New("customer not found"))
			},
			expected:      nil,
			expectedError: errors.New("customer not found"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockRepo := new(MockCustomerRepository)
			service := NewCustomerService(mockRepo)
			tc.mockBehavior(mockRepo)

			// Act
			result, err := service.GetCustomer(context.Background(), tc.id)

			// Assert
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetAllCustomers(t *testing.T) {
	testCases := []struct {
		name          string
		mockBehavior  func(m *MockCustomerRepository)
		expected      []*model.Customer
		expectedError error
	}{
		{
			name: "Customers exist",
			mockBehavior: func(m *MockCustomerRepository) {
				m.On("GetAll", mock.Anything).Return([]*model.Customer{
					{ID: 1, Name: "Alice"},
					{ID: 2, Name: "Bob"},
				}, nil)
			},
			expected: []*model.Customer{
				{ID: 1, Name: "Alice"},
				{ID: 2, Name: "Bob"},
			},
			expectedError: nil,
		},
		{
			name: "No customers",
			mockBehavior: func(m *MockCustomerRepository) {
				m.On("GetAll", mock.Anything).Return([]*model.Customer{}, nil)
			},
			expected:      []*model.Customer{},
			expectedError: nil,
		},
		{
			name: "Repository error",
			mockBehavior: func(m *MockCustomerRepository) {
				m.On("GetAll", mock.Anything).Return([]*model.Customer(nil), errors.New("repository error"))
			},
			expected:      nil,
			expectedError: errors.New("repository error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockRepo := new(MockCustomerRepository)
			service := NewCustomerService(mockRepo)
			tc.mockBehavior(mockRepo)

			// Act
			result, err := service.GetAllCustomers(context.Background())

			// Assert
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateCustomer(t *testing.T) {
	testCases := []struct {
		name          string
		id            string
		input         *model.Customer
		mockBehavior  func(m *MockCustomerRepository, id string, input *model.Customer)
		expected      *model.Customer
		expectedError error
	}{
		{
			name: "Successful update",
			id:   "1",
			input: &model.Customer{
				Name:       "Alice Updated",
				Surname:    "Smith",
				Number:     123,
				Gender:     model.GenderFemale,
				Country:    "UK",
				Dependants: 0,
				BirthDate:  time.Date(2024, 9, 29, 0, 0, 0, 0, time.UTC),
			},
			mockBehavior: func(m *MockCustomerRepository, id string, input *model.Customer) {
				birthDateStr := input.BirthDate.Format("2006-01-02")
				expectedInput := &graphModel.UpdateCustomerInput{
					Name:       &input.Name,
					Surname:    &input.Surname,
					Number:     &input.Number,
					Gender:     (*graphModel.Gender)(&input.Gender),
					Country:    &input.Country,
					Dependants: &input.Dependants,
					BirthDate:  &birthDateStr,
				}
				m.On("Update", mock.Anything, id, expectedInput).Return(input, nil)
			},
			expected:      &model.Customer{ID: 1, Name: "Alice Updated", Surname: "Smith"},
			expectedError: nil,
		},
		{
			name: "Customer not found",
			id:   "2",
			input: &model.Customer{
				Name:       "Bob",
				Surname:    "",
				Number:     0,
				Gender:     model.Gender(""),
				Country:    "",
				Dependants: 0,
				BirthDate:  time.Time{},
			},
			mockBehavior: func(m *MockCustomerRepository, id string, input *model.Customer) {
				birthDateStr := input.BirthDate.Format("2006-01-02")
				expectedInput := &graphModel.UpdateCustomerInput{
					Name:       &input.Name,
					Surname:    &input.Surname,
					Number:     &input.Number,
					Gender:     (*graphModel.Gender)(&input.Gender),
					Country:    &input.Country,
					Dependants: &input.Dependants,
					BirthDate:  &birthDateStr,
				}
				m.On("Update", mock.Anything, id, expectedInput).Return((*model.Customer)(nil), errors.New("customer not found"))
			},
			expected:      nil,
			expectedError: errors.New("customer not found"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockRepo := new(MockCustomerRepository)
			service := NewCustomerService(mockRepo)
			tc.mockBehavior(mockRepo, tc.id, tc.input)

			// Act
			result, err := service.UpdateCustomer(context.Background(), tc.id, tc.input)

			// Assert
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected.Name, result.Name)
				assert.Equal(t, tc.expected.Surname, result.Surname)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteCustomer(t *testing.T) {
	testCases := []struct {
		name          string
		id            string
		mockBehavior  func(m *MockCustomerRepository)
		expected      bool
		expectedError error
	}{
		{
			name: "Successful deletion",
			id:   "1",
			mockBehavior: func(m *MockCustomerRepository) {
				m.On("Delete", mock.Anything, "1").Return(nil)
			},
			expected:      true,
			expectedError: nil,
		},
		{
			name: "Customer not found",
			id:   "2",
			mockBehavior: func(m *MockCustomerRepository) {
				m.On("Delete", mock.Anything, "2").Return(errors.New("customer not found"))
			},
			expected:      false,
			expectedError: errors.New("customer not found"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockRepo := new(MockCustomerRepository)
			service := NewCustomerService(mockRepo)
			tc.mockBehavior(mockRepo)

			// Act
			result, err := service.DeleteCustomer(context.Background(), tc.id)

			// Assert
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.False(t, result)
			} else {
				assert.NoError(t, err)
				assert.True(t, result)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
