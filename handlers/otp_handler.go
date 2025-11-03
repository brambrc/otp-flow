package handlers

import (
	"net/http"
	"prenup/models"

	"github.com/gin-gonic/gin"
)

type OTPService interface {
	RequestOTP(userID string) (*models.OTP, error)
	ValidateOTP(userID, otpCode string) (bool, error)
}

type OTPHandler struct {
	service OTPService
}

func NewOTPHandler(service OTPService) *OTPHandler {
	return &OTPHandler{service: service}
}

func (h *OTPHandler) RequestOTP(c *gin.Context) {
	var req models.CreateOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:            "invalid_request",
			ErrorDescription: "Missing required field: user_id",
		})
		return
	}

	otp, err := h.service.RequestOTP(req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:            "internal_error",
			ErrorDescription: "Failed to generate OTP",
		})
		return
	}

	c.JSON(http.StatusOK, models.CreateOTPResponse{
		UserID: otp.UserID,
		OTP:    otp.Code,
	})
}

func (h *OTPHandler) ValidateOTP(c *gin.Context) {
	var req models.ValidateOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:            "invalid_request",
			ErrorDescription: "Missing required fields: user_id, otp",
		})
		return
	}

	valid, err := h.service.ValidateOTP(req.UserID, req.OTP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:            "internal_error",
			ErrorDescription: "Failed to validate OTP",
		})
		return
	}

	if !valid {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:            "otp_not_found",
			ErrorDescription: "OTP Not Found",
		})
		return
	}

	c.JSON(http.StatusOK, models.ValidateOTPResponse{
		UserID:  req.UserID,
		Message: "OTP validated successfully.",
	})
}
