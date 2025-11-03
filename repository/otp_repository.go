package repository

import (
	"database/sql"
	"prenup/models"
	"time"
)

type OTPRepository struct {
	db *sql.DB
}

func NewOTPRepository(db *sql.DB) *OTPRepository {
	return &OTPRepository{db: db}
}

func (r *OTPRepository) CreateOTP(userID, code string, expiresAt time.Time) (*models.OTP, error) {
	query := `
	INSERT INTO otp (user_id, code, status, expires_at)
	VALUES ($1, $2, $3, $4)
	RETURNING id, user_id, code, status, created_at, expires_at
	`

	otp := &models.OTP{}
	err := r.db.QueryRow(query, userID, code, models.StatusCreated, expiresAt).
		Scan(&otp.ID, &otp.UserID, &otp.Code, &otp.Status, &otp.CreatedAt, &otp.ExpiresAt)

	if err != nil {
		return nil, err
	}

	return otp, nil
}

func (r *OTPRepository) GetOTPByUserID(userID string) (*models.OTP, error) {
	query := `
	SELECT id, user_id, code, status, created_at, expires_at
	FROM otp
	WHERE user_id = $1
	ORDER BY created_at DESC
	LIMIT 1
	`

	otp := &models.OTP{}
	err := r.db.QueryRow(query, userID).
		Scan(&otp.ID, &otp.UserID, &otp.Code, &otp.Status, &otp.CreatedAt, &otp.ExpiresAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return otp, nil
}

func (r *OTPRepository) UpdateOTPStatus(id int, status models.OTPStatus) error {
	query := `UPDATE otp SET status = $1 WHERE id = $2`
	_, err := r.db.Exec(query, status, id)
	return err
}
