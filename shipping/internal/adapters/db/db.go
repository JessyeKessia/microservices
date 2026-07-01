package db

import (
	"context"
	"fmt"

	"github.com/jessyekessia/microservices/shipping/internal/application/core/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ShippingItem struct {
	gorm.Model
	ShippingID  uint
	ProductCode string
	Quantity    int32
}

type Shipping struct {
	gorm.Model
	OrderID      int64
	DeliveryDays int32
	Items        []ShippingItem
}

type Adapter struct {
	db *gorm.DB
}

func NewAdapter(dataSourceUrl string) (*Adapter, error) {
	db, openErr := gorm.Open(mysql.Open(dataSourceUrl), &gorm.Config{})
	if openErr != nil {
		return nil, fmt.Errorf("db connection error: %v", openErr)
	}

	err := db.AutoMigrate(&Shipping{}, &ShippingItem{})
	if err != nil {
		return nil, fmt.Errorf("db migration error: %v", err)
	}

	return &Adapter{db: db}, nil
}

func (a Adapter) Save(ctx context.Context, shipping *domain.Shipping) error {
	var items []ShippingItem
	for _, item := range shipping.Items {
		items = append(items, ShippingItem{
			ProductCode: item.ProductCode,
			Quantity:    item.Quantity,
		})
	}

	model := Shipping{
		OrderID:      shipping.OrderID,
		DeliveryDays: shipping.DeliveryDays,
		Items:        items,
	}

	res := a.db.WithContext(ctx).Create(&model)
	if res.Error == nil {
		shipping.ID = int64(model.ID)
	}
	return res.Error
}

func (a Adapter) Get(ctx context.Context, id string) (domain.Shipping, error) {
	var model Shipping
	res := a.db.WithContext(ctx).Preload("Items").First(&model, id)

	var items []domain.ShippingItem
	for _, item := range model.Items {
		items = append(items, domain.ShippingItem{
			ProductCode: item.ProductCode,
			Quantity:    item.Quantity,
		})
	}

	return domain.Shipping{
		ID:           int64(model.ID),
		OrderID:      model.OrderID,
		DeliveryDays: model.DeliveryDays,
		Items:        items,
		CreatedAt:    model.CreatedAt.UnixNano(),
	}, res.Error
}
