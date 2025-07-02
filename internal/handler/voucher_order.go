package handler

import (
	"github.com/gin-gonic/gin"
	"go-dianping/internal/service"
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

func (h *VoucherOrderHandler) GetVoucherOrder(ctx *gin.Context) {

}
