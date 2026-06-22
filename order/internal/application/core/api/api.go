package api

import (
	"github.com/jessyekessia/microservices/order/internal/application/core/domain"
	"github.com/jessyekessia/microservices/order/internal/ports"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
)

type Application struct {
	db ports.DBPort
	payment ports.PaymentPort
}

func NewApplication (db ports.DBPort, payment ports.PaymentPort ) * Application {
	return & Application {
		 db : db ,
		 payment : payment ,
	}
}

func (a Application) PlaceOrder(
	order domain.Order,
) (domain.Order, error) {

	// Validação dos 50
	if err := order.Validate(); err != nil {
		return domain.Order{}, status.Errorf(
			codes.InvalidArgument,
			"Order cannot have more than 50 items",
		)
	}
	

	err := a.db.Save(&order)
	if err != nil {
		return domain.Order{}, err
	}

	paymentErr := a.payment.Charge(&order)

	if paymentErr != nil {

		order.Status = "Canceled"

		updateErr := a.db.Update(
			order.ID,
			order.Status,
		)

		if updateErr != nil {
			return domain.Order{}, updateErr
		}

		return domain.Order{}, paymentErr
	}

	order.Status = "Paid"

	err = a.db.Update(
		order.ID,
		order.Status,
	)

	if err != nil {
		return domain.Order{}, err
	}

	return order, nil
}