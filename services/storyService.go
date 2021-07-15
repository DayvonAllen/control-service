package services

import (
	"example.com/app/domain"
	"example.com/app/repo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StoryService interface {
	FindAll(string, bool) (*[]domain.Story, error)
	FindById(primitive.ObjectID) (*domain.StoryDto, error)
	DeleteById(primitive.ObjectID) error
}

type DefaultStoryService struct {
	repo repo.StoryRepo
}

func (s DefaultStoryService) FindAll(page string, newStoriesQuery bool) (*[]domain.Story, error) {
	story, err := s.repo.FindAll(page, newStoriesQuery)
	if err != nil {
		return nil, err
	}
	return story, nil
}

func (s DefaultStoryService) FindById(id primitive.ObjectID) (*domain.StoryDto, error) {
	story, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}
	return story, nil
}

func (s DefaultStoryService) DeleteById(id primitive.ObjectID) error {
	err := s.repo.DeleteById(id)
	if err != nil {
		return err
	}
	return nil
}

func NewStoryService(repository repo.StoryRepo) DefaultStoryService {
	return DefaultStoryService{repository}
}
