package handler

import (
	"github.com/gin-gonic/gin"
	"go-dianping/api/v1"
	"go-dianping/internal/service"
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
// @Produce json
// @Success 200 {object} v1.GetShopTypeListResp
// @Router /shop-type/list [get]
func (h *ShopTypeHandler) GetShopTypeList(ctx *gin.Context) {
	shopTypeList, err := h.shopTypeService.GetShopTypeList(ctx.Request.Context())
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	v1.HandleListSuccess(ctx, shopTypeList, len(shopTypeList))
}
