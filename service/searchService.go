package service

import (
	"First/model"
	"First/repository"
)

type SearchService struct {
	repo repository.SearchRepository
}

func NewSearchService(repo repository.SearchRepository) *SearchService {
	return &SearchService{repo}
}

func (s *SearchService) SerachUser(key string) ([]model.User, error) {
	return s.repo.SearchUser(key)
}

func (s *SearchService) SerachPost(key string) ([]model.Post, error) {
	return s.repo.SearchPost(key)
}
