package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"prenup/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOTPService struct {
	mock.Mock
}

func (m *MockOTPService) RequestOTP(userID string) (*models.OTP, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.OTP), args.Error(1)
}

func (m *MockOTPService) ValidateOTP(userID, otpCode string) (bool, error) {
	args := m.Called(userID, otpCode)
	return args.Bool(0), args.Error(1)
}

func TestRequestOTP_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockOTPService)

	userID := "Robert"
	otpCode := "123456"
	now := time.Now()

	mockService.On("RequestOTP", userID).
		Return(&models.OTP{
			ID:        1,
			UserID:    userID,
			Code:      otpCode,
			Status:    models.StatusCreated,
			CreatedAt: now,
			ExpiresAt: now.Add(2 * time.Minute),
		}, nil)

	handler := NewOTPHandler(mockService)

	body := models.CreateOTPRequest{UserID: userID}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/otp/request", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	handler.RequestOTP(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.CreateOTPResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, userID, response.UserID)
	assert.Equal(t, otpCode, response.OTP)
	mockService.AssertExpectations(t)
}

func TestRequestOTP_MissingUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockOTPService)
	handler := NewOTPHandler(mockService)

	body := map[string]string{}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/otp/request", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	handler.RequestOTP(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestValidateOTP_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockOTPService)

	userID := "Robert"
	otpCode := "123456"

	mockService.On("ValidateOTP", userID, otpCode).Return(true, nil)

	handler := NewOTPHandler(mockService)

	body := models.ValidateOTPRequest{UserID: userID, OTP: otpCode}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/otp/validate", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	handler.ValidateOTP(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.ValidateOTPResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, userID, response.UserID)
	assert.Equal(t, "OTP validated successfully.", response.Message)
	mockService.AssertExpectations(t)
}

func TestValidateOTP_InvalidOTP(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockOTPService)

	userID := "Robert"
	otpCode := "000000"

	mockService.On("ValidateOTP", userID, otpCode).Return(false, nil)

	handler := NewOTPHandler(mockService)

	body := models.ValidateOTPRequest{UserID: userID, OTP: otpCode}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/otp/validate", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	handler.ValidateOTP(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response models.ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "otp_not_found", response.Error)
	mockService.AssertExpectations(t)
}

func TestValidateOTP_MissingFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockOTPService)
	handler := NewOTPHandler(mockService)

	body := map[string]string{"user_id": "Robert"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/otp/validate", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	handler.ValidateOTP(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
