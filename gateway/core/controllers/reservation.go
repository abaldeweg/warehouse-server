package controllers

import (
	"net/http"
	"strings"
	"time"

	"github.com/abaldeweg/warehouse-server/gateway/auth"
	"github.com/abaldeweg/warehouse-server/gateway/core/models"
	"github.com/abaldeweg/warehouse-server/gateway/core/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ReservationController handles reservation-related HTTP requests.
type ReservationController struct {
	db              *gorm.DB
	reservationRepo repository.ReservationRepository
}

// NewReservationController creates a new ReservationController.
func NewReservationController(db *gorm.DB) *ReservationController {
	return &ReservationController{
		db:              db,
		reservationRepo: repository.NewReservationRepository(db),
	}
}

// FindAll retrieves all reservations for the current user's branch.
func (rc *ReservationController) FindAll(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}

	reservations, err := rc.reservationRepo.FindAll(uint(user.(auth.User).Branch.Id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, reservations)
}

// ReservationStatus retrieves the number of open reservations for the current user's branch.
func (rc *ReservationController) ReservationStatus(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}

	count, err := rc.reservationRepo.ReservationStatus(uint(user.(auth.User).Branch.Id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"open": count})
}

// FindOne retrieves a reservation by its UUID.
func (rc *ReservationController) FindOne(c *gin.Context) {
	uuid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	user, ok := c.Get("user")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}

	reservation, err := rc.reservationRepo.FindOne(uuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Reservation not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if uint(user.(auth.User).Branch.Id) != reservation.BranchID {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"msg": "Forbidden"})
		return
	}

	c.JSON(http.StatusOK, reservation)
}

// Create creates a new reservation.
func (rc *ReservationController) Create(c *gin.Context) {
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

	reservation.Books = make([]*models.Book, 0)
	for bookID := range strings.SplitSeq(reservationForm.Books, ",") {
		bookID = strings.TrimSpace(bookID)
		if bookID == "" {
			continue
		}
		var book models.Book
		if err := rc.db.First(&book, "id = ?", bookID).Error; err == nil {
			if book.Sold || book.Removed || book.Reserved {
				continue
			}
			book.Reserved = true
			book.ReservedAt = time.Now()
			book.ReservationID = uuid.MustParse(reservation.ID)
			if err := rc.db.Save(&book).Error; err != nil {
				continue
			}
			reservation.Books = append(reservation.Books, &book)
		}
	}

	user, ok := c.Get("user")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}
	reservation.BranchID = uint(user.(auth.User).Branch.Id)

	if !reservation.Validate(rc.db) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reservation data"})
		return
	}

	if err := rc.reservationRepo.Create(&reservation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reservation"})
		return
	}

	createdReservation, err := rc.reservationRepo.FindOne(uuid.MustParse(reservation.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Reservation created, but failed to retrieve"})
		return
	}

	c.JSON(http.StatusCreated, createdReservation)
}

// Update updates an existing reservation by its UUID.
func (rc *ReservationController) Update(c *gin.Context) {
	uuid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	existingReservation, err := rc.reservationRepo.FindOne(uuid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reservation not found"})
		return
	}

	user, ok := c.Get("user")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}

	if uint(user.(auth.User).Branch.Id) != existingReservation.BranchID {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"msg": "Forbidden"})
		return
	}

	var reservationForm models.ReservationForm
	if err := c.ShouldBindJSON(&reservationForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reservation := models.Reservation{
		ID:         existingReservation.ID,
		CreatedAt:  existingReservation.CreatedAt,
		Notes:      reservationForm.Notes,
		Salutation: reservationForm.Salutation,
		Firstname:  reservationForm.Firstname,
		Surname:    reservationForm.Surname,
		Mail:       reservationForm.Mail,
		Phone:      reservationForm.Phone,
		Open:       reservationForm.Open,
	}

	// reservation.Books = make([]*models.Book, 0)
	// for bookID := range strings.SplitSeq(reservationForm.Books, ",") {
	// 	bookID = strings.TrimSpace(bookID)
	// 	if bookID == "" {
	// 		continue
	// 	}
	// 	var book models.Book
	// 	if err := rc.db.First(&book, "id = ?", bookID).Error; err == nil {
	// 		if book.Sold || book.Removed || book.Reserved {
	// 			continue
	// 		}
	// 		reservation.Books = append(reservation.Books, &book)
	// 	}
	// }

	reservation.BranchID = uint(user.(auth.User).Branch.Id)

	if !reservation.Validate(rc.db) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reservation data"})
		return
	}

	// reservation.ID = uuid.String()

	if err := rc.reservationRepo.Update(&reservation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update"})
		return
	}

	updatedReservation, err := rc.reservationRepo.FindOne(uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Reservation updated, but failed to retrieve"})
		return
	}

	c.JSON(http.StatusOK, updatedReservation)
}

// Delete deletes a reservation by its UUID.
func (rc *ReservationController) Delete(c *gin.Context) {
	uuid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, ok := c.Get("user")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}

	existingReservation, err := rc.reservationRepo.FindOne(uuid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reservation not found"})
		return
	}

	if uint(user.(auth.User).Branch.Id) != existingReservation.BranchID {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"msg": "Forbidden"})
		return
	}

	for i := range existingReservation.Books {
		existingReservation.Books[i].Reservation = nil
		if err := rc.db.Save(&existingReservation.Books[i]).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book reservation"})
			return
		}
	}

	if err := rc.reservationRepo.Delete(uuid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete"})
		return
	}

	c.Status(http.StatusNoContent)
}
