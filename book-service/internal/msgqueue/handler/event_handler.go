package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"gitlab.com/narm-group/book-service/internal/models"
	"gitlab.com/narm-group/book-service/internal/msgqueue"
	"gitlab.com/narm-group/book-service/internal/persistence"
	"gitlab.com/narm-group/book-service/pkg/service/book"
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
			offerAppEvt, ok := evt.(*events_api.OfferApproved)
			if !ok {
				logrus.Errorf("event is not of offer submitted format")
				return ErrInvalidEventType
			}

			fmt.Printf("fetched event _---> %#v\n", offerAppEvt)
			_, err := persistence.InsertBook(ctx, toBookOffer(offerAppEvt))
			if err != nil {
				logrus.Errorf("error creating offer-> %v\n", err)
				return err
			}
		}
	}
}

func toBookOffer(evt *events_api.OfferApproved) *models.BookOffer {
	return &models.BookOffer{
		ID:          evt.OfferID,
		OwnerId:     evt.OwnerId,
		Isbn:        evt.Isbn,
		Name:        evt.Name,
		Price:       evt.Price,
		PriceDeal:   evt.PriceDeal,
		Publisher:   evt.Publisher,
		Edition:     evt.Edition,
		Description: evt.Description,
		OfferImages: book.ToOfferImages(evt.ImageUrls),
	}
}
