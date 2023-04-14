package userrepo

import (
	"homework8/internal/entities/user"
)

type repository struct {
	adDataById     map[int64]*user.User
	adDataByTitle  map[string]*user.User
	curIdGenerator int64
}

func New() user.Repository {
	return &repository{adDataById: make(map[int64]*user.User), adDataByTitle: make(map[string]*user.User), curIdGenerator: 0}
}
