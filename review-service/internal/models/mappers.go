package models

import "gitlab.com/narm-group/service-api/api/reviewpb"

func ToOfferList(offers []*BookOffer) *reviewpb.OfferList {
	arr := make([]*reviewpb.Offer, len(offers))
	for i, offer := range offers {
		arr[i] = &reviewpb.Offer{
			Id:          offer.ID,
			OwnerId:     offer.OwnerId,
			Name:        offer.Name,
			Price:       offer.Price,
			PriceDeal:   offer.PriceDeal,
			ImageUrls:   ToImageUrls(offer.OfferImages),
			Isbn:        offer.Isbn,
			Publisher:   offer.Publisher,
			Edition:     offer.Edition,
			Description: offer.Description,
		}
	}

	return &reviewpb.OfferList{
		Offers: arr,
	}
}

func ToImageUrls(offerImages []OfferImage) []string {
	images := []string{}
	for _, img := range offerImages {
		images = append(images, img.ImageUrl)
	}
	return images
}

func ToOfferImages(urls []string) []OfferImage {
	images := []OfferImage{}
	for _, url := range urls {
		images = append(images, OfferImage{ImageUrl: url})
	}
	return images
}
