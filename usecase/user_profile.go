package usecase

import "github.com/zackwn/article-api/entity"

func NewUserProfileUseCase(userRepo UserRepository) *UserProfileUseCase {
	return &UserProfileUseCase{userRepository: userRepo}
}

type UserProfileUseCase struct {
	userRepository UserRepository
}

type UserProfileDTO struct {
	ID string
}

func (userProfileUseCase UserProfileUseCase) Exec(dto *UserProfileDTO) (*entity.User, UseCaseErr) {
	user, found := userProfileUseCase.userRepository.FindByID(dto.ID)
	if !found {
		return nil, ErrUserDoNotExists{}
	}
	// occult password
	user.Password = ""
	return user, nil
}
