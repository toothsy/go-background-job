package repository

import "github/toothsy/go-background-job/internal/models"

type DatabaseRepo interface {
	// UploadImage()
	SignUp(user *models.UserPayload)
	FetchUser(user *models.UserPayload) (*models.UserPayload, error)
	UpdateUserVerification(user *models.UserPayload)
	// Authenticate()
	// SearchUserImage()
}
