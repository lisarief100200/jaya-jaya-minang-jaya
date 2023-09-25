package delivery

import (
	"backend/controllers"
	"backend/domain"
	"backend/models"
	"backend/utils/constants"
	"backend/utils/log"
	"backend/utils/strings"
	"net/http"
	"strconv"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

type Controllers struct {
	AUsecase domain.AuthUsecase
}

// NewAuthController represent auth handler
func NewAuthController(gPublic, gAdmin *gin.RouterGroup, au domain.AuthUsecase) {
	handler := &Controllers{AUsecase: au}
	gPublic.POST("/login/5", handler.LoginAdmin)
	gAdmin.GET("/logout", handler.LogoutAdmin)
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
		resp := models.ResponseV2{RespCode: constants.InvalidSessionCode, RespMessage: "Invalid username or passowrd", Response: nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	// Check whether admin is valid or not
	err = r.AUsecase.CheckLoginAdmin(c, req, uid)
	if err != nil {
		log.Log.Errorf(constants.LoginAdmin+constants.ErrorForRequestId, err.Error(), requestid.Get(c))
		resp := models.ResponseV2{RespCode: constants.InvalidSessionCode, RespMessage: "Invalid username or passowrd", Response: nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	err = r.AUsecase.GetPassword(c, uid, req)
	if err != nil {
		if strings.IsContains(err.Error(), "password") {
			log.Log.Errorf(constants.LoginAdmin+constants.ErrorForRequestId, err.Error(), requestid.Get(c))
			resp := models.ResponseV2{RespCode: constants.InvalidSessionCode, RespMessage: "Invalid password", Response: nil}
			c.JSON(http.StatusOK, resp)
			return
		}
		log.Log.Errorf(constants.LoginAdmin+constants.ErrorForRequestId, err.Error(), requestid.Get(c))
		resp := models.ResponseV2{RespCode: constants.InvalidSessionCode, RespMessage: "Invalid username or password", Response: nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	var typeLogin string = "AUTH"

	// Generate sessionID
	sessionId, err := r.AUsecase.GenerateSessionID(c, typeLogin)
	if err != nil {
		log.Log.Error(constants.LoginAdmin+constants.ErrorForRequestId, err.Error(), requestid.Get(c))
		resp := models.ResponseV2{RespCode: constants.InvalidSessionCode, RespMessage: "Invalid username or password", Response: nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	// Prop Admin
	prop, err := r.AUsecase.GetAdminProp(c, uid)
	if err != nil {
		log.Log.Errorf(constants.LoginAdmin+constants.ErrorForRequestId, err.Error(), requestid.Get(c))
		resp := models.ResponseV2{RespCode: constants.InvalidSessionCode, RespMessage: "Invalid username or passowrd", Response: nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	// Generate token
	tokenString, exp, err := controllers.TokenAuthGenerator(req.Username, prop.Level, sessionId, uid)
	if err != nil {
		log.Log.Errorf(constants.LoginAdmin+constants.ErrorForRequestId, err.Error(), requestid.Get(c))
		resp := models.ResponseV2{RespCode: constants.InvalidSessionCode, RespMessage: "Invalid username or passowrd", Response: nil}
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
