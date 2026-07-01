package api

import (
	"log"

	"github.com/jessyekessia/microservices/order/internal/application/core/domain"
	"github.com/jessyekessia/microservices/order/internal/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Application struct {
	db       ports.DBPort
	payment  ports.PaymentPort
	shipping ports.ShippingPort
}

func NewApplication(db ports.DBPort, payment ports.PaymentPort, shipping ports.ShippingPort) *Application {
	return &Application{
		db:       db,
		payment:  payment,
		shipping: shipping,
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

	// Validação de estoque
	for _, item := range order.OrderItems {
		exists, err := a.db.ProductExists(item.ProductCode)
		if err != nil {
			return domain.Order{}, status.Errorf(
				codes.Internal,
				"failed to check stock for product %s: %v", item.ProductCode, err,
			)
		}
		if !exists {
			return domain.Order{}, status.Errorf(
				codes.NotFound,
				"product %s not found in stock", item.ProductCode,
			)
		}
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

	deliveryDays, shippingErr := a.shipping.Calculate(&order)
	if shippingErr != nil {
		log.Printf("Warning: failed to calculate shipping for order %d: %v", order.ID, shippingErr)
	} else {
		log.Printf("Shipping calculated for order %d: %d delivery days", order.ID, deliveryDays)
	}

	return order, nil
}