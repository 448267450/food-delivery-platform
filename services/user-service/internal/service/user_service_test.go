package service_test

import (
	"errors"
	"testing"

	"github.com/448267450/food-delivery-platform/services/user-service/internal/model"
	"github.com/448267450/food-delivery-platform/services/user-service/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// ==================== Mock Repository ====================

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id uint) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

// ==================== Register Tests ====================

func TestRegister_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := service.NewUserService(mockRepo, "test-secret")

	req := &service.RegisterRequest{
		Name:     "Ryan Ren",
		Email:    "ryan@test.com",
		Password: "123456",
		Phone:    "512-000-0001",
	}

	mockRepo.On("FindByEmail", req.Email).Return(nil, errors.New("not found"))
	mockRepo.On("Create", mock.AnythingOfType("*model.User")).Return(nil)

	resp, err := svc.Register(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, req.Email, resp.User.Email)
	assert.Equal(t, req.Name, resp.User.Name)
	assert.Equal(t, "customer", resp.User.Role)
	mockRepo.AssertExpectations(t)
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := service.NewUserService(mockRepo, "test-secret")

	existingUser := &model.User{ID: 1, Email: "ryan@test.com"}

	req := &service.RegisterRequest{
		Name:     "Ryan Ren",
		Email:    "ryan@test.com",
		Password: "123456",
		Phone:    "512-000-0001",
	}

	mockRepo.On("FindByEmail", req.Email).Return(existingUser, nil)

	resp, err := svc.Register(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "email already registered", err.Error())
	mockRepo.AssertNotCalled(t, "Create")
}

func TestRegister_DatabaseError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := service.NewUserService(mockRepo, "test-secret")

	req := &service.RegisterRequest{
		Name:     "Ryan Ren",
		Email:    "ryan@test.com",
		Password: "123456",
		Phone:    "512-000-0001",
	}

	mockRepo.On("FindByEmail", req.Email).Return(nil, errors.New("not found"))
	mockRepo.On("Create", mock.AnythingOfType("*model.User")).Return(errors.New("db error"))

	resp, err := svc.Register(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

// ==================== Login Tests ====================

func TestLogin_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := service.NewUserService(mockRepo, "test-secret")

	// Generate a real bcrypt hash for "123456"
	hashed, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)

	user := &model.User{
		ID:       1,
		Name:     "Ryan Ren",
		Email:    "ryan@test.com",
		Password: string(hashed),
		Role:     "customer",
	}

	// Login calls FindByEmail once — return the user with hashed password
	mockRepo.On("FindByEmail", "ryan@test.com").Return(user, nil)

	resp, err := svc.Login(&service.LoginRequest{
		Email:    "ryan@test.com",
		Password: "123456",
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, "ryan@test.com", resp.User.Email)
}

func TestLogin_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := service.NewUserService(mockRepo, "test-secret")

	mockRepo.On("FindByEmail", "nobody@test.com").Return(nil, errors.New("not found"))

	resp, err := svc.Login(&service.LoginRequest{
		Email:    "nobody@test.com",
		Password: "123456",
	})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "invalid email or password", err.Error())
}

func TestLogin_WrongPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := service.NewUserService(mockRepo, "test-secret")

	hashed, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)

	user := &model.User{
		ID:       1,
		Email:    "ryan@test.com",
		Password: string(hashed),
	}

	mockRepo.On("FindByEmail", "ryan@test.com").Return(user, nil)

	resp, err := svc.Login(&service.LoginRequest{
		Email:    "ryan@test.com",
		Password: "wrongpassword",
	})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "invalid email or password", err.Error())
}

// ==================== GetProfile Tests ====================

func TestGetProfile_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := service.NewUserService(mockRepo, "test-secret")

	user := &model.User{
		ID:    1,
		Name:  "Ryan Ren",
		Email: "ryan@test.com",
		Role:  "customer",
	}

	mockRepo.On("FindByID", uint(1)).Return(user, nil)

	result, err := svc.GetProfile(1)

	assert.NoError(t, err)
	assert.Equal(t, user.Name, result.Name)
	assert.Equal(t, user.Email, result.Email)
}

func TestGetProfile_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := service.NewUserService(mockRepo, "test-secret")

	mockRepo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	result, err := svc.GetProfile(999)

	assert.Error(t, err)
	assert.Nil(t, result)
}
