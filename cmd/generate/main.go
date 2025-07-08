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

func snakeToLowerCamel(snake string) string {
	parts := strings.Split(snake, "_")
	for i := 1; i < len(parts); i++ {
		if len(parts[i]) > 0 {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}

type Querier interface {
	// GetByID
	// SELECT * FROM @@table WHERE id=@id
	GetByID(id uint64) (*gen.T, error)
}

func main() {
	workdir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	modelPath := filepath.Join(workdir, "internal", "model")
	queryPath := filepath.Join(workdir, "internal", "query")

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
		gen.FieldGORMTag("create_time", func(tag field.GormTag) field.GormTag {
			tag.Set("autoCreateTime", "")
			return tag
		}),
		gen.FieldGORMTag("update_time", func(tag field.GormTag) field.GormTag {
			tag.Set("autoUpdateTime", "")
			return tag
		}),
		gen.FieldJSONTagWithNS(func(columnName string) (tagContent string) {
			return snakeToLowerCamel(columnName)
		}),
		gen.FieldJSONTag("create_time", "-"),
		gen.FieldJSONTag("update_time", "-"),
	}

	g.WithOpts(fieldOpts...)

	g.WithFileNameStrategy(func(table string) string {
		return strings.TrimPrefix(table, "tb_")
	})

	g.WithDataTypeMap(map[string]func(gorm.ColumnType) string{
		"tinyint": func(column gorm.ColumnType) string {
			return "int8"
		},
	})

	dia := mysql.Open("root:123@tcp(127.0.0.1:3306)/hmdp?charset=utf8mb4&parseTime=True&loc=Local")
	db, err := gorm.Open(dia)
	if err != nil {
		fmt.Println(err)
		return
	}
	g.UseDB(db)
	g.ApplyInterface(func(Querier) {},
		g.GenerateModelAs("tb_blog", "Blog",
			gen.FieldNew("Icon", "*string", field.Tag{"gorm": "-", "json": "icon"}),
			gen.FieldNew("Name", "*string", field.Tag{"gorm": "-", "json": "name"}),
			gen.FieldNew("IsLike", "bool", field.Tag{"gorm": "-", "json": "isLike"}),
		),
		g.GenerateModelAs("tb_blog_comments", "BlogComments"),
		g.GenerateModelAs("tb_follow", "Follow"),
		g.GenerateModelAs("tb_seckill_voucher", "SeckillVoucher"),
		g.GenerateModelAs("tb_shop", "Shop"),
		g.GenerateModelAs("tb_shop_type", "ShopType"),
		g.GenerateModelAs("tb_sign", "Sign"),
		g.GenerateModelAs("tb_user", "User"),
		g.GenerateModelAs("tb_user_info", "UserInfo"),
		g.GenerateModelAs("tb_voucher", "Voucher",
			gen.FieldNew("Stock", "int", field.Tag{"gorm": "-", "json": "stock"}),
			gen.FieldNew("BeginTime", "time.Time", field.Tag{"gorm": "-", "json": "beginTime"}),
			gen.FieldNew("EndTime", "time.Time", field.Tag{"gorm": "-", "json": "endTime"}),
		),
		g.GenerateModelAs("tb_voucher_order", "VoucherOrder"),
	)
	g.Execute()
}
