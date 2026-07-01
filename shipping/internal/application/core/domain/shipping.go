package domain

import "time"

type ShippingItem struct {
	ProductCode string `json:"product_code"`
	Quantity    int32  `json:"quantity"`
}

type Shipping struct {
	ID           int64          `json:"id"`
	OrderID      int64          `json:"order_id"`
	DeliveryDays int32          `json:"delivery_days"`
	Items        []ShippingItem `json:"items"`
	CreatedAt    int64          `json:"created_at"`
}

func NewShipping(orderId int64, items []ShippingItem) Shipping {
	deliveryDays := calculateDeliveryDays(items)
	return Shipping{
		OrderID:      orderId,
		DeliveryDays: deliveryDays,
		Items:        items,
		CreatedAt:    time.Now().Unix(),
	}
}

func calculateDeliveryDays(items []ShippingItem) int32 {
	var totalQuantity int32
	for _, item := range items {
		totalQuantity += item.Quantity
	}
	deliveryDays := int32(1) + (totalQuantity / 5)
	return deliveryDays
}
