package handler

import (
	"fmt"
	"github.com/samber/lo"
	v1 "go-dianping/api/v1"
	"go-dianping/internal/handler"
	"go-dianping/internal/repository"
	"go-dianping/internal/service"
	"net/http"
	"testing"
	"time"
)

func TestVoucherHandler_AddSeckillVoucher(t *testing.T) {
	layout := "2006-01-02T15:04:05"
	beginTime, err := time.Parse(layout, "2022-01-26T10:09:17")
	if err != nil {
		panic(err)
	}
	endTime, err := time.Parse(layout, "2022-01-26T23:09:04")
	if err != nil {
		panic(err)
	}

	params := v1.AddSeckillVoucherReq{
		ShopID:      lo.ToPtr(uint64(2)),
		Title:       "100元代金券",
		SubTitle:    lo.ToPtr("周一至周五均可使用"),
		Rules:       lo.ToPtr("全场通用\\n无需预约\\n可无限叠加\\不兑现，不找零\\n仅限堂食"),
		PayValue:    8000,
		ActualValue: 10000,
		Type:        1,
		Stock:       100,
		BeginTime:   beginTime,
		EndTime:     endTime,
	}

	voucherRepo := repository.NewVoucherRepository(repo)
	seckillRepo := repository.NewSeckillVoucherRepository(repo)
	voucherService := service.NewVoucherService(srv, voucherRepo, seckillRepo)

	voucherHandler := handler.NewVoucherHandler(hdl, voucherService)
	router.POST("/api/voucher/seckill", voucherHandler.AddSeckillVoucher)

	e := newHttpExcept(t, router)

	fmt.Println(e)

	obj := e.POST("/api/voucher/seckill").
		WithHeader("Content-Type", "application/json").
		WithJSON(params).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()

	fmt.Println(obj)

	obj.Value("success").IsEqual(true)
	obj.Value("data").IsNumber()
}
