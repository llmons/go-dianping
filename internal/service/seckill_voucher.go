package service

import (
	"go-dianping/internal/repository"
)

type SeckillVoucherService interface {
}

func NewSeckillVoucherService(
	service *Service,
	seckillVoucherRepository repository.SeckillVoucherRepository,
) SeckillVoucherService {
	return &seckillVoucherService{
		Service:                  service,
		seckillVoucherRepository: seckillVoucherRepository,
	}
}

type seckillVoucherService struct {
	*Service
	seckillVoucherRepository repository.SeckillVoucherRepository
}
