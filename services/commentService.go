package services

import (
	"example.com/app/repo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentService interface {
	DeleteById(id primitive.ObjectID) error
}

type DefaultCommentService struct {
	repo repo.CommentRepo
}

func (c DefaultCommentService) DeleteById(id primitive.ObjectID) error {
	err := c.repo.DeleteById(id)
	if err != nil {
		return err
	}
	return nil
}

func NewCommentService(repository repo.CommentRepo) DefaultCommentService {
	return DefaultCommentService{repository}
}