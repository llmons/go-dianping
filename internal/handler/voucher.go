package handler

import (
	"github.com/gin-gonic/gin"
	v1 "go-dianping/api/v1"
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

func (h *VoucherHandler) AddVoucher(ctx *gin.Context) {}

func (h *VoucherHandler) AddSeckillVoucher(ctx *gin.Context) {
	var req v1.AddSeckillVoucherReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := h.voucherService.AddSeckillVoucher(ctx, &req); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	v1.HandleSuccess(ctx, req.ID)
}
