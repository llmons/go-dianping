package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Success  bool   `json:"success"`
	ErrorMsg string `json:"errorMsg,omitempty"`
	Data     any    `json:"data,omitempty"`
	Total    int    `json:"total,omitempty"`
}

func HandleSuccess(ctx *gin.Context, data any) {
	resp := Response{Success: true, Data: data}
	ctx.JSON(http.StatusOK, resp)
}

func HandleListSuccess(ctx *gin.Context, data any, total int) {
	resp := Response{Success: true, Data: data, Total: total}
	ctx.JSON(http.StatusOK, resp)
}

func HandleError(ctx *gin.Context, httpCode int, message string, data any) {
	resp := Response{Success: false, ErrorMsg: message, Data: data}
	ctx.JSON(httpCode, resp)
}

type Error struct {
	Code    int
	Message string
}

var errorCodeMap = map[error]int{}

func newError(code int, msg string) error {
	err := errors.New(msg)
	errorCodeMap[err] = code
	return err
}
func (e Error) Error() string {
	return e.Message
}
