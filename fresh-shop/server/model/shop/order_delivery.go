package shop

import (
	"fresh-shop/server/global"
	"fresh-shop/server/model/business"
	"time"
)

// OrderDelivery 结构体
type OrderDelivery struct {
	global.DbModel
	OrderId       *int                  `json:"orderId" form:"orderId" gorm:"column:order_id;comment:订单Id;size:20;"`
	ScheduledTime time.Time             `json:"scheduledTime" form:"scheduledTime" gorm:"column:scheduled_time;comment:预计到达时间;"`
	DeliverName   string                `json:"deliverName" form:"deliverName" gorm:"column:deliver_name;comment:送货人姓名;"`
	DeliveryId    *int                  `json:"deliveryId" form:"deliveryId" gorm:"column:delivery_id;comment:送货人ID;size:11;"`
	DeliverMobile string                `json:"deliverMobile" form:"deliverMobile" gorm:"column:deliver_mobile;comment:送货人联系电话;size:11;"`
	ReceiptTime   *time.Time            `json:"receiptTime" form:"receiptTime" gorm:"column:receipt_time;comment:收货时间;"`
	UserDelivery  business.UserDelivery `json:"user" gorm:"foreignKey:id;references:delivery_id"`
}

// TableName OrderDelivery 表名
func (OrderDelivery) TableName() string {
	return "shop_order_delivery"
}
