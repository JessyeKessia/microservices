package shipping_adapter

import (
	"context"
	"log"
	"time"

	"github.com/jessyekessia/microservices-proto/golang/shipping"
	"github.com/jessyekessia/microservices/order/internal/application/core/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Adapter struct {
	shipping shipping.ShippingClient
}

func NewAdapter(shippingServiceUrl string) (*Adapter, error) {
	conn, err := grpc.Dial(
		shippingServiceUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	client := shipping.NewShippingClient(conn)
	return &Adapter{shipping: client}, nil
}

func (a *Adapter) Calculate(order *domain.Order) (int32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var items []*shipping.OrderItem
	for _, item := range order.OrderItems {
		items = append(items, &shipping.OrderItem{
			ProductCode: item.ProductCode,
			Quantity:    item.Quantity,
		})
	}

	resp, err := a.shipping.CalculateShipping(ctx, &shipping.CreateShippingRequest{
		OrderId:    order.ID,
		OrderItems: items,
	})

	if err != nil {
		if status.Code(err) == codes.DeadlineExceeded {
			log.Println("Deadline excedido ao chamar o microsserviço Shipping.")
		}
		return 0, err
	}

	return resp.DeliveryDays, nil
}
