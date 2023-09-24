package controllers

import (
	"api-lisa/middleware"
	"api-lisa/models"
	"api-lisa/pkg/v1/mysql"
	"api-lisa/utils/config"
	"api-lisa/utils/constants"
	"api-lisa/utils/log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	if config.MyConfig.Environment == "PROD" {
		gin.SetMode(gin.ReleaseMode)
		log.Log.Printf("Starting %s on PRODUCTION Environment", config.MyConfig.AppName)
	} else {
		log.Log.Printf("Starting %s on DEVELOPMENT Environment", config.MyConfig.AppName)
	}

	middleware.SetupLogger()
	err := mysql.InitDBConnection()
	if err != nil {
		log.Log.Fatal(err.Error())
	}
}

func HandleNoRoutes(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, models.CreateResponse(c, http.StatusNotFound, constants.REJECT, constants.UndefinedProcess, constants.WarnUndefinedProcess, nil))
}

func HandleNoMethod(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, models.CreateResponse(c, http.StatusMethodNotAllowed, constants.REJECT, constants.UndefinedMethods, constants.WarnUndefinedMethod, nil))
}

func HandlePanic(c *gin.Context, err interface{}) {
	log.Log.Error(err.(error).Error())
	c.JSON(http.StatusInternalServerError, models.CreateResponse(c, http.StatusInternalServerError, constants.FAILED, constants.InternalServerError, constants.WarnInternalError, nil))
}
