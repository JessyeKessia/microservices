package ports

import (
	"context"

	"github.com/jessyekessia/microservices/shipping/internal/application/core/domain"
)

type APIPort interface {
	CalculateShipping(ctx context.Context, shipping domain.Shipping) (domain.Shipping, error)
}
