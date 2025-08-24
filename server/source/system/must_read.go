package system

import (
	"context"

	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"gorm.io/gorm"
)

const initOrderMustRead = initOrderAuthority + 10

type initMustRead struct{}

// auto run
func init() {
	system.RegisterInit(initOrderMustRead, &initMustRead{})
}

func (i *initMustRead) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(&sysModel.SysMustRead{})
}

func (i *initMustRead) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&sysModel.SysMustRead{})
}

func (i *initMustRead) InitializerName() string {
	return sysModel.SysMustRead{}.TableName()
}

func (i *initMustRead) InitializeData(ctx context.Context) (next context.Context, err error) {
	// 这里可以添加一些初始的必读内容数据（可选）
	// 目前暂时不添加初始数据，只创建表结构
	return ctx, nil
}

func (i *initMustRead) DataInserted(ctx context.Context) bool {
	// 由于没有初始数据，直接返回true
	return true
}
