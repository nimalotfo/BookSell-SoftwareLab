package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"gitlab.com/narm-group/offer-service/internal/models"
	"gorm.io/gorm"
)

var (
	ErrOfferNotFound = errors.New("offer not found")
)

func GetOffer(ctx context.Context, db *gorm.DB, offerId int64) (*models.BookOffer, error) {
	offer := &models.BookOffer{}
	filter := &models.BookOffer{
		ID: offerId,
	}
	err := db.Where(filter).First(&offer).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOfferNotFound
		}
		return nil, err
	}

	return offer, nil
}

func execTx(db *gorm.DB, fn func(tx *gorm.DB) error) error {
	ok := false
	tx := db.Begin(nil)

	defer func() {
		if !ok {
			if err := tx.Rollback().Error; err != nil {
				logrus.Errorf("rollback error : %v\n", err)
			}
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	err := fn(tx)
	if err != nil {
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		logrus.Errorf("error commiting transaction: %v\n", err)
		return err
	}

	ok = true
	return nil
}

func CreateOffer(db *gorm.DB, offer *models.BookOffer, imageUrls []string) (int64, error) {
	imageUrlQuery := `INSERT INTO offer_images VALUES(default, $1, $2)`

	err := execTx(db, func(tx *gorm.DB) error {
		err := tx.Model(offer).Create(&offer).Error
		if err != nil {
			return err
		}

		for _, url := range imageUrls {
			err = tx.Exec(imageUrlQuery, offer.ID, url).Error
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		logrus.Errorf("error creating offer %d: %v\n", offer.ID, err)
		return 0, err
	}

	fmt.Printf("inserted offer id : %d\n", offer.ID)
	return offer.ID, nil
}
