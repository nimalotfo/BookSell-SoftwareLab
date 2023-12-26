package handler

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
	"gitlab.com/narm-group/review-service/internal/database"
	"gitlab.com/narm-group/review-service/internal/models"
	"gitlab.com/narm-group/review-service/internal/msgqueue"
	"gitlab.com/narm-group/review-service/internal/persistence"
	events_api "gitlab.com/narm-group/service-api/events"
)

var (
	ErrInvalidEventType = errors.New("event is not of offer submitted format")
)

func HandleEvents(ctx context.Context, events <-chan msgqueue.Event, errors <-chan error) error {
	for {
		select {
		case lisErr := <-errors:
			logrus.Errorf("listener error : %v\n", lisErr)
		case evt := <-events:
			offerSubEvt, ok := evt.(*events_api.OfferSubmitted)
			if !ok {
				logrus.Errorf("event is not of offer submitted format")
				return ErrInvalidEventType
			}
			db := database.GetDB(ctx)
			_, err := persistence.CreateOffer(db, toOffer(offerSubEvt))
			if err != nil {
				logrus.Errorf("error creating offer-> %v\n", err)
				return err
			}
		}
	}
}

func toOffer(evt *events_api.OfferSubmitted) *models.BookOffer {
	return &models.BookOffer{
		Name:        evt.Name,
		OwnerId:     evt.OwnerId,
		Price:       evt.Price,
		PriceDeal:   evt.PriceDeal,
		Isbn:        evt.Isbn,
		Publisher:   evt.Publisher,
		Edition:     evt.Edition,
		Description: evt.Description,
		OfferImages: models.ToOfferImages(evt.ImageUrls),
		Status:      models.OfferStatus(evt.Status),
	}
}
