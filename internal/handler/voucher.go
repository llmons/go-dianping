package handler

import (
	"github.com/gin-gonic/gin"
	v1 "go-dianping/api/v1"
	"go-dianping/internal/model"
	"go-dianping/internal/service"
	"net/http"
)

type VoucherHandler struct {
	*Handler
	voucherService service.VoucherService
}

func NewVoucherHandler(
	handler *Handler,
	voucherService service.VoucherService,
) *VoucherHandler {
	return &VoucherHandler{
		Handler:        handler,
		voucherService: voucherService,
	}
}

// AddSeckillVoucher godoc
// @Summary 新增秒杀券
// @Schemes
// @Description
// @Tags voucher
// @Accept json
// @Produce json
// @Param request body model.Voucher true "优惠券信息，包含秒杀信息"
// @Success 200 {object} v1.AddSeckillVoucherResp
// @Router /voucher/seckill [post]
func (h *VoucherHandler) AddSeckillVoucher(ctx *gin.Context) {
	var req model.Voucher
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := h.voucherService.AddSeckillVoucher(ctx.Request.Context(), &req); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	v1.HandleSuccess(ctx, req.ID)
}

// AddVoucher godoc
// @Summary 新增普通券
// @Schemes
// @Description
// @Tags voucher
// @Accept json
// @Produce json
// @Param request body model.Voucher true "优惠券信息"
// @Success 200 {object} v1.AddVoucherResp
// @Router /voucher [post]
func (h *VoucherHandler) AddVoucher(ctx *gin.Context) {
	var req model.Voucher
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := h.voucherService.AddVoucher(ctx.Request.Context(), &req); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	v1.HandleSuccess(ctx, req.ID)
}

// QueryVoucherOfShop godoc
// @Summary 查询店铺的优惠券列表
// @Schemes
// @Description
// @Tags voucher
// @Accept json
// @Produce json
// @Param shopId path uint64 true "商铺 id"
// @Success 200 {object} v1.QueryVoucherOfShopResp
// @Router /voucher/list/{shopId} [get]
func (h *VoucherHandler) QueryVoucherOfShop(ctx *gin.Context) {
	var req v1.QueryVoucherOfShopReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}
	list, err := h.voucherService.QueryVoucherOfShop(ctx.Request.Context(), &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	v1.HandleListSuccess(ctx, list, len(list))
}
