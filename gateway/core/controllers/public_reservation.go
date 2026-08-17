package controllers

import (
	"net/http"
	"strings"
	"time"

	"github.com/abaldeweg/warehouse-server/gateway/core/models"
	"github.com/abaldeweg/warehouse-server/gateway/core/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PublicReservationController handles reservation-related HTTP requests.
type PublicReservationController struct {
	db              *gorm.DB
	reservationRepo repository.ReservationRepository
}

// NewPublicReservationController creates a new PublicReservationController.
func NewPublicReservationController(db *gorm.DB) *PublicReservationController {
	return &PublicReservationController{
		db:              db,
		reservationRepo: repository.NewReservationRepository(db),
	}
}

// Create creates a new reservation.
func (rc *PublicReservationController) Create(c *gin.Context) {
	var reservationForm models.ReservationForm
	if err := c.ShouldBindJSON(&reservationForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reservation := models.Reservation{
		ID:         uuid.New().String(),
		CreatedAt:  time.Now(),
		Notes:      reservationForm.Notes,
		Salutation: reservationForm.Salutation,
		Firstname:  reservationForm.Firstname,
		Surname:    reservationForm.Surname,
		Mail:       reservationForm.Mail,
		Phone:      reservationForm.Phone,
		Open:       true,
	}

	var booksToUpdate []*models.Book
	for bookID := range strings.SplitSeq(reservationForm.Books, ",") {
		bookID = strings.TrimSpace(bookID)
		if bookID == "" {
			continue
		}
		var book models.Book
		if err := rc.db.Preload("Branch").First(&book, "id = ?", bookID).Error; err == nil {
			if book.Sold || book.Removed || book.Reserved {
				continue
			}
			booksToUpdate = append(booksToUpdate, &book)
			reservation.Books = append(reservation.Books, &book)
		}
	}

	if len(reservation.Books) > 0 && reservation.Books[0].BranchID != nil {
		reservation.BranchID = *reservation.Books[0].BranchID
	}

	if !reservation.Validate(rc.db) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reservation data"})
		return
	}

	if err := rc.reservationRepo.Create(&reservation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reservation"})
		return
	}

	rid := uuid.MustParse(reservation.ID)
	now := time.Now()
	for _, book := range booksToUpdate {
		book.Reserved = true
		book.ReservedAt = &now
		book.ReservationID = &rid
		rc.db.Select("reserved", "reserved_at", "reservation_id").Save(book)
	}

	_, err := rc.reservationRepo.FindOne(uuid.MustParse(reservation.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Reservation created, but failed to retrieve"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"msg": "SUCCESS"})
}
