package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	workdir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}
	modelPath := filepath.Join(workdir, "internal/entity")
	queryPath := filepath.Join(workdir, "internal/query")

	g := gen.NewGenerator(gen.Config{
		ModelPkgPath:     modelPath,
		OutPath:          queryPath,
		Mode:             gen.WithoutContext | gen.WithQueryInterface,
		FieldNullable:    true,
		FieldCoverable:   true,
		FieldSignable:    true,
		FieldWithTypeTag: true,
	})

	fieldOpts := []gen.ModelOpt{
		gen.FieldGORMTag("update_time", func(tag field.GormTag) field.GormTag {
			tag.Set("autoUpdateTime", "")
			return tag
		}),
		gen.FieldGORMTag("create_time", func(tag field.GormTag) field.GormTag {
			tag.Set("autoCreateTime", "")
			return tag
		}),
	}

	g.WithOpts(fieldOpts...)

	g.WithFileNameStrategy(func(table string) string {
		return strings.TrimPrefix(table, "tb_")
	})

	dia := mysql.Open("root:123@tcp(127.0.0.1:3306)/hmdp?charset=utf8mb4&parseTime=True&loc=Local")
	db, err := gorm.Open(dia)
	if err != nil {
		fmt.Println(err)
		return
	}
	g.UseDB(db)
	g.ApplyBasic(
		g.GenerateModelAs("tb_blog", "Blog"),
		g.GenerateModelAs("tb_blog_comments", "BlogComments"),
		g.GenerateModelAs("tb_follow", "Follow"),
		g.GenerateModelAs("tb_seckill_voucher", "SeckillVoucher"),
		g.GenerateModelAs("tb_shop", "Shop"),
		g.GenerateModelAs("tb_shop_type", "ShopType"),
		g.GenerateModelAs("tb_sign", "Sign"),
		g.GenerateModelAs("tb_user", "User"),
		g.GenerateModelAs("tb_user_info", "UserInfo"),
		g.GenerateModelAs("tb_voucher", "Voucher"),
		g.GenerateModelAs("tb_voucher_order", "VoucherOrder"),
	)
	g.Execute()
}
