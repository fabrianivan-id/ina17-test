package models

import "time"

type Payment struct {
	ID        string    `gorm:"primaryKey;type:char(36)"`
	BookingID string    `gorm:"not null;uniqueIndex"`
	Amount    float64   `gorm:"not null;type:decimal(10,2)"`
	Status    string    `gorm:"not null;default:'pending'"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type PaymentRequest struct {
	BookingID string  `json:"booking_id" binding:"required"`
	Amount    float64 `json:"amount" binding:"required"`
}
