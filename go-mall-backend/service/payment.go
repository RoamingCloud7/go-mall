package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	logging "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"mall/pkg/e"
	util "mall/pkg/utils"
	dao2 "mall/repository/db/dao"
	model2 "mall/repository/db/model"
	"mall/serializer"
)

type OrderPay struct {
	OrderId   uint    `form:"order_id" json:"order_id"`
	Money     float64 `form:"money" json:"money"`
	OrderNo   string  `form:"orderNo" json:"orderNo"`
	ProductID int     `form:"product_id" json:"product_id"`
	PayTime   string  `form:"payTime" json:"payTime" `
	Sign      string  `form:"sign" json:"sign" `
	BossID    int     `form:"boss_id" json:"boss_id"`
	BossName  string  `form:"boss_name" json:"boss_name"`
	Num       int     `form:"num" json:"num"`
	Key       string  `form:"key" json:"key"`
}

func (service *OrderPay) PayDown(ctx context.Context, uId uint) serializer.Response {
	code := e.SUCCESS

	err := dao2.NewOrderDao(ctx).Transaction(func(tx *gorm.DB) error {
		util.Encrypt.SetKey(service.Key)
		orderDao := dao2.NewOrderDaoByDB(tx)
		// Retrieve the order by ID
		order, err := orderDao.GetOrderById(service.OrderId)
		if err != nil {
			logging.Info(err)
			return err
		}
		money := order.Money
		num := order.Num
		money = money * float64(num)

		userDao := dao2.NewUserDaoByDB(tx)
		user, err := userDao.GetUserById(uId)
		if err != nil {
			logging.Info(err)
			code = e.ErrorDatabase
			return err
		}

		// Decrypt user money, subtract the order amount, and re-encrypt
		moneyStr := util.Encrypt.AesDecoding(user.Money)
		moneyFloat, _ := strconv.ParseFloat(moneyStr, 64)
		if moneyFloat-money < 0.0 { // Insufficient funds, rollback
			logging.Info(err)
			code = e.ErrorDatabase
			return errors.New("金币不足")
		}

		finMoney := fmt.Sprintf("%f", moneyFloat-money)
		user.Money = util.Encrypt.AesEncoding(finMoney)

		// Update user balance
		err = userDao.UpdateUserById(uId, user)
		if err != nil {
			logging.Info(err)
			code = e.ErrorDatabase
			return err
		}
		boss := new(model2.User)
		boss, err = userDao.GetUserById(uint(service.BossID))
		moneyStr = util.Encrypt.AesDecoding(boss.Money)
		moneyFloat, _ = strconv.ParseFloat(moneyStr, 64)
		finMoney = fmt.Sprintf("%f", moneyFloat+money)
		boss.Money = util.Encrypt.AesEncoding(finMoney)

		// Update boss balance
		err = userDao.UpdateUserById(uint(service.BossID), boss)
		if err != nil {
			logging.Info(err)
			code = e.ErrorDatabase
			return err
		}

		// Update product quantity
		product := new(model2.Product)
		productDao := dao2.NewProductDaoByDB(tx)
		product, err = productDao.GetProductById(uint(service.ProductID))
		if err != nil {
			return err
		}
		product.Num -= num
		err = productDao.UpdateProduct(uint(service.ProductID), product)
		if err != nil {
			logging.Info(err)
			code = e.ErrorDatabase
			return err
		}

		// Update order status to 'paid'
		order.Type = 2
		err = orderDao.UpdateOrderById(service.OrderId, order)
		if err != nil {
			logging.Info(err)
			code = e.ErrorDatabase
			return err
		}

		productUser := model2.Product{
			Name:          product.Name,
			CategoryID:    product.CategoryID,
			Title:         product.Title,
			Info:          product.Info,
			ImgPath:       product.ImgPath,
			Price:         product.Price,
			DiscountPrice: product.DiscountPrice,
			Num:           num,
			OnSale:        false,
			BossID:        uId,
			BossName:      user.UserName,
			BossAvatar:    user.Avatar,
		}

		// Create a new product record for the user
		err = productDao.CreateProduct(&productUser)
		if err != nil {
			logging.Info(err)
			code = e.ErrorDatabase
			return err
		}

		return nil

	})

	if err != nil {
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}
