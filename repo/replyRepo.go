package repo

import (
	"example.com/app/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ReplyRepo interface {
	Create(comment *domain.Reply) error
	UpdateById(id primitive.ObjectID, newContent string, edited bool, updatedTime time.Time) error
	DeleteById(id primitive.ObjectID, username string) error
}

