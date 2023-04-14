package userapp

import (
	"homework8/internal/entities/user"
)

type App interface {
}

type app struct {
	repo user.Repository
}

func NewApp(repo user.Repository) App {
	return app{repo: repo}
}
