package repository

import (
	"github.com/abaldeweg/warehouse-server/gateway/core/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ReservationRepository interface defines methods for interacting with reservation data.
type ReservationRepository interface {
	FindAll(uint) ([]models.Reservation, error)
	ReservationStatus(branchID uint) (int64, error)
	FindOne(id uuid.UUID) (*models.Reservation, error)
	Create(reservation *models.Reservation) error
	Update(reservation *models.Reservation) error
	Delete(id uuid.UUID) error
}

// reservationRepository implements the ReservationRepository interface.
type reservationRepository struct {
	db *gorm.DB
}

// NewReservationRepository creates a new reservation repository.
func NewReservationRepository(db *gorm.DB) ReservationRepository {
	return &reservationRepository{db: db}
}

// FindAll retrieves all reservations for the given branch ID, ordered by creation date in descending order.
func (rr *reservationRepository) FindAll(branchID uint) ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := rr.db.Preload("Branch").Preload("Books").Where("branch_id = ?", branchID).Order("created_at DESC").Find(&reservations).Error
	return reservations, err
}

// ReservationStatus retrieves the number of open reservations for a given branch ID.
func (rr *reservationRepository) ReservationStatus(branchID uint) (int64, error) {
	var count int64
	err := rr.db.Model(&models.Reservation{}).Where("branch_id = ? AND open = ?", branchID, true).Count(&count).Error
	return count, err
}

// FindOne retrieves a reservation by its UUID.
func (rr *reservationRepository) FindOne(id uuid.UUID) (*models.Reservation, error) {
	var reservation models.Reservation
	err := rr.db.Preload("Branch").Preload("Books").First(&reservation, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &reservation, nil
}

// Create creates a new reservation.
func (rr *reservationRepository) Create(reservation *models.Reservation) error {
	return rr.db.Create(reservation).Error
}

// Update updates an existing reservation.
func (rr *reservationRepository) Update(reservation *models.Reservation) error {
	return rr.db.Save(reservation).Error
}

// Delete deletes a reservation by its UUID.
func (rr *reservationRepository) Delete(id uuid.UUID) error {
	return rr.db.Delete(&models.Reservation{}, "id = ?", id).Error
}
