package services

import (
	"fmt"

	"github.com/mrrizkin/omniscan/app/models"
	"github.com/mrrizkin/omniscan/app/providers/hashing"
	"github.com/mrrizkin/omniscan/app/repositories"
)

type UserService struct {
	repo    *repositories.UserRepository
	hashing *hashing.Hashing
}

func (*UserService) Construct() interface{} {
	return func(repo *repositories.UserRepository, hashing *hashing.Hashing) *UserService {
		return &UserService{repo, hashing}
	}
}

func (s *UserService) Create(user *models.User) (*models.User, error) {
	if user.Password == nil {
		return nil, fmt.Errorf("password is required")
	}

	hash, err := s.hashing.Generate(*user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = &hash
	err = s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) FindAll(page, perPage int) (map[string]interface{}, error) {
	users, err := s.repo.FindAll(page, perPage)
	if err != nil {
		return nil, err
	}

	usersCount, err := s.repo.FindAllCount()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"result": users,
		"total":  int(usersCount),
	}, nil
}

func (s *UserService) FindByID(id uint) (*models.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Update(id uint, user *models.User) (*models.User, error) {
	userExist, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if user.Password != nil {
		if *user.Password != "" {
			hash, err := s.hashing.Generate(*user.Password)
			if err != nil {
				return nil, err
			}

			userExist.Password = &hash
		}
	}

	userExist.Name = user.Name
	userExist.Email = user.Email
	userExist.Username = user.Username
	userExist.RoleID = user.RoleID

	err = s.repo.Update(userExist)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Delete(userLogin *models.User, id uint) error {
	if userLogin.ID == uint(id) {
		return fmt.Errorf("cannot delete yourself")
	}

	return s.repo.Delete(id)
}
