package grpc

import (
	"context"
	"fmt"
	"log"

	"github.com/jessyekessia/microservices-proto/golang/shipping"
	"github.com/jessyekessia/microservices/shipping/internal/application/core/domain"
)

func (a Adapter) CalculateShipping(
	ctx context.Context,
	request *shipping.CreateShippingRequest,
) (*shipping.CreateShippingResponse, error) {

	log.Printf("CalculateShipping called for order_id=%d", request.OrderId)

	var items []domain.ShippingItem
	for _, item := range request.OrderItems {
		items = append(items, domain.ShippingItem{
			ProductCode: item.ProductCode,
			Quantity:    item.Quantity,
		})
	}

	newShipping := domain.NewShipping(request.OrderId, items)

	result, err := a.api.CalculateShipping(ctx, newShipping)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate shipping: %v", err)
	}

	return &shipping.CreateShippingResponse{
		DeliveryDays: result.DeliveryDays,
	}, nil
}
