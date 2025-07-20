package service

import (
	"First/model"
	"First/repository"
)

type PostService struct {
	repo repository.PostRepository
}

func NewPostService(repo repository.PostRepository) *PostService {
	return &PostService{repo}
}

func (s *PostService) CreatePost(post *model.Post) error {
	return s.repo.Create(post)
}

func (s *PostService) DeletePost(id, uid int, isAdmin bool) error {
	return s.repo.Delete(id, uid, isAdmin)
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
