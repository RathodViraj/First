package service

import (
	"First/model"
	"First/repository"
	"fmt"
)

type CommentsService struct {
	repo repository.CommentsRepository
	gr   repository.Graph
}

func NewCommentsService(repo repository.CommentsRepository, gr repository.Graph) *CommentsService {
	return &CommentsService{repo, gr}
}

func (s *CommentsService) GetComments(id int) ([]model.Post, error) {
	return s.repo.GetComments(id)
}

func (s *CommentsService) CreateComment(comment *model.Post) error {
	if err := s.repo.AddComment(comment); err != nil {
		return err
	}

	if err := s.gr.CreatePostNode(comment.Id, comment.Uid, *comment.ParentId); err != nil {
		// s.repo.DeleteSQL(post.Id)
		return fmt.Errorf("Can't create post node in neo4j")
	}

	return nil
}

func (s *CommentsService) UpdateComment(comment *model.Post) error {
	return s.repo.UpdateComment(comment)
}
