package api

import (
	"context"

	"github.com/jessyekessia/microservices/shipping/internal/application/core/domain"
	"github.com/jessyekessia/microservices/shipping/internal/ports"
)

type Application struct {
	db ports.DBPort
}

func NewApplication(db ports.DBPort) *Application {
	return &Application{db: db}
}

func (a Application) CalculateShipping(ctx context.Context, shipping domain.Shipping) (domain.Shipping, error) {
	err := a.db.Save(ctx, &shipping)
	if err != nil {
		return domain.Shipping{}, err
	}
	return shipping, nil
}
