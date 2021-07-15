package services

import (
	"context"
	"example.com/app/domain"
	"example.com/app/repo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService interface {
	GetAllUsers(string, context.Context) (*domain.UserResponse, error)
	DeleteByID(primitive.ObjectID) error
}

// DefaultUserService the service has a dependency of the repo
type DefaultUserService struct {
	repo repo.UserRepo
}

func (s DefaultUserService) GetAllUsers(page string, ctx context.Context) (*domain.UserResponse, error) {
	//childSpan := opentracing.StartSpan("child", opentracing.ChildOf(span.Context()))
	//defer childSpan.Finish()
	u, err := s.repo.FindAll(page, ctx)
	if err != nil {
		return nil, err
	}
	return  u, nil
}

func (s DefaultUserService) DeleteByID(id primitive.ObjectID) error {
	err := s.repo.DeleteByID(id)
	if err != nil {
		return err
	}
	return nil
}

func NewUserService(repository repo.UserRepo) DefaultUserService {
	return DefaultUserService{repository}
}