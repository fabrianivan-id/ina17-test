package booking

import (
	"concert-booking/pkg/models"
	"errors"
)

type BookingService struct {
	repo *BookingRepository
}

func NewBookingService(repo *BookingRepository) *BookingService {
	return &BookingService{repo: repo}
}

func (s *BookingService) CreateBooking(userID string, req *models.BookingRequest) (*models.Booking, error) {
	concert, err := s.repo.FindConcert(req.ConcertID)
	if err != nil {
		return nil, errors.New("concert not found")
	}

	bookedCount, err := s.repo.GetBookedCount(concert.ID)
	if err != nil {
		return nil, errors.New("failed to check availability")
	}

	if bookedCount+req.Quantity > concert.Capacity {
		return nil, errors.New("not enough tickets available")
	}

	booking := &models.Booking{
		UserID:    userID,
		ConcertID: concert.ID,
		Quantity:  req.Quantity,
		Status:    "pending",
	}

	if err := s.repo.Create(booking); err != nil {
		return nil, errors.New("failed to create booking")
	}

	return booking, nil
}
