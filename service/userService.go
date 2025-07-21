package service

import (
	"First/model"
	"First/repository"
)

type UserService struct {
	Repo repository.UserRepository
}

func NewUserService(Repo repository.UserRepository) *UserService {
	return &UserService{Repo}
}

func (s *UserService) RegisterUser(user *model.User) error {
	return s.Repo.Create(user)
}

func (s *UserService) GetUser(id int) (*model.User, error) {
	return s.Repo.GetByID(id)
}

func (s *UserService) DeleteUser(id, userID int, isAdmin bool) error {
	return s.Repo.Delete(id, userID, isAdmin)
}

func (s *UserService) UserFeed(id, offset int) ([]model.Post, error) {
	return s.Repo.GetUserFeed(id, offset)
}

func (s *UserService) GetUserByEmail(UserEmail string) (*model.User, error) {
	return s.Repo.GetByEmail(UserEmail)
}

func (s *UserService) UpdateUser(user *model.User) error {
	return s.Repo.Update(user)
}
