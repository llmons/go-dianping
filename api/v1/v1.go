package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type response struct {
	Success  bool   `json:"success"`
	ErrorMsg string `json:"error_msg,omitempty"`
	Data     any    `json:"data,omitempty"`
	Total    int    `json:"total,omitempty"`
}

func HandleSuccess(ctx *gin.Context, data any) {
	resp := response{Success: true, Data: data}
	ctx.JSON(http.StatusOK, resp)
}

func HandleListSuccess(ctx *gin.Context, data any, total int) {
	resp := response{Success: true, Data: data, Total: total}
	ctx.JSON(http.StatusOK, resp)
}

func HandleError(ctx *gin.Context, httpCode int, message string, data any) {
	resp := response{Success: false, ErrorMsg: message, Data: data}
	ctx.JSON(httpCode, resp)
}
