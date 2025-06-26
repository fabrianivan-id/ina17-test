package payment

import (
	"concert-booking/pkg/models"

	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Create(payment *models.Payment) error {
	return r.db.Create(payment).Error
}

func (r *PaymentRepository) UpdateBookingStatus(bookingID, status string) error {
	return r.db.Model(&models.Booking{}).
		Where("id = ?", bookingID).
		Update("status", status).Error
}
