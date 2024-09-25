package user

import (
	"github.com/mrrizkin/omniscan/app/models"
	"github.com/mrrizkin/omniscan/system/database"
	"github.com/mrrizkin/omniscan/third-party/hashing"
)

type Repo struct {
	db *database.Database
}

type Service struct {
	repo    *Repo
	hashing hashing.Hashing
}

type PaginatedUser struct {
	Result []models.User
	Total  int
}
