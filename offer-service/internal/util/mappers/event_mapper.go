package mappers

import (
	"gitlab.com/narm-group/offer-service/internal/models"
	"gitlab.com/narm-group/service-api/events"
)

func NewOfferSubmittedEvent(offer *models.BookOffer, imageUrls []string) *events.OfferSubmitted {
	return &events.OfferSubmitted{
		OfferID:     offer.ID,
		OwnerId:     offer.OwnerId,
		Name:        offer.Name,
		Price:       offer.Price,
		PriceDeal:   offer.PriceDeal,
		Isbn:        offer.Isbn,
		Publisher:   offer.Publisher,
		Edition:     offer.Edition,
		Description: offer.Description,
		ImageUrls:   imageUrls,
		Status:      int32(offer.Status),
	}
}
