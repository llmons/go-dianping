package handler

import (
	"github.com/gin-gonic/gin"
	"go-dianping/api"
	"go-dianping/internal/service"
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

// GetShopTypeList godoc
// @Summary 获取商品类别列表
// @Schemes
// @Description
// @Tags shop-type
// @Accept json
// @Produce json
// @Success 200 {object} api.GetShopTypeListRespData
// @Router /shop-type/list [get]
func (h *ShopTypeHandler) GetShopTypeList(ctx *gin.Context) {
	shopTypeList, err := h.shopTypeService.GetShopTypeList(ctx.Request.Context())
	h.logger.Info("GetShopTypeList", zap.Int("shopTypeListLength", len(shopTypeList)))
	if err != nil {
		api.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	api.HandleListSuccess(ctx, shopTypeList, len(shopTypeList))
}
