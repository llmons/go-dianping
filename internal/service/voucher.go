package service

import (
	"context"
	"github.com/jinzhu/copier"
	v1 "go-dianping/api/v1"
	"go-dianping/internal/model"
)

type VoucherService interface {
	AddSeckillVoucher(ctx context.Context, req *v1.AddSeckillVoucherReq) error
}

func NewVoucherService(
	service *Service,
) VoucherService {
	return &voucherService{
		Service: service,
	}
}

type voucherService struct {
	*Service
}

func (s *voucherService) AddSeckillVoucher(_ context.Context, req *v1.AddSeckillVoucherReq) error {
	// 保存优惠券
	if err := s.query.Voucher.Save((*model.Voucher)(req)); err != nil {
		return err
	}
	// 保存秒杀信息
	var seckillVoucher model.SeckillVoucher
	seckillVoucher.VoucherID = req.ID
	if err := copier.Copy(&seckillVoucher, &req); err != nil {
		return err
	}
	return s.query.SeckillVoucher.Save(&seckillVoucher)
}
