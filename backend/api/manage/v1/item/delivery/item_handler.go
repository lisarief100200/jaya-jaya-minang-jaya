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
	IUsecase domain.ItemUsecase
}

// NewItemController represent item handler constructor
func NewItemController(gUser *gin.RouterGroup, iu domain.ItemUsecase) {
	handler := &Controllers{IUsecase: iu}

	gUser.GET("/item", handler.GetItems)
	gUser.POST("/item", handler.CreateItem)
	gUser.PUT("/item", handler.UpdateItem)
	gUser.DELETE("/item", handler.DeleteItem)
}

func (r *Controllers) GetItems(c *gin.Context) {
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

	respObj, err := r.IUsecase.GetItems(c, uid, level)
	if err != nil {
		log.Log.Errorf(constants.GetItem+constants.ErrorForRequestId, err.Error(), requestid.Get(c))
		resp := models.ResponseV2{RespCode: constants.InvalidSessionCode, RespMessage: "Bad Request", Response: nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	log.Log.Error(constants.GetItem+constants.SuccessReqestId, requestid.Get(c))
	resp := models.ResponseV2{RespCode: constants.SuccessCode, RespMessage: "Success", Response: respObj}
	c.JSON(http.StatusOK, resp)
}

func (r *Controllers) CreateItem(c *gin.Context) {
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
		log.Log.Error(constants.CreateItem, constants.ErrorForRequestId, "File Oversize", requestid.Get(c))
		res := models.ResponseV2{RespCode: constants.BadRequestCode, RespMessage: "Bad file", Response: nil}
		c.JSON(http.StatusOK, res)
		return
	}

	var req models.ReqCreateItem
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
	err = r.IUsecase.CreateItem(c, &req, uidInt)
	if err != nil {
		log.Log.Errorf(constants.CreateItem+constants.ErrorForRequestId, err.Error(), requestid.Get(c))
		resp := models.ResponseV2{RespCode: constants.BadRequestCode, RespMessage: "Bad Request", Response: nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	log.Log.Error(constants.CreateItem+constants.SuccessReqestId, requestid.Get(c))
	resp := models.ResponseV2{RespCode: constants.SuccessCode, RespMessage: "Success", Response: nil}
	c.JSON(http.StatusOK, resp)
}

func (r *Controllers) UpdateItem(c *gin.Context) {
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
		log.Log.Error(constants.UpdateItem, constants.ErrorForRequestId, "File Oversize", requestid.Get(c))
		res := models.ResponseV2{RespCode: constants.BadRequestCode, RespMessage: "Bad file", Response: nil}
		c.JSON(http.StatusOK, res)
		return
	}

	var req models.ReqUpdateItem
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
	err = r.IUsecase.UpdateItem(c, &req, uidInt)
	if err != nil {
		log.Log.Errorf(constants.UpdateItem+constants.ErrorForRequestId, err.Error(), requestid.Get(c))
		resp := models.ResponseV2{RespCode: constants.BadRequestCode, RespMessage: "Bad Request", Response: nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	log.Log.Error(constants.UpdateItem+constants.SuccessReqestId, requestid.Get(c))
	resp := models.ResponseV2{RespCode: constants.SuccessCode, RespMessage: "Success", Response: nil}
	c.JSON(http.StatusOK, resp)
}

func (r *Controllers) DeleteItem(c *gin.Context) {
	var (
		req models.ReqDeleteItem
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Log.Errorf(constants.DeleteItem+constants.ErrorForRequestId, err.Error(), requestid.Get(c))
		resp := models.ResponseV2{RespCode: constants.BadRequestCode, RespMessage: "Failed binding request from JSON", Response: nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	err := r.IUsecase.DeleteItem(c, req)
	if err != nil {
		log.Log.Errorf(constants.DeleteItem+constants.ErrorForRequestId, err.Error(), requestid.Get(c))
		resp := models.ResponseV2{RespCode: constants.BadRequestCode, RespMessage: "Bad Request", Response: nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	log.Log.Error(constants.DeleteItem+constants.SuccessReqestId, requestid.Get(c))
	resp := models.ResponseV2{RespCode: constants.SuccessCode, RespMessage: "Success", Response: nil}
	c.JSON(http.StatusOK, resp)
}
