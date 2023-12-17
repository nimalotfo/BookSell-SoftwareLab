package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.com/narm-group/book-service/internal/database"
	"gitlab.com/narm-group/book-service/internal/models"
	"gorm.io/gorm"
)

var (
	ErrOfferNotFound = errors.New("offer not found")
)

func GetBookInfo(ctx context.Context, offerId int64) (*models.BookOffer, error) {
	offer := &models.BookOffer{}
	filter := &models.BookOffer{
		ID: offerId,
	}

	db := database.GetDB(ctx)

	err := db.Preload("OfferImages").Where(filter).First(&offer).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOfferNotFound
		}
		return nil, err
	}

	return offer, nil
}

func GetBookOffers(ctx context.Context, fp FilterParams) ([]*models.BookOffer, error) {
	db := database.GetDB(ctx)
	query := db.Model(&models.BookOffer{}).Where("status = ?", models.ACTIVE)

	if fp.FromDate.UnixNano() > 0 {
		query = query.Where("updated_at >= ?", fp.FromDate)
	}
	if fp.ToDate.UnixNano() > 0 {
		query = query.Where("updated_at <= ?", fp.ToDate)
	}

	if fp.FromPrice > 0 {
		query = query.Where("price >= ?", fp.FromPrice)
	}
	if fp.ToPrice > 0 {
		query = query.Where("price <= ?", fp.ToPrice)
	}

	if fp.UserId > 0 {
		query = query.Where("owner_id = ?", fp.UserId)
	}

	if fp.PriceDealStatus == CHECKED || fp.PriceDealStatus == UNCHECKED {
		priceDeal := fp.PriceDealStatus == CHECKED
		query = query.Where("price_deal = ?", priceDeal)
	}

	if len(fp.Name) > 0 {
		query = query.Where("name LIKE ?", fmt.Sprintf("%%%s%%", fp.Name))
	}

	var offers []*models.BookOffer
	err := query.Find(&offers).Error
	if err != nil {
		return nil, err
	}

	for _, offer := range offers {
		var images []models.OfferImage
		err = db.Table("offer_images").Where("offer_id = ?", offer.ID).Scan(&images).Error
		if err != nil {
			return nil, err
		}
		offer.OfferImages = images
	}

	return offers, nil
}

func InsertBook(ctx context.Context, model *models.BookOffer) (insertedId int64, err error) {
	db := database.GetDB(ctx)

	model.Status = models.ACTIVE

	err = execTx(db, func(tx *gorm.DB) error {

		err = tx.Create(&model).Error
		if err != nil {
			return err
		}

		for _, img := range model.OfferImages {
			img.OfferID = model.ID
			fmt.Printf("image : %#v\n", img)
			err = tx.Save(img).Error
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return
	}

	return model.ID, nil
}

type PriceDealStatus int32

const (
	ANY PriceDealStatus = iota
	CHECKED
	UNCHECKED
)

type FilterParams struct {
	FromDate        time.Time
	ToDate          time.Time
	FromPrice       int64
	ToPrice         int64
	PriceDealStatus PriceDealStatus
	UserId          int64
	Name            string
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
