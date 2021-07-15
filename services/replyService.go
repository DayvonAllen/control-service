package services

import (
	"example.com/app/repo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReplyService interface {
	DeleteById(id primitive.ObjectID, username string) error
}

type DefaultReplyService struct {
	repo repo.ReplyRepo
}

func (r DefaultReplyService) DeleteById(id primitive.ObjectID, username string) error {
	err := r.repo.DeleteById(id, username)
	if err != nil {
		return err
	}
	return nil
}

func NewReplyService(repository repo.ReplyRepo) DefaultReplyService {
	return DefaultReplyService{repository}
}
