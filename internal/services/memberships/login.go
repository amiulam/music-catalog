package memberships

import (
	"errors"

	"github.com/amiulam/music-catalog/internal/models/memberships"
	"github.com/amiulam/music-catalog/pkg/jwt"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (s *service) Login(request memberships.LoginRequest) (string, error) {
	userDetail, err := s.membershipRepo.GetUser(request.Email, "", 0)

	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error().Err(err).Msg("error while get user from database")
		return "", err
	}

	if userDetail == nil {
		return "", errors.New("user no exist")
	}

	err = bcrypt.CompareHashAndPassword([]byte(userDetail.Password), []byte(request.Password))

	if err != nil {
		return "", errors.New("invalid credentials")
	}

	accessToken, err := jwt.CreateToken(int64(userDetail.ID), userDetail.Username, s.cfg.Service.SecretJWT)

	if err != nil {
		log.Error().Err(err).Msg("fail to create JWT token")
		return "", err
	}

	return accessToken, nil

}
