package domain

import (
	"backend/models"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthUsecase interface {
	CheckIdUser(c *gin.Context, req models.ReqLoginUser) (string, error)
	CheckLoginUser(c *gin.Context, req models.ReqLoginUser, uid string) error
	GetPassword(c *gin.Context, uid string, req models.ReqLoginUser) error
	GenerateSessionID(c *gin.Context, typeLogin string) (string, error)
	GetUserProp(c *gin.Context, uid string) (models.UserProp, error)
	UpdateLastLogin(c *gin.Context, uid string) error
	InsertSession(c *gin.Context, sessionId string, uid string, typeLogin string, token string, exp time.Time) error
	LogoutUser(c *gin.Context, uid int) error
}

type AuthRepositories interface {
	CheckIdUser(c *gin.Context, username string) (string, error)
	CheckLoginUser(c *gin.Context, req models.ReqLoginUser, uid string) (exist bool, err error)
	GetPassword(c *gin.Context, uid string) (string, error)
	GenerateSessionID(c *gin.Context, processName string) (string, error)
	GetUserProp(c *gin.Context, uid string) (models.UserProp, error)
	UpdateLastLogin(c *gin.Context, uid string) error
	InsertSession(c *gin.Context, session, uid, tipe, token string, exptime time.Time) error
	LogoutUser(c *gin.Context, uid int) error
}

type ItemUsecase interface {
	CreateItem(c *gin.Context, req *models.ReqCreateItem, uid int) error
	GetItems(c *gin.Context, uid, level string) ([]models.RespGetList, error)
	UpdateItem(c *gin.Context, req *models.ReqUpdateItem, uid int) error
	DeleteItem(c *gin.Context, req models.ReqDeleteItem) error
}

type ItemRepositories interface {
	CreateItem(c *gin.Context, req models.ReqCreateItem, imagePath []byte, uid int) error
	GetItems(c *gin.Context) ([]models.RespGetList, error)
	GetItemsById(c *gin.Context, uid int) ([]models.RespGetList, error)
	GetFilePath(c *gin.Context, req *models.ReqUpdateItem) (string, error)
	UpdateItem(c *gin.Context, req models.ReqUpdateItem, imagePath []byte, uid int) error
	DeleteItem(c *gin.Context, req models.ReqDeleteItem) error
}
