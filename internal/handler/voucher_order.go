package handler

import (
	"github.com/gin-gonic/gin"
	v1 "go-dianping/api/v1"
	"go-dianping/internal/service"
	"net/http"
)

type VoucherOrderHandler struct {
	*Handler
	voucherOrderService service.VoucherOrderService
}

func NewVoucherOrderHandler(
	handler *Handler,
	voucherOrderService service.VoucherOrderService,
) *VoucherOrderHandler {
	return &VoucherOrderHandler{
		Handler:             handler,
		voucherOrderService: voucherOrderService,
	}
}

func (h *VoucherOrderHandler) SeckillVoucher(ctx *gin.Context) {
	var req v1.SeckillVoucherReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	orderId, err := h.voucherOrderService.SeckillVoucher(ctx.Request.Context(), &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	v1.HandleSuccess(ctx, orderId)
}
