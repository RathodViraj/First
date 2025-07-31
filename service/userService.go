package service

import (
	"First/model"
	"First/repository"
	"fmt"
)

type UserService struct {
	repo repository.UserRepository
	gr   repository.Graph
}

func NewUserService(repo repository.UserRepository, gr repository.Graph) *UserService {
	return &UserService{repo: repo, gr: gr}
}

func (s *UserService) RegisterUser(user model.User) error {
	err := s.repo.CreateUserSQL(&user)
	if err != nil {
		return fmt.Errorf("failed to create user in SQL: %v", err)
	}

	err = s.gr.CreateUserNode(user.Id)
	if err != nil {
		_ = s.repo.Delete(user.Id, user.Id, true)
		return fmt.Errorf("failed to create user in Neo4j: %v", err)
	}

	return nil
}

func (s *UserService) GetUser(id int) (*model.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) DeleteUser(id, userID int, isAdmin bool) error {
	err := s.repo.Delete(id, userID, isAdmin)
	if err != nil {
		return err
	}

	return s.gr.DeleteUserNode(userID)
}

func (s *UserService) GetUserFeed(userID, offset int) ([]model.Post, error) {
	posts, err := s.repo.GetUserFeed(userID, offset, s.gr)
	if err != nil {
		return nil, fmt.Errorf("failed to get user feed: %w", err)
	}
	return posts, nil
}

func (s *UserService) GetUserByEmail(email string) (*model.User, error) {
	return s.repo.GetByEmail(email)
}

func (s *UserService) UpdateUser(user *model.User) error {
	return s.repo.Update(user)
}
