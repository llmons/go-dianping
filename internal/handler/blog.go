package handler

import (
	"github.com/gin-gonic/gin"
	v1 "go-dianping/api/v1"
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
// @Success 200 {object} v1.Response
// @Router /blog/ [post]
func (h *BlogHandler) SaveBlog(ctx *gin.Context) {
	var blog model.Blog
	if err := ctx.ShouldBindJSON(&blog); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}
	id, err := h.blogService.SaveBlog(ctx, &blog)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	v1.HandleSuccess(ctx, id)
}

// LikeBlog godoc
// @Summary 点赞博文
// @Schemes
// @Description
// @Tags shop
// @Accept json
// @Produce json
// @Param id path uint64 true "博文 FollowUserId"
// @Success 200 {object} v1.Response
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
	v1.HandleSuccess(ctx, nil)
}

// QueryMyBlog godoc
// @Summary 查询我的博文
// @Schemes
// @Description
// @Tags shop
// @Accept json
// @Produce json
// @Param current query int true "当前页码"
// @Success 200 {object} v1.QueryMyBlogResp
// @Router /blog/of/me [get]
func (h *BlogHandler) QueryMyBlog(ctx *gin.Context) {
	var req struct {
		Current int `form:"id" binding:"required"`
	}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	records, err := h.blogService.QueryMyBlog(ctx, req.Current)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	v1.HandleListSuccess(ctx, records, len(records))
}

// QueryHotBlog godoc
// @Summary 查询热门博文
// @Schemes
// @Description
// @Tags shop
// @Accept json
// @Produce json
// @Param current query int true "当前页码"
// @Success 200 {object} v1.QueryHotBlogResp
// @Router /blog/hot [get]
func (h *BlogHandler) QueryHotBlog(ctx *gin.Context) {
	var req struct {
		Current int `form:"id" binding:"required"`
	}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	records, err := h.blogService.QueryHotBlog(ctx, req.Current)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	v1.HandleListSuccess(ctx, records, len(records))
}

// QueryById godoc
// @Summary 根据 FollowUserId 查询博文
// @Schemes
// @Description
// @Tags shop
// @Accept json
// @Produce json
// @Param id path uint64 true "博文 FollowUserId"
// @Success 200 {object} v1.QueryBlogByIDResp
// @Router /blog/{id} [get]
func (h *BlogHandler) QueryById(ctx *gin.Context) {
	var req struct {
		ID uint64 `uri:"id" binding:"required"`
	}
	if err := ctx.ShouldBindUri(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	blog, err := h.blogService.QueryBlogById(ctx, req.ID)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	v1.HandleSuccess(ctx, blog)
}

// QueryBlogLikes godoc
// @Summary 根据博文的点赞信息
// @Schemes
// @Description
// @Tags shop
// @Accept json
// @Produce json
// @Param id path uint64 true "博文 FollowUserId"
// @Success 200 {object} v1.QueryBlogByIDResp
// @Router /blog/{id} [get]
func (h *BlogHandler) QueryBlogLikes(ctx *gin.Context) {
	var req struct {
		ID uint64 `uri:"id" binding:"required"`
	}
	if err := ctx.ShouldBindUri(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	blog, err := h.blogService.QueryBlogLikes(ctx, req.ID)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	v1.HandleSuccess(ctx, blog)
}

// QueryBlogByUserID godoc
// @Summary 根据 userID 查询博文
// @Schemes
// @Description
// @Tags shop
// @Accept json
// @Produce json
// @Param id query uint64 true "userID"
// @Success 200 {object} v1.QueryBlogByUserIDResp
// @Router /blog/of/user [get]
func (h *BlogHandler) QueryBlogByUserID(ctx *gin.Context) {
	var req struct {
		ID      uint64 `form:"id" binding:"required"`
		Current int    `form:"current,default=1" binding:"required"`
	}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	blogs, err := h.blogService.QueryBlogByUserID(ctx, req.ID, req.Current)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	v1.HandleListSuccess(ctx, blogs, len(blogs))
}
