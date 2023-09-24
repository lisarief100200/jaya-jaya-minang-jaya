package domain

import (
	"api-lisa/models"
	"time"

	"github.com/gin-gonic/gin"
)

type AdminUsecase interface {
	CheckIdAdmin(c *gin.Context, req models.ReqLoginAdmin) (string, error)
	CheckLoginAdmin(c *gin.Context, req models.ReqLoginAdmin, uid string) error
	GetPassword(c *gin.Context, uid string, req models.ReqLoginAdmin) error
	GenerateSessionID(c *gin.Context, typeLogin string) (string, error)
	GetAdminProp(c *gin.Context, uid string) (models.AdminProp, error)
	UpdateLastLogin(c *gin.Context, uid string) error
	InsertSession(c *gin.Context, sessionId string, uid string, typeLogin string, token string, exp time.Time) error
	CreateItem(c *gin.Context, req *models.ReqCreateItem, uid int) error
	GetItems(c *gin.Context) ([]models.RespGetList, error)
	UpdateItem(c *gin.Context, req *models.ReqUpdateItem, uid int) error
	DeleteItem(c *gin.Context, req models.ReqDeleteItem) error
	LogoutAdmin(c *gin.Context, uid int) error
}

type AdminRepositories interface {
	CheckIdAdmin(c *gin.Context, username string) (string, error)
	CheckLoginAdmin(c *gin.Context, req models.ReqLoginAdmin, uid string) (exist bool, err error)
	GetPassword(c *gin.Context, uid string) (string, error)
	GenerateSessionID(c *gin.Context, processName string) (string, error)
	GetAdminProp(c *gin.Context, uid string) (models.AdminProp, error)
	UpdateLastLogin(c *gin.Context, uid string) error
	InsertSession(c *gin.Context, session, uid, tipe, token string, exptime time.Time) error
	CreateItem(c *gin.Context, req models.ReqCreateItem, imagePath string, uid int) error
	GetItems(c *gin.Context) ([]models.RespGetList, error)
	GetFilePath(c *gin.Context, req *models.ReqUpdateItem) (string, error)
	UpdateItem(c *gin.Context, req models.ReqUpdateItem, imagePath string, uid int) error
	DeleteItem(c *gin.Context, req models.ReqDeleteItem) error
	LogoutAdmin(c *gin.Context, uid int) error
}
