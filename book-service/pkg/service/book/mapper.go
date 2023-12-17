package book

import (
	"gitlab.com/narm-group/book-service/internal/models"
	"gitlab.com/narm-group/service-api/api/bookpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func toPbOffers(offers []*models.BookOffer) []*bookpb.Offer {
	arr := make([]*bookpb.Offer, len(offers))
	for i, offer := range offers {
		arr[i] = toPbOffer(offer)
	}

	return arr
}

func toPbOffer(offer *models.BookOffer) *bookpb.Offer {
	return &bookpb.Offer{
		Id:          offer.ID,
		Title:       offer.Name,
		BookTitle:   offer.Name,
		Isbn:        offer.Isbn,
		Publisher:   offer.Publisher,
		Edition:     offer.Edition,
		OwnerId:     offer.OwnerId,
		Price:       offer.Price,
		PriceDeal:   offer.PriceDeal,
		Description: offer.Description,
		ImageUrls:   ToImageUrls(offer.OfferImages),
		CreatedAt:   timestamppb.New(offer.CreatedAt),
		UpdatedAt:   timestamppb.New(offer.UpdatedAt),
	}
}

func ToImageUrls(offerImages []models.OfferImage) []string {
	images := []string{}
	for _, img := range offerImages {
		images = append(images, img.ImageUrl)
	}
	return images
}

func ToOfferImages(urls []string) []models.OfferImage {
	images := []models.OfferImage{}
	for _, url := range urls {
		images = append(images, models.OfferImage{ImageUrl: url})
	}
	return images
}
