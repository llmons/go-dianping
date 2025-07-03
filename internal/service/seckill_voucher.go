package service

type SeckillVoucherService interface {
}

func NewSeckillVoucherService(
	service *Service,
) SeckillVoucherService {
	return &seckillVoucherService{
		Service: service,
	}
}

type seckillVoucherService struct {
	*Service
}
