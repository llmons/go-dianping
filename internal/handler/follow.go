package handler

import (
	"github.com/gin-gonic/gin"
	v1 "go-dianping/api/v1"
	"go-dianping/internal/service"
	"net/http"
)

type FollowHandler struct {
	*Handler
	followService service.FollowService
}

func NewFollowHandler(
	handler *Handler,
	followService service.FollowService,
) *FollowHandler {
	return &FollowHandler{
		Handler:       handler,
		followService: followService,
	}
}

// Follow godoc
// @Summary 点赞
// @Schemes
// @Description
// @Tags shop
// @Accept json
// @Produce json
// @Param id path uint64 true "点赞目标的用户ID"
// @Param isFollow path bool true "是否点赞"
// @Success 200 {object} v1.Response
// @Router /follow/{id}/{isFollow} [put]
func (h *FollowHandler) Follow(ctx *gin.Context) {
	var req struct {
		FollowUserID uint64 `uri:"id" binding:"required"`
		IsFollow     bool   `uri:"isFollow" binding:"required"`
	}
	if err := ctx.ShouldBindUri(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := h.followService.Follow(ctx, req.FollowUserID, req.IsFollow); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// IsFollow godoc
// @Summary 查询是否点赞
// @Schemes
// @Description
// @Tags shop
// @Accept json
// @Produce json
// @Param id path uint64 true "点赞目标的用户ID"
// @Success 200 {object} v1.IsFollowResp
// @Router /follow/or/not/{id} [get]
func (h *FollowHandler) IsFollow(ctx *gin.Context) {
	var req struct {
		FollowUserID uint64 `uri:"id" binding:"required"`
	}
	if err := ctx.ShouldBindUri(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	isFollow, err := h.followService.IsFollow(ctx, req.FollowUserID)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	v1.HandleSuccess(ctx, isFollow)
}

// FollowCommons godoc
// @Summary 查询是否点赞
// @Schemes
// @Description
// @Tags shop
// @Accept json
// @Produce json
// @Param id path uint64 true "ID"
// @Success 200 {object} v1.FollowCommonsResp
// @Router /follow/common/{id} [get]
func (h *FollowHandler) FollowCommons(ctx *gin.Context) {
	var req struct {
		ID uint64 `uri:"id" binding:"required"`
	}
	if err := ctx.ShouldBindUri(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	users, err := h.followService.FollowCommons(ctx, req.ID)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	v1.HandleListSuccess(ctx, users, len(users))
}
