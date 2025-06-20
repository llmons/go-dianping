package handler

import (
	"github.com/gin-gonic/gin"
	"go-dianping/api"
	"go-dianping/internal/service"
	"go.uber.org/zap"
	"net/http"
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

// GetShopById godoc
// @Summary 根据 id 获取商铺
// @Schemes
// @Description
// @Tags shop
// @Produce json
// @Params id path string true "id"
// @Success 200 {object} api.GetShopByIdResp
// @Router /shop/:id [get]
func (h *ShopHandler) GetShopById(ctx *gin.Context) {
	var params api.GetShopByIdReq
	if err := ctx.ShouldBindUri(&params); err != nil {
		api.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	data, err := h.shopService.GetShopById(ctx.Request.Context(), &params)
	h.logger.Info("GetShopById", zap.Any("data", data))
	if err != nil {
		api.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}
	api.HandleSuccess(ctx, data)
}

// UpdateShop godoc
// @Summary 更新商铺
// @Schemes
// @Description
// @Tags shop
// @Produce json
// @Params request body api.UpdateShopReq true "商铺信息"
// @Success 200 {object} api.UpdateShopResp
// @Router /shop [put]
func (h *ShopHandler) UpdateShop(ctx *gin.Context) {
	var req api.UpdateShopReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		api.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err := h.shopService.UpdateShop(ctx.Request.Context(), &req)
	h.logger.Info("UpdateShop", zap.Any("update data", req))
	if err != nil {
		api.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}
	api.HandleSuccess(ctx, nil)
}
