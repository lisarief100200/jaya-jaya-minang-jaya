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
func NewAuthController(gPublic, gUser *gin.RouterGroup, au domain.AuthUsecase) {
	handler := &Controllers{AUsecase: au}
	gPublic.POST("/login/5", handler.LoginUser)
	gUser.GET("/logout", handler.LogoutUser)
}

func (r *Controllers) LoginUser(c *gin.Context) {
	var (
		req     models.ReqLoginUser
		respObj models.RespLoginUser
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Log.Errorf(constants.LoginUser+constants.ErrorForRequestId, err.Error(), requestid.Get(c))
		resp := models.ResponseV2{RespCode: constants.BadRequestCode, RespMessage: "Failed binding request from JSON", Response: nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	// Get Id User
	uid, err := r.AUsecase.CheckIdUser(c, req)
	if err != nil {
		log.Log.Errorf(constants.LoginUser+constants.ErrorForRequestId, err.Error(), requestid.Get(c))
		resp := models.ResponseV2{RespCode: constants.InvalidSessionCode, RespMessage: "Invalid username or passowrd", Response: nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	// Check whether user is valid or not
	err = r.AUsecase.CheckLoginUser(c, req, uid)
	if err != nil {
		log.Log.Errorf(constants.LoginUser+constants.ErrorForRequestId, err.Error(), requestid.Get(c))
		resp := models.ResponseV2{RespCode: constants.InvalidSessionCode, RespMessage: "Invalid username or passowrd", Response: nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	err = r.AUsecase.GetPassword(c, uid, req)
	if err != nil {
		if strings.IsContains(err.Error(), "password") {
			log.Log.Errorf(constants.LoginUser+constants.ErrorForRequestId, err.Error(), requestid.Get(c))
			resp := models.ResponseV2{RespCode: constants.InvalidSessionCode, RespMessage: "Invalid password", Response: nil}
			c.JSON(http.StatusOK, resp)
			return
		}
		log.Log.Errorf(constants.LoginUser+constants.ErrorForRequestId, err.Error(), requestid.Get(c))
		resp := models.ResponseV2{RespCode: constants.InvalidSessionCode, RespMessage: "Invalid username or password", Response: nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	var typeLogin string = "AUTH"

	// Generate sessionID
	sessionId, err := r.AUsecase.GenerateSessionID(c, typeLogin)
	if err != nil {
		log.Log.Error(constants.LoginUser+constants.ErrorForRequestId, err.Error(), requestid.Get(c))
		resp := models.ResponseV2{RespCode: constants.InvalidSessionCode, RespMessage: "Invalid username or password", Response: nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	// Prop User
	prop, err := r.AUsecase.GetUserProp(c, uid)
	if err != nil {
		log.Log.Errorf(constants.LoginUser+constants.ErrorForRequestId, err.Error(), requestid.Get(c))
		resp := models.ResponseV2{RespCode: constants.InvalidSessionCode, RespMessage: "Invalid username or passowrd", Response: nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	// Generate token
	tokenString, exp, err := controllers.TokenAuthGenerator(req.Username, prop.Level, sessionId, uid)
	if err != nil {
		log.Log.Errorf(constants.LoginUser+constants.ErrorForRequestId, err.Error(), requestid.Get(c))
		resp := models.ResponseV2{RespCode: constants.InvalidSessionCode, RespMessage: "Invalid username or passowrd", Response: nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	err = r.AUsecase.UpdateLastLogin(c, uid)
	if err != nil {
		log.Log.Errorf(constants.LoginUser+constants.ErrorForRequestId, err.Error(), requestid.Get(c))
		resp := models.ResponseV2{RespCode: constants.InvalidSessionCode, RespMessage: "Invalid username or password", Response: nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	err = r.AUsecase.InsertSession(c, sessionId, uid, typeLogin, tokenString, exp)
	if err != nil {
		log.Log.Errorf(constants.LoginUser+constants.ErrorForRequestId, err.Error(), requestid.Get(c))
		resp := models.ResponseV2{RespCode: constants.InvalidSessionCode, RespMessage: "Invalid username or password", Response: nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	respObj = models.RespLoginUser{
		Token:    tokenString,
		RoleCode: prop.Level,
	}

	log.Log.Error(constants.LoginUser+constants.SuccessReqestId, requestid.Get(c))
	resp := models.ResponseV2{RespCode: constants.SuccessCode, RespMessage: "Success", Response: respObj}
	c.JSON(http.StatusOK, resp)
}

func (r *Controllers) LogoutUser(c *gin.Context) {
	uid, err := controllers.GetUid(c)
	if err != nil {
		res := models.ResponseV2{RespCode: constants.InvalidSessionCode, RespMessage: "Invalid Session", Response: nil}
		c.JSON(http.StatusOK, res)
		return
	}

	uidInt, _ := strconv.Atoi(uid)
	err = r.AUsecase.LogoutUser(c, uidInt)
	if err != nil {
		log.Log.Errorf(constants.LogoutUser+constants.ErrorForRequestId, err.Error(), requestid.Get(c))
		resp := models.ResponseV2{RespCode: constants.BadRequestCode, RespMessage: "Bad Request", Response: nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	log.Log.Error(constants.LogoutUser+constants.SuccessReqestId, requestid.Get(c))
	resp := models.ResponseV2{RespCode: constants.SuccessCode, RespMessage: "Success", Response: nil}
	c.JSON(http.StatusOK, resp)
}
