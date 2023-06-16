package repository

type DatabaseRepo interface {
	// UploadImage()
	SignUp(username string, pHash string)
	// Authenticate()
	// SearchUserImage()
}
