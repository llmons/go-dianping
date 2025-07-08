package handler

import (
	"github.com/gin-gonic/gin"
	v1 "go-dianping/api/v1"
	"go-dianping/internal/base/user_holder"
	"go-dianping/internal/model"
	"go-dianping/internal/service"
	"net/http"
)

type BlogHandler struct {
	*Handler
	blogService service.BlogService
}

func NewBlogHandler(
	handler *Handler,
	blogService service.BlogService,
) *BlogHandler {
	return &BlogHandler{
		Handler:     handler,
		blogService: blogService,
	}
}

// SaveBlog godoc
// @Summary 保存博文
// @Schemes
// @Description
// @Tags shop
// @Accept json
// @Produce json
// @Param blog body model.Blog true "博文信息"
// @Success 200 {object} v1.QueryShopByIDResp
// @Router /blog/ [post]
func (h *BlogHandler) SaveBlog(ctx *gin.Context) {
	var blog model.Blog
	if err := ctx.ShouldBindJSON(&blog); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}
	// 获取登录用户
	user := user_holder.GetUser(ctx)
	blog.ID = *user.ID
	//	保存探店博文
	if err := h.blogService.SaveBlog(ctx, &blog); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	//	返回 id
	v1.HandleSuccess(ctx, blog.ID)
}

// LikeBlog godoc
// @Summary 点赞博文
// @Schemes
// @Description
// @Tags shop
// @Accept json
// @Produce json
// @Param id path uint64 true "博文 ID"
// @Success 200 {object} v1.QueryShopByIDResp
// @Router /blog/like/{id} [put]
func (h *BlogHandler) LikeBlog(ctx *gin.Context) {
	var req struct {
		ID uint64 `uri:"id" binding:"required"`
	}
	if err := ctx.ShouldBindUri(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := h.blogService.LikeBlog(ctx, req.ID); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	//	返回 id
	v1.HandleSuccess(ctx, nil)
}
