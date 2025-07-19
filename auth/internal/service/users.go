package service

import "auth/internal/repository"

type UsersService struct {
	store *repository.Storage
}
