package ports

import "github.com/jessyekessia/microservices/order/internal/application/core/domain"

type ShippingPort interface {
	Calculate(order *domain.Order) (int32, error)
}
