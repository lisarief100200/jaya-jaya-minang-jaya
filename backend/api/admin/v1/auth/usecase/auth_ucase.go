package usecase

import (
	"backend/domain"
	"backend/models"
	"backend/utils/helpers"
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

func (r *AuthUsecase) CheckIdAdmin(c *gin.Context, req models.ReqLoginAdmin) (string, error) {
	uid, err := r.ARepositories.CheckIdAdmin(c, strings.TrimSpace(req.Username))
	if err != nil {
		return "", err
	}
	return uid, nil
}

func (r *AuthUsecase) CheckLoginAdmin(c *gin.Context, req models.ReqLoginAdmin, uid string) error {
	validAdmin, err := r.ARepositories.CheckLoginAdmin(c, req, uid)
	if err != nil {
		return err
	}

	if !validAdmin {
		return err
	}

	return nil
}

func (r *AuthUsecase) GetPassword(c *gin.Context, uid string, req models.ReqLoginAdmin) error {
	pass, err := r.ARepositories.GetPassword(c, uid)
	if err != nil {
		return err
	}

	ps := helpers.Hashed(req.Password)

	if ps != pass {
		return err
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

func (r *AuthUsecase) GetAdminProp(c *gin.Context, uid string) (models.AdminProp, error) {
	prop, err := r.ARepositories.GetAdminProp(c, uid)
	if err != nil {
		return models.AdminProp{}, err
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

func (r *AuthUsecase) LogoutAdmin(c *gin.Context, uid int) error {
	err := r.ARepositories.LogoutAdmin(c, uid)
	if err != nil {
		return err
	}
	return nil
}
