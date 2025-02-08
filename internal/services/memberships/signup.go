package memberships

import (
	"errors"

	"github.com/amiulam/music-catalog/internal/models/memberships"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (s *service) SignUp(request memberships.SignUpRequest) error {
	existingUser, err := s.membershipRepo.GetUser(request.Email, request.Username, 0)

	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error().Err(err).Msg("error while retrieving user data")
		return err
	}

	if existingUser != nil {
		return errors.New("email or username already exist")
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	if err != nil {
		log.Error().Err(err).Msg("error hashed password")
		return err
	}

	model := memberships.User{
		Email:     request.Email,
		Password:  string(pass),
		Username:  request.Username,
		CreatedBy: request.Email,
		UpdatedBy: request.Email,
	}

	return s.membershipRepo.CreateUser(model)
}
