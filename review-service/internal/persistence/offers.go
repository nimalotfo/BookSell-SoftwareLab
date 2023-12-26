package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"gitlab.com/narm-group/review-service/internal/database"
	"gitlab.com/narm-group/review-service/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

var (
	ErrOfferNotFound = errors.New("offer not found")
)

func GetOffer(ctx context.Context, offerId int64) (*models.BookOffer, error) {
	offer := &models.BookOffer{}
	filter := &models.BookOffer{
		ID: offerId,
	}

	db := database.GetDB(ctx)

	err := db.Preload("OfferImages").Where(filter).First(&offer).Error
	//	err := db.Where(filter).First(&offer).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOfferNotFound
		}
		return nil, err
	}

	fmt.Printf("getOffer-> %#v\n", offer)
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

func CreateOffer(db *gorm.DB, offer *models.BookOffer) (int64, error) {
	//imageUrlQuery := `INSERT INTO offer_images VALUES(default, $1, $2)`

	err := execTx(db, func(tx *gorm.DB) error {
		err := tx.Create(&offer).Error
		if err != nil {
			return err
		}

		fmt.Printf("offer -> %#v\n", offer)

		for _, offerImage := range offer.OfferImages {
			offerImage.OfferID = offer.ID
			err = tx.Save(offerImage).Error
			//err = tx.Exec(imageUrlQuery, offer.ID, offerImage).Error
			if err != nil {
				return err
			}
		}

		//offer.OfferImages = offer.OfferImages
		return nil
	})

	if err != nil {
		logrus.Errorf("error creating offer %d: %v\n", offer.ID, err)
		return 0, err
	}
	fmt.Printf("after creating offer -> %#v\n", offer)

	return offer.ID, nil
}

func GetPendingOffers(db *gorm.DB, count int) (offers []*models.BookOffer, err error) {
	//offers := make([]*models.BookOffer, 0)

	fmt.Println("in getPendingOffers persistence")
	err = db.Model(models.BookOffer{}).
		Preload("OfferImages").
		Where(&models.BookOffer{Status: models.Pending}).
		Order("created_at ASC").
		Limit(count).
		Scan(&offers).
		Error

	return
}

func GetUserOffers(ctx context.Context, userId int64, offerStatus models.OfferStatus, count int) (offers []*models.BookOffer, err error) {
	db := database.GetDB(ctx)

	if offerStatus.Value() == models.Unknown {
		err = status.Errorf(codes.InvalidArgument, "offer status is invalid")
		return
	}

	err = db.Model(models.BookOffer{}).
		Preload("OfferImages").
		Where(&models.BookOffer{Status: offerStatus, OwnerId: userId}).
		Order("created_at ASC").
		Limit(count).
		Scan(&offers).
		Error

	return
}
