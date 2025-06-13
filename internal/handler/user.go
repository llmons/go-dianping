package handler

import (
	"github.com/gin-gonic/gin"
	"go-dianping/internal/dto"
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
	var params dto.LoginForm
	if err := ctx.ShouldBind(&params); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err := h.userService.Login(ctx, &params)
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp.HandleSuccess(ctx, nil)
}

func (h *UserHandler) Me(ctx *gin.Context) {
	user, err := h.userService.Me(ctx)
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	resp.HandleSuccess(ctx, user)
}
