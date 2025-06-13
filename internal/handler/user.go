package handler

import (
	"github.com/gin-gonic/gin"
	"go-dianping/internal/service"
	"go-dianping/pkg/helper/resp"
	"net/http"
)

func NewUserHandler(handler *Handler,
	userService service.UserService,
) *UserHandler {
	return &UserHandler{
		Handler:     handler,
		userService: userService,
	}
}

type UserHandler struct {
	*Handler
	userService service.UserService
}

func (h *UserHandler) SendCode(ctx *gin.Context) {
	var params struct {
		Phone string `form:"phone" binding:"required"`
	}
	if err := ctx.ShouldBind(&params); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err := h.userService.SendCode(ctx, params.Phone)
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	resp.HandleSuccess(ctx, nil)
}

func (h *UserHandler) Login(ctx *gin.Context) {
	var params struct {
		Phone    string `json:"phone" binding:"required"`
		Code     string `json:"code"`
		Password string `json:"password"`
	}
	if err := ctx.ShouldBind(&params); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err := h.userService.Login(ctx, params.Phone, params.Code, params.Password)
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp.HandleSuccess(ctx, nil)
}
