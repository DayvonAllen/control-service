package repo

import (
	"example.com/app/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StoryRepo interface {
	FindAll(string, bool) (*[]domain.Story, error)
	FindById(primitive.ObjectID) (*domain.StoryDto, error)
	Create(story *domain.Story) error
	UpdateById(primitive.ObjectID, string, string, string, *[]domain.Tag, bool) error
	DeleteById(primitive.ObjectID) error
}
