package handler

import (
	"github.com/gin-gonic/gin"
	"go-dianping/internal/service"
	"go-dianping/pkg/helper/resp"
	"go.uber.org/zap"
	"net/http"
)

type ShopTypeHandler struct {
	*Handler
	shopTypeService service.ShopTypeService
}

func NewShopTypeHandler(
	handler *Handler,
	shopTypeService service.ShopTypeService,
) *ShopTypeHandler {
	return &ShopTypeHandler{
		Handler:         handler,
		shopTypeService: shopTypeService,
	}
}

func (h *ShopTypeHandler) GetShopTypeList(ctx *gin.Context) {
	shopTypeList, err := h.shopTypeService.GetShopTypeList(ctx)
	h.logger.Info("GetShopTypeList", zap.Any("shopTypeList", shopTypeList))
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	resp.HandleListSuccess(ctx, shopTypeList, len(shopTypeList))
}
