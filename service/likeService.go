package service

import (
	"First/model"
	"First/repository"
)

type LikeService struct {
	repo repository.LikeRepository
}

func NewLikeService(repo repository.LikeRepository) *LikeService {
	return &LikeService{repo}
}

func (s *LikeService) Like(like *model.Like) error {
	return s.repo.AddLike(like)
}

func (s *LikeService) Unlike(like *model.Like) error {
	return s.repo.RemoveLike(like)
}

func (s *LikeService) GetLikes(id int) ([]model.User, error) {
	return s.repo.WhoLiked(id)
}
