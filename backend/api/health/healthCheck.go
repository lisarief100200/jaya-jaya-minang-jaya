package health

import (
	"api-lisa/models"
	"api-lisa/utils/constants"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, models.CreateResponseV2(c, constants.SuccessCode, constants.HealthCheck, constants.WarnHealthSuccess, nil))
}
