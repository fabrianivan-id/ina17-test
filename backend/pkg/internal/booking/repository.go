package booking

import (
	"concert-booking/pkg/models"
	"errors"

	"gorm.io/gorm"
)

type BookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) Create(booking *models.Booking) error {
	return r.db.Create(booking).Error
}

func (r *BookingRepository) FindByID(id string) (*models.Booking, error) {
	var booking models.Booking
	err := r.db.Preload("Concert").First(&booking, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &booking, err
}

func (r *BookingRepository) FindConcert(id uint) (*models.Concert, error) {
	var concert models.Concert
	err := r.db.First(&concert, id).Error
	return &concert, err
}

func (r *BookingRepository) GetBookedCount(concertID uint) (int, error) {
	var count int64
	err := r.db.Model(&models.Booking{}).
		Where("concert_id = ? AND status = 'confirmed'", concertID).
		Count(&count).Error
	return int(count), err
}
