package handler

import (
	"github.com/gin-gonic/gin"
	"go-dianping/internal/service"
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

func (h *BlogHandler) GetBlog(ctx *gin.Context) {

}
