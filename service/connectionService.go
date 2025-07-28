package service

import (
	"First/model"
	"First/repository"
)

type ConnectionService struct {
	repo  repository.ConnectionRepository
	graph repository.Graph
}

func NewConnectionService(repo repository.ConnectionRepository, graph repository.Graph) *ConnectionService {
	return &ConnectionService{repo: repo, graph: graph}
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
	followerIDs, err := s.graph.GetFollowersIDs(userID)
	if err != nil {
		return []model.User{}, err
	}

	return s.repo.GetFollowers(followerIDs), nil
}

func (s *ConnectionService) GetFollowings(userID int) ([]model.User, error) {
	followingIDs, err := s.graph.GetFollowingsIDs(userID)
	if err != nil {
		return []model.User{}, err
	}

	return s.repo.GetFollowings(followingIDs), nil
}

func (s *ConnectionService) Mutual(userID int) ([]model.User, error) {
	mutualIDs, err := s.graph.GetMutualIDs(userID)
	if err != nil {
		return []model.User{}, err
	}

	return s.repo.GetMutual(mutualIDs)
}
