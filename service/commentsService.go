package service

import (
	"First/model"
	"First/repository"
)

type CommentsService struct {
	repo repository.CommentsRepository
}

func NewCommentsService(repo repository.CommentsRepository) *CommentsService {
	return &CommentsService{repo}
}

func (s *CommentsService) GetComments(id int) ([]model.Post, error) {
	return s.repo.GetComments(id)
}

func (s *CommentsService) CreateComment(comment *model.Post) error {
	return s.repo.AddComment(comment)
}

func (s *CommentsService) UpdateComment(comment *model.Post) error {
	return s.repo.UpdateComment(comment)
}
