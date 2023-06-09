package account

import (
	"fresh-shop/server/global"
)

// AccountGroup 结构体
type AccountGroup struct {
	global.DbModel
	NameEn string   `json:"nameEn" form:"nameEn" gorm:"column:name_en;comment:币中英文名;size:20;"`
	NameCn string   `json:"nameCn" form:"nameCn" gorm:"column:name_cn;comment:币种中文名;size:20;"`
	Places *float64 `json:"places" form:"places" gorm:"column:places;default:4;comment:小数点位数;"`
	Status *int     `json:"status" form:"status" gorm:"column:status;default:1;comment:状态(0禁用 1启用);"`
	Sync   *int     `json:"sync" form:"sync" gorm:"column:sync;default:0;comment:同步状态 生成用户对应的币种账户和币种流水表;"`
}

// TableName AccountGroup 表名
func (AccountGroup) TableName() string {
	return "user_account_group"
}
