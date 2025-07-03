package service

import (
	"context"
	"github.com/jinzhu/copier"
	v1 "go-dianping/api/v1"
	"go-dianping/internal/model"
	"go-dianping/internal/repository"
)

type VoucherService interface {
	AddSeckillVoucher(ctx context.Context, req *v1.AddSeckillVoucherReq) error
}

func NewVoucherService(
	service *Service,
	voucherRepository repository.VoucherRepository,
	seckillVoucherRepository repository.SeckillVoucherRepository,
) VoucherService {
	return &voucherService{
		Service:            service,
		voucherRepo:        voucherRepository,
		seckillVoucherRepo: seckillVoucherRepository,
	}
}

type voucherService struct {
	*Service
	voucherRepo        repository.VoucherRepository
	seckillVoucherRepo repository.SeckillVoucherRepository
}

func (s *voucherService) AddSeckillVoucher(ctx context.Context, req *v1.AddSeckillVoucherReq) error {
	// 保存优惠券
	if err := s.voucherRepo.Save(ctx, (*model.Voucher)(req)); err != nil {
		return err
	}
	// 保存秒杀信息
	var seckillVoucher model.SeckillVoucher
	seckillVoucher.VoucherID = req.ID
	if err := copier.Copy(&seckillVoucher, &req); err != nil {
		return err
	}
	return s.seckillVoucherRepo.Save(ctx, &seckillVoucher)
}
