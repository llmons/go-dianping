package handler

import (
	"github.com/gin-gonic/gin"
	"go-dianping/api/v1"
	"go-dianping/internal/model"
	"go-dianping/internal/service"
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

// QueryShopById godoc
// @Summary 根据 id 查询商铺信息
// @Schemes
// @Description
// @Tags shop
// @Produce json
// @Param id path string true "商铺 id"
// @Success 200 {object} v1.QueryShopByIDResp
// @Router /shop/{id} [get]
func (h *ShopHandler) QueryShopById(ctx *gin.Context) {
	var req v1.QueryShopByIDReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	data, err := h.shopService.QueryById(ctx.Request.Context(), &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}
	v1.HandleSuccess(ctx, data)
}

// UpdateShop godoc
// @Summary 更新商铺信息
// @Schemes
// @Description
// @Tags shop
// @Accept json
// @Produce json
// @Param request body model.Shop true "商铺数据"
// @Success 200 {object} v1.Response
// @Router /api/shop [put]
func (h *ShopHandler) UpdateShop(ctx *gin.Context) {
	var req model.Shop
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err := h.shopService.UpdateShop(ctx.Request.Context(), &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// QueryShopOfType godoc
// @Summary 根据商铺类型分页查询商铺信息
// @Schemes
// @Description
// @Tags shop
// @Accept json
// @Produce json
// @Param typeId query uint64 true "商铺类型 ID"
// @Param current query int true "页码"
// @Param x query float64 true "纬度"
// @Param y query float64 true "经度"
// @Success 200 {object} v1.Response
// @Router /shop/of/type [get]
func (h *ShopHandler) QueryShopOfType(ctx *gin.Context) {
	var req struct {
		TypeId  uint64   `form:"typeId" binding:"required"`
		Current int      `form:"current,default=1" binding:"required"`
		X       *float64 `form:"x"`
		Y       *float64 `form:"y"`
	}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err := h.shopService.QueryShopOfType(ctx.Request.Context(), req.TypeId, req.Current, req.X, req.Y)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}
