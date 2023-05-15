package shop

import (
	"errors"
	"fmt"
	"fresh-shop/server/global"
	"fresh-shop/server/model/common/request"
	"fresh-shop/server/model/shop"
	shopReq "fresh-shop/server/model/shop/request"
	systemReq "fresh-shop/server/model/system/request"
	"fresh-shop/server/service/common"
	"fresh-shop/server/utils"
	"gorm.io/gorm"
	"strconv"
)

type OrderService struct {
}

// CreateOrder 创建Order记录
// Author [piexlmax](https://github.com/likfees)
func (orderService *OrderService) CreateOrder(order shop.Order, userClaims *systemReq.CustomClaims) (err error) {
	// 获取收货地址信息
	var address shop.UserAddress
	if errors.Is(global.DB.Where("id = ?", order.AddressId).First(&address).Error, gorm.ErrRecordNotFound) {
		return errors.New("收货地址不存在")
	}
	// 获取购物车列表数据
	var cartList []shop.Cart
	global.DB.Where("user_id = ? and id in ?", order.UserId, order.CartIds).Preload("Goods.Images").Find(&cartList)
	if len(cartList) <= 0 {
		return errors.New("商品查询失败")
	}
	var orderDetailList []shop.OrderDetails
	pointCfg, err := common.GetSysConfig("point")
	pointSwitch := true
	if err != nil {
		if errors.Is(err, common.ErrConfigDisabled) {
			pointSwitch = false
		} else {
			global.SugarLog.Errorf("创建订单时查询积分配置参数异常, err:%v \n", err)
			return err
		}
	}
	// 判断库存是否充足  以后可以上锁，解决高并发
	for _, c := range cartList {
		// 购物车数量大于库存
		if *c.Num > *c.Goods.Store {
			return errors.New("商品库存不足")
		}
		// 计算总数量
		*order.Num = *order.Num + *c.Num

		// 计算总金额 如果优惠价小于成本价
		if *c.Goods.Price < *c.Goods.CostPrice {
			*order.Total += float64(*c.Num) * *c.Goods.Price
		} else {
			*order.Total += float64(*c.Num) * *c.Goods.CostPrice
		}
		// 组织订单详情数据
		imgUrl := ""
		if len(c.Goods.Images) > 0 {
			imgUrl = c.Goods.Images[0].Url
		}
		orderDetail := shop.OrderDetails{}
		orderDetail.GoodsId = c.Goods.ID
		orderDetail.GoodsName = c.Goods.Name
		orderDetail.GoodsImage = imgUrl
		orderDetail.Unit = c.Goods.Unit
		orderDetail.Num = c.Num
		orderDetail.Price = c.Goods.Price
		// 计算单个商品多个数量的总金额
		if *c.Goods.Price < *c.Goods.CostPrice {
			*orderDetail.Total = float64(*c.Num) * *c.Goods.Price
		} else {
			*orderDetail.Total = float64(*c.Num) * *c.Goods.CostPrice
		}
		// 规格id 现在只开发了单规格订单，多规格以后在支持
		orderDetail.SpecId = utils.Pointer(0)
		orderDetail.SpecKeyName = ""

		// 计算赠送积分
		if pointSwitch {
			point, err := strconv.Atoi(pointCfg)
			if err != nil {
				global.SugarLog.Errorf("创建订单详情信息时转换积分配置参数异常, err:%v \n", err)
				return err
			}
			*orderDetail.GiftPoints = *orderDetail.Total * (float64(point) / 100)
		}
		orderDetailList = append(orderDetailList, orderDetail)
	}

	// 设置订单基本信息
	order.OrderSn = utils.GenerateOrderNumber("SN")
	// 商品区域默认为普通商品
	if order.GoodsArea == nil || *order.GoodsArea == 0 {
		*order.GoodsArea = 0
	}
	order.ShipmentName = address.Name
	order.ShipmentMobile = address.Mobile
	order.ShipmentAddress = address.Address + address.Title + address.Detail
	*order.Status = 0 // 未付款状态
	// 计算总赠送积分
	if pointSwitch {
		point, err := strconv.Atoi(pointCfg)
		if err != nil {
			global.SugarLog.Errorf("创建订单时转换积分配置参数异常, err:%v \n", err)
			return err
		}
		// 公式 总金额 * n%
		*order.GiftPoints = *order.Total * (float64(point) / 100)
	}

	log := fmt.Sprintf("[OrderService] CreateOrder submit data:%+v; \n", order)
	// 启动事务
	txDB := global.DB.Begin()
	// 创建订单
	if err = global.DB.Create(&order).Error; err != nil {
		txDB.Rollback()
		global.SugarLog.Errorf("log:%s,err:%v \n", log, err)
		return errors.New("订单创建失败")
	}
	// 创建订单详情
	// 设置订单详情 orderId
	for _, v := range orderDetailList {
		v.OrderId = order.ID
	}
	if err = global.DB.Create(&orderDetailList).Error; err != nil {
		txDB.Rollback()
		global.SugarLog.Errorf("log:%s,err:%v \n", log, err)
		return errors.New("订单详情创建失败")
	}
	// 删除购物车列表
	if err = global.DB.Delete(&cartList).Error; err != nil {
		txDB.Rollback()
		global.SugarLog.Errorf("log:%s,err:%v \n", log, err)
		return errors.New("购物车删除失败")
	}
	// 提交事务
	txDB.Commit()
	return nil
}

