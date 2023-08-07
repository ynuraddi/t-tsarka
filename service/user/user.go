package user

import (
	"github.com/ynuraddi/t-tsarka/repository"
)

type userService struct {
	repository.IUserRepository
}

func New(repo repository.IUserRepository) *userService {
	return &userService{IUserRepository: repo}
}

// func (s *userService) Create(ctx context.Context, user model.User) (id int64, err error) {
// 	return s.repo.Create(ctx, user)
// }

// func (s *userService) Get(ctx context.Context, id int64) (user model.User, err error) {
// 	return s.repo.Get(ctx, id)
// }

// func (s *userService) Update(ctx context.Context, user model.User) (dbuser model.User, err error) {
// 	return s.repo.Update(ctx, user)
// }

// func (s *userService) Delete(ctx context.Context, id int64) error {
// 	return s.repo.Delete(ctx, id)
// }
