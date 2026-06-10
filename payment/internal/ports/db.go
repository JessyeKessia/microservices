package ports

import (
	"context"
	"github.com/jessyekessia/microservices/payment/internal/application/core/domain"
)

type DBPort interface {
	Get(ctx context.Context, id string) (domain.Payment, error)
	Save(ctx context.Context, payment *domain.Payment) error
}
