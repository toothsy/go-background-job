package models

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Job defines what workers pass around to execute and load to db
type Job struct {
	Id        string
	Status    int // refer to contants
	JobType   int // refer to contants
	User      UserPayload
	Image     ImagePayload
	CreatedAt time.Time
}

func (j Job) String() string {
	return fmt.Sprintf(`
	ID:%s
	Status:%d
	JobType:%d
	User:[%+v]
	Image:[%+v]
	Created at:[%+v]`, j.Id, j.Status, j.JobType, j.User, j.User, j.CreatedAt)
}

func (u UserPayload) String() string {
	return fmt.Sprintf(`
	Name:%s
	Email:%s
	PassWordHash:%s
	Created at :%s`, u.UserName, u.Email, u.PasswordHash, u.CreatedAt)
}

func (i ImagePayload) String() string {
	return fmt.Sprintf(`
	FileName:%s
	Email:%s
	Created at :%s`, i.Filename, i.Email, i.CreatedAt)
}

// used to hold the user data that will be inserted
type UserPayload struct {
	UserName     string    `json:"userName"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"passwordHash"`
	IsVerified   bool      `json:"isVerified"`
	CreatedAt    time.Time `json:"CreatedAt"`
}

// used to hold the image data that will be inserted
type ImagePayload struct {
	Email     string
	Filename  string
	FileData  primitive.Binary
	CreatedAt time.Time
}
