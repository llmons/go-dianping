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

func (h *ShopHandler) GetShopById(ctx *gin.Context) {
	//var params api.GetShopByIdReq
	//if err := ctx.ShouldBind(&params); err != nil {
	//	api.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
	//	return
	//}

	id := ctx.Param("id")
	var params api.GetShopByIdReq
	params.Id = id
	data, err := h.shopService.GetShopById(ctx.Request.Context(), &params)
	h.logger.Info("GetShopById", zap.Any("data", data))
	if err != nil {
		api.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}
	api.HandleSuccess(ctx, data)
}
