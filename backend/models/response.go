package models

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Response struct {
	StatusCode int         `json:"statusCode"`
	Status     string      `json:"status"`
	Message    string      `json:"message"`
	DisplayMsg string      `json:"displayMessage"`
	Response   interface{} `json:"response"`
}
type ResponseV2 struct {
	RespCode    string      `json:"respCode"`
	RespMessage string      `json:"respMessage"`
	IdMessage   string      `json:"idMessage"`
	Response    interface{} `json:"response"`
	//Error       Error       `json:"error,omitempty"`
}

type Error struct {
	ErrorCode    int    `json:"errorCode,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

func CreateResponse(c *gin.Context, statusCode int, status, message, displayMessage string, response interface{}) Response {
	c.Set("RespCode", strconv.Itoa(statusCode))
	return Response{StatusCode: statusCode, Status: status, Message: message, DisplayMsg: displayMessage, Response: response}
}

func CreateResponseV2(c *gin.Context, respCode, RespMessage, idMessage string, response interface{}) ResponseV2 {
	c.Set("RespCode", respCode)
	return ResponseV2{RespCode: respCode, RespMessage: RespMessage, IdMessage: idMessage, Response: response}
}
