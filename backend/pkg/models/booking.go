package models

import "time"

type Concert struct {
	ID       uint   `gorm:"primaryKey"`
	Title    string `gorm:"not null"`
	Date     time.Time
	Venue    string `gorm:"not null"`
	Capacity int    `gorm:"not null"`
}

type Booking struct {
	ID        string    `gorm:"primaryKey;type:char(36)"`
	UserID    string    `gorm:"not null;index"`
	ConcertID uint      `gorm:"not null;index"`
	Quantity  int       `gorm:"not null;check:quantity > 0"`
	Status    string    `gorm:"not null;default:'pending'"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type BookingRequest struct {
	ConcertID uint `json:"concert_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}
