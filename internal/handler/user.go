package handler

import (
	"github.com/gin-gonic/gin"
	"go-dianping/api/v1"
	"go-dianping/internal/service"
	"go.uber.org/zap"
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
// @Produce json
// @Param phone query string true "手机号"
// @Success 200 {object} api.SendCodeResp
// @Router /user/code [post]
func (h *UserHandler) SendCode(ctx *gin.Context) {
	var params v1.SendCodeReq
	if err := ctx.ShouldBind(&params); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err := h.userService.SendCode(ctx.Request.Context(), &params)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// Login godoc
// @Summary 登录
// @Schemes
// @Description
// @Tags user
// @Accept json
// @Produce json
// @Params request body api.LoginReq true "手机+验证码"
// @Success 200 {object} api.LoginResp
// @Router /user/login [post]
func (h *UserHandler) Login(ctx *gin.Context) {
	var params v1.LoginReq
	if err := ctx.ShouldBind(&params); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	data, err := h.userService.Login(ctx.Request.Context(), &params)
	h.logger.Info("Login", zap.Any("data", data))
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	v1.HandleSuccess(ctx, data)
}

// GetMe godoc
// @Summary 获取当前登录的用户信息
// @Schemes
// @Description
// @Tags user
// @Produce json
// @Security Bearer
// @Success 200 {object} api.GetMeResp
// @Router /user/me [get]
func (h *UserHandler) GetMe(ctx *gin.Context) {
	user, err := h.userService.GetMe(ctx.Request.Context())
	h.logger.Info("Get", zap.Any("user", user))
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	v1.HandleSuccess(ctx, user)
}
