package delivery

import (
	"backend/controllers"
	"backend/domain"
	"backend/models"
	"backend/utils/constants"
	"backend/utils/log"
	"net/http"
	"strconv"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

type Controllers struct {
	UUsecase domain.UtilitiesUsecase
}

// NewUtilitiesController represent utilities handler constructor
func NewUtilitiesController(gUser *gin.RouterGroup, uu domain.UtilitiesUsecase) {
	handler := &Controllers{UUsecase: uu}

	gUser.GET("/utilities", handler.GetUtilities)
	gUser.POST("/utilities", handler.CreateUtilities)
	gUser.PUT("/utilities", handler.UpdateUtilities)
	gUser.DELETE("/utilities", handler.DeleteUtilities)
}

func (r *Controllers) GetUtilities(c *gin.Context) {
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

	respObj, err := r.UUsecase.GetUtilities(c, uid, level)
	if err != nil {
		log.Log.Errorf(constants.GetUtilities+constants.ErrorForRequestId, err.Error(), requestid.Get(c))
		resp := models.ResponseV2{RespCode: constants.InvalidSessionCode, RespMessage: "Bad Request", Response: nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	log.Log.Error(constants.GetUtilities+constants.SuccessReqestId, requestid.Get(c))
	resp := models.ResponseV2{RespCode: constants.SuccessCode, RespMessage: "Success", Response: respObj}
	c.JSON(http.StatusOK, resp)
}

func (r *Controllers) CreateUtilities(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		log.Log.Error(err.Error())
		res := models.ResponseV2{RespCode: constants.BadRequestCode, RespMessage: "Bad Request", Response: nil}
		c.JSON(http.StatusOK, res)
		return
	}

	files := form.File["image"]
	var fileSize int64
	for _, v := range files {
		fileSize += v.Size
	}

	if fileSize >= 10000000 {
		log.Log.Error(constants.CreateUtilities, constants.ErrorForRequestId, "File Oversize", requestid.Get(c))
		res := models.ResponseV2{RespCode: constants.BadRequestCode, RespMessage: "Bad file", Response: nil}
		c.JSON(http.StatusOK, res)
		return
	}

	var req models.ReqCreateUtilities
	err = c.Bind(&req)
	if err != nil {
		log.Log.Error(err.Error())
		res := models.ResponseV2{RespCode: constants.BadRequestCode, RespMessage: "Bad Request", Response: nil}
		c.JSON(http.StatusOK, res)
		return
	}

	uid, err := controllers.GetUid(c)
	if err != nil {
		res := models.ResponseV2{RespCode: constants.InvalidSessionCode, RespMessage: "Invalid Session", Response: nil}
		c.JSON(http.StatusOK, res)
		return
	}

	uidInt, _ := strconv.Atoi(uid)
	err = r.UUsecase.CreateUtilities(c, &req, uidInt)
	if err != nil {
		log.Log.Errorf(constants.CreateUtilities+constants.ErrorForRequestId, err.Error(), requestid.Get(c))
		resp := models.ResponseV2{RespCode: constants.SuccessCode, RespMessage: "Success", Response: nil}
		c.JSON(http.StatusOK, resp)
	}
}

func (r *Controllers) UpdateUtilities(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		log.Log.Error(err.Error())
		res := models.ResponseV2{RespCode: constants.BadRequestCode, RespMessage: "Bad request", Response: nil}
		c.JSON(http.StatusOK, res)
		return
	}

	files := form.File["image"]
	var fileSize int64
	for _, v := range files {
		fileSize += v.Size
	}

	if fileSize >= 10000000 {
		log.Log.Error(constants.UpdateUtilities, constants.ErrorForRequestId, "File Oversize", requestid.Get(c))
		res := models.ResponseV2{RespCode: constants.BadRequestCode, RespMessage: "Bad file", Response: nil}
		c.JSON(http.StatusOK, res)
		return
	}

	var req models.ReqUpdateUtilities
	err = c.Bind(&req)
	if err != nil {
		log.Log.Error(err.Error())
		res := models.ResponseV2{RespCode: constants.BadRequestCode, RespMessage: "Bad Request", Response: nil}
		c.JSON(http.StatusOK, res)
		return
	}

	uid, err := controllers.GetUid(c)
	if err != nil {
		res := models.ResponseV2{RespCode: constants.InvalidSessionCode, RespMessage: "Invalid Session", Response: nil}
		c.JSON(http.StatusOK, res)
		return
	}

	uidInt, _ := strconv.Atoi(uid)
	err = r.UUsecase.UpdateUtilities(c, &req, uidInt)
	if err != nil {
		log.Log.Errorf(constants.UpdateUtilities+constants.ErrorForRequestId, err.Error(), requestid.Get(c))
		resp := models.ResponseV2{RespCode: constants.BadRequestCode, RespMessage: "Bad Request", Response: nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	log.Log.Error(constants.UpdateUtilities+constants.SuccessReqestId, requestid.Get(c))
	resp := models.ResponseV2{RespCode: constants.SuccessCode, RespMessage: "Success", Response: nil}
	c.JSON(http.StatusOK, resp)
}

func (r *Controllers) DeleteUtilities(c *gin.Context) {
	var (
		req models.ReqDeleteUtilities
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Log.Errorf(constants.DeleteUtilities+constants.ErrorForRequestId, err.Error(), requestid.Get(c))
		resp := models.ResponseV2{RespCode: constants.BadRequestCode, RespMessage: "Failed binding requests from JSON", Response: nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	err := r.UUsecase.DeleteUtilities(c, req)
	if err != nil {
		log.Log.Errorf(constants.DeleteUtilities+constants.ErrorForRequestId, err.Error(), requestid.Get(c))
		resp := models.ResponseV2{RespCode: constants.BadRequestCode, RespMessage: "Bad Request", Response: nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	log.Log.Error(constants.DeleteUtilities+constants.SuccessReqestId, requestid.Get(c))
	resp := models.ResponseV2{RespCode: constants.SuccessCode, RespMessage: "Success", Response: nil}
	c.JSON(http.StatusOK, resp)
}
