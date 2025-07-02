package service

import (
	"go-dianping/internal/repository"
)

type VoucherOrderService interface {
}

func NewVoucherOrderService(
	service *Service,
	voucherOrderRepository repository.VoucherOrderRepository,
) VoucherOrderService {
	return &voucherOrderService{
		Service:                service,
		voucherOrderRepository: voucherOrderRepository,
	}
}

type voucherOrderService struct {
	*Service
	voucherOrderRepository repository.VoucherOrderRepository
}
