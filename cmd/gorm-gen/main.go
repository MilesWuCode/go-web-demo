package main

import (
	"log"
	"web-demo/config"
	"web-demo/database"
	"web-demo/model"

	"gorm.io/gen"
)

// Dynamic SQL
type Querier interface {
}

func main() {
	// 設定值
	cfg := config.Get()

	// 建立 GORM Gen 生成器
	g := gen.NewGenerator(gen.Config{
		OutPath:      "./generated/query",
		ModelPkgPath: "./model",
		Mode:         gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,

		// 預設值全為false
		FieldNullable:     true, // 允許指標類型映射 NULL
		FieldCoverable:    true, // 生成可覆蓋的欄位邏輯
		FieldSignable:     true, // 支持 unsigned 類型
		FieldWithIndexTag: true, // 保留索引標籤
		FieldWithTypeTag:  true, // 保留資料庫具體類型標籤
	})

	// 建立資料庫連線
	db, err := database.NewDB(cfg)
	if err != nil {
		log.Fatalf("無法連接資料庫: %v", err)
	}
	// 使用 defer 確保程式結束時關閉資料庫
	defer database.CloseDB(db)

	// 注入資料庫連線
	g.UseDB(db)

	// Generate basic type-safe DAO API for struct `model.User` following conventions
	g.ApplyBasic(model.User{})

	// Generate Type Safe API with Dynamic SQL defined on Querier interface for `model.User` and `model.Company`
	g.ApplyInterface(func(Querier) {}, model.User{}, model.Post{})

	// Generate the code
	g.Execute()
}
