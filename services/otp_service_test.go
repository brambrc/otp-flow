package services

import (
	"testing"
	"time"

	"prenup/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOTPRepository struct {
	mock.Mock
}

func (m *MockOTPRepository) CreateOTP(userID, code string, expiresAt time.Time) (*models.OTP, error) {
	args := m.Called(userID, code, expiresAt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.OTP), args.Error(1)
}

func (m *MockOTPRepository) GetOTPByUserID(userID string) (*models.OTP, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.OTP), args.Error(1)
}

func (m *MockOTPRepository) UpdateOTPStatus(id int, status models.OTPStatus) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func TestGenerateOTP(t *testing.T) {
	mockRepo := new(MockOTPRepository)
	service := NewOTPService(mockRepo)

	code := service.GenerateOTP()

	assert.NotEmpty(t, code)
	assert.Len(t, code, 6)
	assert.GreaterOrEqual(t, code, "100000")
	assert.LessOrEqual(t, code, "999999")
}

func TestRequestOTP(t *testing.T) {
	mockRepo := new(MockOTPRepository)
	service := NewOTPService(mockRepo)

	userID := "Robert"
	now := time.Now()

	mockRepo.On("CreateOTP", userID, mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).
		Return(&models.OTP{
			ID:        1,
			UserID:    userID,
			Code:      "123456",
			Status:    models.StatusCreated,
			CreatedAt: now,
			ExpiresAt: now.Add(2 * time.Minute),
		}, nil)

	otp, err := service.RequestOTP(userID)

	assert.NoError(t, err)
	assert.NotNil(t, otp)
	assert.Equal(t, userID, otp.UserID)
	assert.Equal(t, models.StatusCreated, otp.Status)
	mockRepo.AssertExpectations(t)
}

func TestValidateOTP_Success(t *testing.T) {
	mockRepo := new(MockOTPRepository)
	service := NewOTPService(mockRepo)

	userID := "Robert"
	code := "123456"
	now := time.Now()

	mockRepo.On("GetOTPByUserID", userID).
		Return(&models.OTP{
			ID:        1,
			UserID:    userID,
			Code:      code,
			Status:    models.StatusCreated,
			CreatedAt: now,
			ExpiresAt: now.Add(2 * time.Minute),
		}, nil)

	mockRepo.On("UpdateOTPStatus", 1, models.StatusValidated).Return(nil)

	valid, err := service.ValidateOTP(userID, code)

	assert.NoError(t, err)
	assert.True(t, valid)
	mockRepo.AssertExpectations(t)
}

func TestValidateOTP_InvalidCode(t *testing.T) {
	mockRepo := new(MockOTPRepository)
	service := NewOTPService(mockRepo)

	userID := "Robert"
	now := time.Now()

	mockRepo.On("GetOTPByUserID", userID).
		Return(&models.OTP{
			ID:        1,
			UserID:    userID,
			Code:      "123456",
			Status:    models.StatusCreated,
			CreatedAt: now,
			ExpiresAt: now.Add(2 * time.Minute),
		}, nil)

	valid, err := service.ValidateOTP(userID, "000000")

	assert.NoError(t, err)
	assert.False(t, valid)
	mockRepo.AssertExpectations(t)
}

func TestValidateOTP_Expired(t *testing.T) {
	mockRepo := new(MockOTPRepository)
	service := NewOTPService(mockRepo)

	userID := "Robert"
	now := time.Now()

	mockRepo.On("GetOTPByUserID", userID).
		Return(&models.OTP{
			ID:        1,
			UserID:    userID,
			Code:      "123456",
			Status:    models.StatusCreated,
			CreatedAt: now.Add(-3 * time.Minute),
			ExpiresAt: now.Add(-1 * time.Minute),
		}, nil)

	mockRepo.On("UpdateOTPStatus", 1, models.StatusExpired).Return(nil)

	valid, err := service.ValidateOTP(userID, "123456")

	assert.NoError(t, err)
	assert.False(t, valid)
	mockRepo.AssertExpectations(t)
}

func TestValidateOTP_NotFound(t *testing.T) {
	mockRepo := new(MockOTPRepository)
	service := NewOTPService(mockRepo)

	userID := "Unknown"

	mockRepo.On("GetOTPByUserID", userID).Return(nil, nil)

	valid, err := service.ValidateOTP(userID, "123456")

	assert.NoError(t, err)
	assert.False(t, valid)
	mockRepo.AssertExpectations(t)
}
