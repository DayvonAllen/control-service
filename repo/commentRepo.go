package repo

import (
	"example.com/app/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type CommentRepo interface {
	Create(comment *domain.Comment) error
	UpdateById(id primitive.ObjectID, newContent string, edited bool, updatedTime time.Time, username string) error
	DeleteById(id primitive.ObjectID) error
	DeleteManyById(id primitive.ObjectID) error
}
