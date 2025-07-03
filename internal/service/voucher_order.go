package service

import (
	"context"
	"go-dianping/api/v1"
	"go-dianping/internal/base/redis_worker"
	"go-dianping/internal/base/user_holder"
	"go-dianping/internal/entity"
	"go-dianping/internal/repository"
	"time"
)

type VoucherOrderService interface {
	SeckillVoucher(ctx context.Context, req *v1.SeckillVoucherReq) (int64, error)
}

func NewVoucherOrderService(
	service *Service,
	voucherOrderRepository repository.VoucherOrderRepository,
	seckillVoucherRepository repository.SeckillVoucherRepository,
	redisWorker redis_worker.RedisWorker,
) VoucherOrderService {
	return &voucherOrderService{
		Service:                  service,
		voucherOrderRepository:   voucherOrderRepository,
		seckillVoucherRepository: seckillVoucherRepository,
		redisWorker:              redisWorker,
	}
}

type voucherOrderService struct {
	*Service
	voucherOrderRepository   repository.VoucherOrderRepository
	seckillVoucherRepository repository.SeckillVoucherRepository
	redisWorker              redis_worker.RedisWorker
}

func (s *voucherOrderService) SeckillVoucher(ctx context.Context, req *v1.SeckillVoucherReq) (int64, error) {
	//	1. 查询优惠券
	voucher, err := s.seckillVoucherRepository.GetByID(ctx, req.ID)
	if err != nil {
		return 0, err
	}
	//	2. 判断秒杀是否开始
	if voucher.BeginTime.After(time.Now()) {
		return 0, v1.ErrSeckillNotStart
	}
	//	3. 判断秒杀是否已经结束
	if voucher.EndTime.Before(time.Now()) {
		return 0, v1.ErrSeckillIsEnd
	}
	//	4.
	if voucher.Stock < 1 {
		return 0, v1.ErrInsufficientStock
	}
	//	5. 扣减库存
	var orderId int64
	// 7. 返回订单 id
	return orderId, s.tm.Transaction(ctx, func(ctx context.Context) error {
		info, err := s.seckillVoucherRepository.DecStock(ctx, voucher.VoucherID)
		if err != nil {
			return err
		}
		if info.Error != nil {
			// 扣减失败
			return v1.ErrInsufficientStock
		}
		//	6. 创建订单
		var voucherOrder entity.VoucherOrder
		// 6.1. 订单 id
		orderId, err = s.redisWorker.NextId(ctx, "order")
		if err != nil {
			return err
		}
		voucherOrder.ID = orderId
		// 6.2. 用户 id
		userId := user_holder.GetUser(ctx).ID
		voucherOrder.UserID = *userId
		// 6.3. 代金券 id
		voucherOrder.VoucherID = voucher.VoucherID
		if err := s.voucherOrderRepository.Save(ctx, &voucherOrder); err != nil {
			return err
		}

		return nil
	})
}
