package usecase

import (
	"backend/domain"
	"backend/models"
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthUsecase struct {
	ARepositories domain.AuthRepositories
}

func NewAuthUsecase(authRepo domain.AuthRepositories) domain.AuthUsecase {
	return &AuthUsecase{
		ARepositories: authRepo,
	}
}

func (r *AuthUsecase) CheckIdUser(c *gin.Context, req models.ReqLoginUser) (string, error) {
	uid, err := r.ARepositories.CheckIdUser(c, strings.TrimSpace(req.Username))
	if err != nil {
		return "", err
	}
	return uid, nil
}

func (r *AuthUsecase) CheckLoginUser(c *gin.Context, req models.ReqLoginUser, uid string) error {
	validUser, err := r.ARepositories.CheckLoginUser(c, req, uid)
	if err != nil {
		return err
	}

	if !validUser {
		return err
	}

	return nil
}

func (r *AuthUsecase) GetPassword(c *gin.Context, uid string, req models.ReqLoginUser) error {
	pass, err := r.ARepositories.GetPassword(c, uid)
	if err != nil {
		return err
	}

	//ps := helpers.Hashed(req.Password)
	ps := req.Password

	if ps != pass {
		return errors.New("password")
	}

	return nil
}

func (r *AuthUsecase) GenerateSessionID(c *gin.Context, typeLogin string) (string, error) {
	sessionId, err := r.ARepositories.GenerateSessionID(c, typeLogin)
	if err != nil {
		return "", err
	}
	return sessionId, nil
}

func (r *AuthUsecase) GetUserProp(c *gin.Context, uid string) (models.UserProp, error) {
	prop, err := r.ARepositories.GetUserProp(c, uid)
	if err != nil {
		return models.UserProp{}, err
	}
	return prop, nil
}

func (r *AuthUsecase) UpdateLastLogin(c *gin.Context, uid string) error {
	err := r.ARepositories.UpdateLastLogin(c, uid)
	if err != nil {
		return err
	}
	return nil
}

func (r *AuthUsecase) InsertSession(c *gin.Context, sessionId string, uid string, typeLogin string, token string, exp time.Time) error {
	err := r.ARepositories.InsertSession(c, sessionId, uid, typeLogin, token, exp)
	if err != nil {
		return err
	}
	return nil
}

func (r *AuthUsecase) LogoutUser(c *gin.Context, uid int) error {
	err := r.ARepositories.LogoutUser(c, uid)
	if err != nil {
		return err
	}
	return nil
}
