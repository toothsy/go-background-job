package repository

import (
	"encoding/json"
	"github/toothsy/go-background-job/internal/models"
)

type DatabaseRepo interface {
	// UploadImage()
	// inserts the user parameter in to the database
	SignUp(user *models.UserPayload)
	// FetchUser fetches the user based on the email filter
	FetchUser(user *models.UserPayload) (*models.UserPayload, error)

	// marks the user verified if the generated uuid matches one sent via get request
	UpdateUserVerification(user *models.UserPayload)
	// Authenticate()
	// searches image using email and returns array of images that match
	SearchUserImage(email string) ([]json.RawMessage, error)
}
