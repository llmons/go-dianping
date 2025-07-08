package handler

import (
	"github.com/gin-gonic/gin"
	"go-dianping/internal/service"
)

type BlogCommentsHandler struct {
	*Handler
	blogCommentsService service.BlogCommentsService
}

func NewBlogCommentsHandler(
	handler *Handler,
	blogCommentsService service.BlogCommentsService,
) *BlogCommentsHandler {
	return &BlogCommentsHandler{
		Handler:             handler,
		blogCommentsService: blogCommentsService,
	}
}

func (h *BlogCommentsHandler) GetBlogComments(ctx *gin.Context) {

}