// DeleteOrder 删除Order记录
// Author [piexlmax](https://github.com/likfees)
func (orderService *OrderService) DeleteOrder(order shop.Order) (err error) {
	var detail shop.OrderDetails
	err = global.DB.Transaction(func(tx *gorm.DB) error {
		err = global.DB.Where("order_id = ?", order.ID).Delete(&detail).Error
		if err != nil {
			global.SugarLog.Errorf("删除订单详情失败 %d, err: %v", order.ID, err)
			return err
		}
		err = global.DB.Delete(&order).Error
		if err != nil {
			global.SugarLog.Errorf("删除订单失败 %d, err: %v", order.ID, err)
			return err
		}
		return nil
	})
	return err
}

// DeleteOrderByIds 批量删除Order记录
// Author [piexlmax](https://github.com/likfees)
func (orderService *OrderService) DeleteOrderByIds(ids request.IdsReq) (err error) {
	err = global.DB.Delete(&[]shop.Order{}, "id in ?", ids.Ids).Error
	return err
}

// UpdateOrder 更新Order记录
// Author [piexlmax](https://github.com/likfees)
func (orderService *OrderService) UpdateOrder(order shop.Order) (err error) {
	err = global.DB.Save(&order).Error
	return err
}

// GetOrder 根据id获取Order记录
// Author [piexlmax](https://github.com/likfees)
func (orderService *OrderService) GetOrder(id uint) (order shop.Order, err error) {
	err = global.DB.Where("id = ?", id).First(&order).Error
	return
}

// GetOrderInfoList 分页获取Order记录
// Author [piexlmax](https://github.com/likfees)
func (orderService *OrderService) GetOrderInfoList(info shopReq.OrderSearch) (list []shop.Order, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.DB.Model(&shop.Order{})
	var orders []shop.Order
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.OrderSn != "" {
		db = db.Where("order_sn LIKE ?", "%"+info.OrderSn+"%")
	}
	if info.ShipmentName != "" {
		db = db.Where("shipment_name LIKE ?", "%"+info.ShipmentName+"%")
	}
	if info.ShipmentMobile != "" {
		db = db.Where("shipment_mobile LIKE ?", "%"+info.ShipmentMobile+"%")
	}
	if info.ShipmentAddress != "" {
		db = db.Where("shipment_address LIKE ?", "%"+info.ShipmentAddress+"%")
	}
	if info.Payment != nil {
		db = db.Where("payment = ?", info.Payment)
	}
	if info.Status != nil {
		db = db.Where("status = ?", info.Status)
	}
	if info.StartShipmentTime != nil && info.EndShipmentTime != nil {
		db = db.Where("shipment_time BETWEEN ? AND ? ", info.StartShipmentTime, info.EndShipmentTime)
	}
	if info.StartReceiveTime != nil && info.EndReceiveTime != nil {
		db = db.Where("receive_time BETWEEN ? AND ? ", info.StartReceiveTime, info.EndReceiveTime)
	}
	if info.StartCancelTime != nil && info.EndCancelTime != nil {
		db = db.Where("cancel_time BETWEEN ? AND ? ", info.StartCancelTime, info.EndCancelTime)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Find(&orders).Error
	return orders, total, err
}
