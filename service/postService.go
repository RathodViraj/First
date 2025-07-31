package service

import (
	"First/model"
	"First/repository"
)

type PostService struct {
	repo repository.PostRepository
	gr   repository.Graph
}

func NewPostService(repo repository.PostRepository, gr repository.Graph) *PostService {
	return &PostService{repo, gr}
}

func (s *PostService) CreatePost(post *model.Post) error {
	if err := s.repo.Create(post); err != nil {
		return err
	}

	if err := s.gr.CreatePostNode(post.Id, post.Uid, -1); err != nil {
		s.repo.Delete(post.Id, post.Uid, true)
		return err
	}

	return nil
}

func (s *PostService) DeletePost(id, uid int, isAdmin bool) error {
	if err := s.repo.Delete(id, uid, isAdmin); err != nil {
		return err
	}

	if err := s.gr.DeletePostNode(id); err != nil {
		return err
	}

	return nil
}

func (s *PostService) GetPost(id int) (*model.Post, error) {
	return s.repo.GetByID(id)
}

func (s *PostService) NewPosts(offset int) ([]model.Post, error) {
	return s.repo.GetRecentPosts(offset)
}

func (s *PostService) UserPosts(id int) (*[]model.Post, error) {
	return s.repo.GetAllUserPosts(id)
}

func (s *PostService) UpdatePost(post *model.Post) error {
	return s.repo.Update(post)
}
