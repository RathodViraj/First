package service

import (
	"First/model"
	"First/repository"
)

type ConnectionService struct {
	repo repository.ConnectionRepository
}

func NewConnectionService(repo repository.ConnectionRepository) *ConnectionService {
	return &ConnectionService{repo}
}

func (s *ConnectionService) Follow(followerID, followingID int) error {
	return s.repo.CreateConnection(&model.Connection{
		FollowerID:  followerID,
		FollowingID: followingID,
	})
}

func (s *ConnectionService) Unfollow(followerID, followingID int) error {
	return s.repo.DeleteConnection(followerID, followingID)
}

func (s *ConnectionService) GetFollowers(userID int) ([]model.User, error) {
	return s.repo.GetFollowers(userID)
}

func (s *ConnectionService) GetFollowings(userID int) ([]model.User, error) {
	return s.repo.GetFollowings(userID)
}

func (s *ConnectionService) Mutual(id int) ([]model.User, error) {
	return s.repo.GetMutual(id)
}
