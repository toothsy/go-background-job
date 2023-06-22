package helpers

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"github/toothsy/go-background-job/internal/models"
	"log"
	"os"
)

func GenerateToken(u *models.UserPayload) string {
	secretKey := os.Getenv("SECRET")
	data := []byte(u.Email + u.UserName + secretKey)
	key := []byte(secretKey)

	// Generate OTP token using HMAC-SHA256
	hmacSHA256 := hmac.New(sha256.New, key)
	hmacSHA256.Write(data)
	token := hex.EncodeToString(hmacSHA256.Sum(nil))

	return token
}

func VerifyToken(token string, u *models.UserPayload) bool {
	recreatedToken := GenerateToken(u)
	return subtle.ConstantTimeCompare([]byte(token), []byte(recreatedToken)) == 1
}

func ReadFile(filePath string) (string, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("Error in reading file:", err)
		return "", err
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Read the file line by line and append to a string
	var content string
	for scanner.Scan() {
		content += scanner.Text() + "\n"
	}
	return content, nil
}
