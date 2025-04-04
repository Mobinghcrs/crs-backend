package services

import "errors"

var (
	ErrFlightNotFound   = errors.New("flight not found")
	ErrNotEnoughSeats   = errors.New("not enough seats")
	ErrDuplicateBooking = errors.New("duplicate booking")
	ErrUserExists       = errors.New("user already exists")
)