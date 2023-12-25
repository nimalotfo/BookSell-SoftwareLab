package events

type OfferSubmitted struct {
	OfferID     int64    `json:"offer_id"`
	OwnerId     int64    `json:"owner_id"`
	Name        string   `json:"name"`
	Price       int64    `json:"price"`
	PriceDeal   bool     `json:"price_deal"`
	Isbn        string   `json:"isbn"`
	Publisher   string   `json:"publisher"`
	Edition     int32    `json:"edition"`
	Description string   `json:"description"`
	ImageUrls   []string `json:"image_urls"`
	Status      int32    `json:"status"`
}

func (event *OfferSubmitted) EventName() string {
	return "offer_submitted"
}
