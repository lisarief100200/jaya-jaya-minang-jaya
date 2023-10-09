package delivery

import (
	"backend/controllers"
	"backend/domain"
	"backend/models"
	"backend/utils/constants"
	"backend/utils/log"
	"net/http"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

type Controllers struct {
	CFUsecase domain.CashFlowUsecase
}

// NewCashFlowController represent cash flow handler
func NewCashFlowController(gUser *gin.RouterGroup, cfu domain.CashFlowUsecase) {
	handler := &Controllers{CFUsecase: cfu}
	gUser.GET("/cash-flow", handler.GetCashFlow)
}

func (r *Controllers) GetCashFlow(c *gin.Context) {
	uid, err := controllers.GetUid(c)
	if err != nil {
		res := models.ResponseV2{RespCode: constants.InvalidSessionCode, RespMessage: "Invalid Session", Response: nil}
		c.JSON(http.StatusOK, res)
		return
	}

	level, err := controllers.GetLevel(c)
	if err != nil {
		res := models.ResponseV2{RespCode: constants.InvalidSessionCode, RespMessage: "Invalid Session", Response: nil}
		c.JSON(http.StatusOK, res)
		return
	}

	respObj, err := r.CFUsecase.GetCashFlow(c, uid, level)
	if err != nil {
		log.Log.Errorf(constants.GetCashFlow+constants.ErrorForRequestId, err.Error(), requestid.Get(c))
		resp := models.ResponseV2{RespCode: constants.InvalidSessionCode, RespMessage: "Bad Request", Response: nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	log.Log.Error(constants.GetCashFlow+constants.SuccessReqestId, requestid.Get(c))
	resp := models.ResponseV2{RespCode: constants.SuccessCode, RespMessage: "Success", Response: respObj}
	c.JSON(http.StatusOK, resp)
}
