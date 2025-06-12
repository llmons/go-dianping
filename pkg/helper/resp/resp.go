package resp

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type response struct {
	Success  bool        `json:"success"`
	ErrorMsg string      `json:"error_msg"`
	Data     interface{} `json:"data"`
	Total    int         `json:"total"`
}

func HandleSuccess(ctx *gin.Context, data interface{}) {
	if data == nil {
		data = map[string]string{}
	}
	resp := response{Success: true, ErrorMsg: "success", Data: data}
	ctx.JSON(http.StatusOK, resp)
}

func HandleListSuccess(ctx *gin.Context, data interface{}, total int) {
	if data == nil {
		data = map[string]string{}
	}
	resp := response{Success: true, ErrorMsg: "success", Data: data, Total: total}
	ctx.JSON(http.StatusOK, resp)
}

func HandleError(ctx *gin.Context, httpCode int, message string, data interface{}) {
	if data == nil {
		data = map[string]string{}
	}
	resp := response{Success: false, ErrorMsg: message, Data: data}
	ctx.JSON(httpCode, resp)
}
