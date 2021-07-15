package repo

import (
	"context"
	"example.com/app/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepo interface {
	FindAll(string, context.Context) (*domain.UserResponse, error)
	Create(user *domain.User) error
	UpdateByID(user *domain.User) error
	FindByUsername(string) (*domain.UserDto, error)
	DeleteByID(primitive.ObjectID) error
}
