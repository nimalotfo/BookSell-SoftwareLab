package events

type OfferApproved struct {
	OfferID     int64    `json:"offer_id"`
	OwnerId     int64    `json:"owner_id"`
	ReviewerId  int64    `json:"reviewer_id"`
	Name        string   `json:"name"`
	Price       int64    `json:"price"`
	PriceDeal   bool     `json:"price_deal"`
	Isbn        string   `json:"isbn"`
	Publisher   string   `json:"publisher"`
	Edition     int32    `json:"edition"`
	Description string   `json:"description"`
	ImageUrls   []string `json:"image_urls"`
}

func (event *OfferApproved) EventName() string {
	return "offer_approved"
}
