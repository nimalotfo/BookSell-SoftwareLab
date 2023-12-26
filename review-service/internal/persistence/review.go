package persistence

import (
	"errors"

	"gitlab.com/narm-group/review-service/internal/models"
	"gorm.io/gorm"
)

var (
	ErrInvalidReview = errors.New("review shall not have id")
)

func SubmitReview(db *gorm.DB, review *models.Review) (insertedId int64, err error) {
	err = execTx(db, func(tx *gorm.DB) error {
		err = db.Model(&models.Review{}).Save(&review).Error
		if err != nil {
			return err
		}

		if review.OfferStatus != models.Pending {
			err = db.Model(&models.BookOffer{}).
				Where("id = ?", review.OfferId).
				Updates(&models.BookOffer{Status: review.OfferStatus}).
				Error
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return
	}

	return review.ID, nil
}

func GetReview(db *gorm.DB, reviewId int64) (review *models.Review, err error) {
	err = db.Take(review, &models.Review{ID: reviewId}).Error
	if err != nil {
		return nil, err
	}

	return
}
