package payment

import (
	"concert-booking/pkg/models"
	"errors"
)

type PaymentService struct {
	repo *PaymentRepository
}

func NewPaymentService(repo *PaymentRepository) *PaymentService {
	return &PaymentService{repo: repo}
}

func (s *PaymentService) ProcessPayment(req *models.PaymentRequest) (*models.Payment, error) {
	payment := &models.Payment{
		BookingID: req.BookingID,
		Amount:    req.Amount,
		Status:    "success",
	}

	if err := s.repo.Create(payment); err != nil {
		return nil, errors.New("failed to process payment")
	}

	// Update booking status
	if err := s.repo.UpdateBookingStatus(req.BookingID, "confirmed"); err != nil {
		return nil, errors.New("failed to update booking status")
	}

	return payment, nil
}
