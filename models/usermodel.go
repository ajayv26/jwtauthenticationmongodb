package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"id"`
	FirstName    *string            `json:"firstName" validate:"required, min=2, max=20"`
	LastName     *string            `json:"lastName" validate:"required, min=2, max=20"`
	Password     *string            `json:"password" validate:"required, min=2"`
	Email        *string            `json:"email" validate:"email, required"`
	Phone        *string            `json:"phone" validate:"required"`
	Token        *string            `json:"token"`
	UserType     *string            `json:"userType" validate:"required, eq=ADMIN|ew=USER "`
	RefreshToken *string            `json:"refreshToken"`
	CreatedAt    time.Time          `json:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt"`
	UserID       string             `json:"userID"`
}
