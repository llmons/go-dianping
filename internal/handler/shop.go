package handler

import (
	"github.com/gin-gonic/gin"
	"go-dianping/internal/service"
)

type ShopHandler struct {
	*Handler
	shopService service.ShopService
}

func NewShopHandler(
	handler *Handler,
	shopService service.ShopService,
) *ShopHandler {
	return &ShopHandler{
		Handler:     handler,
		shopService: shopService,
	}
}

func (h *ShopHandler) GetShopById(ctx *gin.Context) {

}
