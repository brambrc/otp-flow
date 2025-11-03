package models

import "time"

type OTPStatus string

const (
	StatusCreated   OTPStatus = "created"
	StatusValidated OTPStatus = "validated"
	StatusExpired   OTPStatus = "expired"
)

type OTP struct {
	ID        int       `json:"id"`
	UserID    string    `json:"user_id"`
	Code      string    `json:"otp"`
	Status    OTPStatus `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type CreateOTPRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

type ValidateOTPRequest struct {
	UserID string `json:"user_id" binding:"required"`
	OTP    string `json:"otp" binding:"required"`
}

type CreateOTPResponse struct {
	UserID string `json:"user_id"`
	OTP    string `json:"otp"`
}

type ValidateOTPResponse struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}
