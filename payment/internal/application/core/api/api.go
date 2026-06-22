package api

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/jessyekessia/microservices/payment/internal/application/core/domain"
	"github.com/jessyekessia/microservices/payment/internal/ports"
)

type Application struct {
	db ports.DBPort
}

func NewApplication(db ports.DBPort) *Application {
	return &Application{
		db: db,
	}
}

func (a Application) Charge( ctx context.Context, payment domain.Payment) (domain.Payment, error) {

    if payment.TotalPrice > 1000 {
        return domain.Payment{},
            status.Errorf(
                codes.InvalidArgument,
                "Payment over 1000 is not allowed.",
            )
    }

    err := a.db.Save(ctx, &payment)
    if err != nil {
        return domain.Payment{}, err
    }

    return payment, nil
}
