package delivery

import (
	"api-lisa/controllers"
	"api-lisa/domain"
	"api-lisa/models"
	"api-lisa/utils/constants"
	"api-lisa/utils/log"
	"net/http"
	"strconv"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

type Controllers struct {
	AUsecase domain.AdminUsecase
}

// NewAdminController represent admin handler constructor
func NewAdminController(gPublic, gAdmin *gin.RouterGroup, au domain.AdminUsecase) {
	handler := &Controllers{AUsecase: au}
	gPublic.POST("/login/5", handler.LoginAdmin)
	gAdmin.GET("/logout", handler.LogoutAdmin)

	gAdmin.GET("/item", handler.GetItems)
	gAdmin.POST("/item", handler.CreateItem)
	gAdmin.PATCH("/item", handler.UpdateItem)
	gAdmin.DELETE("/item", handler.DeleteItem)
}

func (r *Controllers) LoginAdmin(c *gin.Context) {
	var (
		req     models.ReqLoginAdmin
		respObj models.RespLoginAdmin
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Log.Errorf(constants.LoginAdmin+constants.ErrorForRequestId, err.Error(), requestid.Get(c))
		resp := models.ResponseV2{RespCode: constants.BadRequestCode, RespMessage: "Failed binding request from JSON", Response: nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	// Get Id Admin
	uid, err := r.AUsecase.CheckIdAdmin(c, req)
	if err != nil {
		log.Log.Errorf(constants.LoginAdmin+constants.ErrorForRequestId, err.Error(), requestid.Get(c))
		resp := models.ResponseV2{RespCode: constants.InvalidSessionCode, RespMessage: "Invalid username or password", Response: nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	// Check whether admin is valid or not
	err = r.AUsecase.CheckLoginAdmin(c, req, uid)
	if err != nil {
		log.Log.Errorf(constants.LoginAdmin+constants.ErrorForRequestId, err.Error(), requestid.Get(c))
		resp := models.ResponseV2{RespCode: constants.InvalidSessionCode, RespMessage: "Invalid username or password", Response: nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	err = r.AUsecase.GetPassword(c, uid, req)
	if err != nil {
		log.Log.Errorf(constants.LoginAdmin+constants.ErrorForRequestId, err.Error(), requestid.Get(c))
		resp := models.ResponseV2{RespCode: constants.InvalidSessionCode, RespMessage: "Invalid username or password", Response: nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	var typeLogin string = "AUTH"

	// Generate sessionID
	sessionId, err := r.AUsecase.GenerateSessionID(c, typeLogin)
	if err != nil {
		log.Log.Errorf(constants.LoginAdmin+constants.ErrorForRequestId, err.Error(), requestid.Get(c))
		resp := models.ResponseV2{RespCode: constants.InvalidSessionCode, RespMessage: "Invalid username or password", Response: nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	// Prop admin
	prop, err := r.AUsecase.GetAdminProp(c, uid)
	if err != nil {
		log.Log.Errorf(constants.LoginAdmin+constants.ErrorForRequestId, err.Error(), requestid.Get(c))
		resp := models.ResponseV2{RespCode: constants.InvalidSessionCode, RespMessage: "Invalid username or password", Response: nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	// Generate token
	tokenString, exp, err := controllers.TokenAuthGenerator(req.Username, prop.Level, sessionId, uid)
	if err != nil {
		log.Log.Errorf(constants.LoginAdmin+constants.ErrorForRequestId, err.Error(), requestid.Get(c))
		resp := models.ResponseV2{RespCode: constants.InvalidSessionCode, RespMessage: "Invalid username or password", Response: nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	err = r.AUsecase.UpdateLastLogin(c, uid)
	if err != nil {
		log.Log.Errorf(constants.LoginAdmin+constants.ErrorForRequestId, err.Error(), requestid.Get(c))
		resp := models.ResponseV2{RespCode: constants.InvalidSessionCode, RespMessage: "Invalid username or password", Response: nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	err = r.AUsecase.InsertSession(c, sessionId, uid, typeLogin, tokenString, exp)
	if err != nil {
		log.Log.Errorf(constants.LoginAdmin+constants.ErrorForRequestId, err.Error(), requestid.Get(c))
		resp := models.ResponseV2{RespCode: constants.InvalidSessionCode, RespMessage: "Invalid username or password", Response: nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	respObj = models.RespLoginAdmin{
		Token:    tokenString,
		RoleCode: prop.Level,
	}

	log.Log.Error(constants.LoginAdmin+constants.SuccessReqestId, requestid.Get(c))
	resp := models.ResponseV2{RespCode: constants.SuccessCode, RespMessage: "Success", Response: respObj}
	c.JSON(http.StatusOK, resp)
}

func (r *Controllers) GetItems(c *gin.Context) {
	respObj, err := r.AUsecase.GetItems(c)
	if err != nil {
		log.Log.Errorf(constants.GetItem+constants.ErrorForRequestId, err.Error(), requestid.Get(c))
		resp := models.ResponseV2{RespCode: constants.InvalidSessionCode, RespMessage: "Invalid username or password", Response: nil}
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
	err = r.AUsecase.CreateItem(c, &req, uidInt)
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
	err = r.AUsecase.UpdateItem(c, &req, uidInt)
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

	err := r.AUsecase.DeleteItem(c, req)
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

func (r *Controllers) LogoutAdmin(c *gin.Context) {
	uid, err := controllers.GetUid(c)
	if err != nil {
		res := models.ResponseV2{RespCode: constants.InvalidSessionCode, RespMessage: "Invalid Session", Response: nil}
		c.JSON(http.StatusOK, res)
		return
	}

	uidInt, _ := strconv.Atoi(uid)
	err = r.AUsecase.LogoutAdmin(c, uidInt)
	if err != nil {
		log.Log.Errorf(constants.LogoutAdmin+constants.ErrorForRequestId, err.Error(), requestid.Get(c))
		resp := models.ResponseV2{RespCode: constants.BadRequestCode, RespMessage: "Bad Request", Response: nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	log.Log.Error(constants.LogoutAdmin+constants.SuccessReqestId, requestid.Get(c))
	resp := models.ResponseV2{RespCode: constants.SuccessCode, RespMessage: "Success", Response: nil}
	c.JSON(http.StatusOK, resp)
}
