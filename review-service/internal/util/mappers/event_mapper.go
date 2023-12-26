package mappers

import (
	"gitlab.com/narm-group/review-service/internal/models"
	"gitlab.com/narm-group/service-api/events"
)

func NewReviewSubmittedEvent(review models.Review) *events.ReviewSubmitted {
	return &events.ReviewSubmitted{
		OfferId:     review.OfferId,
		OfferStatus: int32(review.OfferStatus),
		Description: review.Description,
		CreatedAt:   review.CreatedAt,
	}
}

func NewApprovedOfferEvent(offer *models.BookOffer, review *models.Review) *events.OfferApproved {
	images := []string{}
	for _, img := range offer.OfferImages {
		images = append(images, img.ImageUrl)
	}

	return &events.OfferApproved{
		OfferID:     offer.ID,
		OwnerId:     offer.OwnerId,
		ReviewerId:  review.ReviewerId,
		Name:        offer.Name,
		Price:       offer.Price,
		PriceDeal:   offer.PriceDeal,
		Isbn:        offer.Isbn,
		Publisher:   offer.Publisher,
		Edition:     offer.Edition,
		Description: offer.Description,
		ImageUrls:   images,
	}
}
