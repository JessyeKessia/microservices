package ports

import "github.com/jessyekessia/microservices/order/internal/application/core/domain"

type DBPort interface {
	Get(id string) (domain.Order, error)
	Save(*domain.Order) error
	Update(id int64, status string) error
}