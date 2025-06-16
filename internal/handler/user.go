package handler

import (
	"github.com/gin-gonic/gin"
	"go-dianping/api"
	"go-dianping/internal/service"
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

// SendCode godoc
// @Summary 发送验证码
// @Schemes
// @Description
// @Tags user
// @Param phone query string true	"手机号"
// @Success 200 {object} api.SendCodeResp
// @Router /user/code [post]
func (h *UserHandler) SendCode(ctx *gin.Context) {
	var params api.SendCodeReq
	if err := ctx.ShouldBind(&params); err != nil {
		api.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err := h.userService.SendCode(ctx, params.Phone)
	if err != nil {
		api.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	api.HandleSuccess(ctx, nil)
}

// Login godoc
// @Summary 登录
// @Schemes
// @Description
// @Tags user
// @Accept json
// @Produce json
// @Security Bearer
// @Params request body api.LoginReq true "手机+验证码/密码"
// @Success 200 {object} api.LoginResp
// @Router /user [get]
func (h *UserHandler) Login(ctx *gin.Context) {
	var params api.LoginReq
	if err := ctx.ShouldBind(&params); err != nil {
		api.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	data, err := h.userService.Login(ctx, &params)
	if err != nil {
		api.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	api.HandleSuccess(ctx, data)
}

func (h *UserHandler) GetMe(ctx *gin.Context) {
	user, err := h.userService.GetMe(ctx)
	if err != nil {
		api.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	api.HandleSuccess(ctx, user)
}
