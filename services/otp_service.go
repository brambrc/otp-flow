package services

import (
	"fmt"
	"math/rand"
	"prenup/models"
	"time"
)

type OTPRepository interface {
	CreateOTP(userID, code string, expiresAt time.Time) (*models.OTP, error)
	GetOTPByUserID(userID string) (*models.OTP, error)
	UpdateOTPStatus(id int, status models.OTPStatus) error
}

type OTPService struct {
	repo OTPRepository
}

func NewOTPService(repo OTPRepository) *OTPService {
	return &OTPService{repo: repo}
}

func (s *OTPService) GenerateOTP() string {
	code := rand.Intn(900000) + 100000
	return fmt.Sprintf("%d", code)
}

func (s *OTPService) RequestOTP(userID string) (*models.OTP, error) {
	code := s.GenerateOTP()
	expiresAt := time.Now().Add(2 * time.Minute)

	otp, err := s.repo.CreateOTP(userID, code, expiresAt)
	if err != nil {
		return nil, err
	}

	return otp, nil
}

func (s *OTPService) ValidateOTP(userID, otpCode string) (bool, error) {
	otp, err := s.repo.GetOTPByUserID(userID)
	if err != nil {
		return false, err
	}

	if otp == nil {
		return false, nil
	}

	if time.Now().After(otp.ExpiresAt) {
		_ = s.repo.UpdateOTPStatus(otp.ID, models.StatusExpired)
		return false, nil
	}

	if otp.Code != otpCode {
		return false, nil
	}

	err = s.repo.UpdateOTPStatus(otp.ID, models.StatusValidated)
	if err != nil {
		return false, err
	}

	return true, nil
}
