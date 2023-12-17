package models

type BookOffer struct {
	ID          int64 `json:"id"`
	OwnerId     int64 `json:"owner_id"`
	Name        string	`json:"name"`
	Price       int64	`json:"price"`
	PriceDeal   bool	`json:"price_deal"`
	Isbn        string	`json:"isbn"`
	Publisher   string	`json:"publisher"`
	Edition     int32	`json:"edition"`
	Description string	`json:"description"`
	Status      OfferStatus	`json:"status"`
}

func (BookOffer) TableName() string {
	return "offers"
}

type OfferStatus int32

const (
	Unknown OfferStatus = iota
	Pending
	Accepted
	Rejected
)

func (s OfferStatus) Value() OfferStatus {
	switch s {
	case Pending, Accepted, Rejected:
		return s
	default:
		return Unknown
	}
}

type OfferDto struct {
	ID          int64
	OwnerId     int64
	Name        string
	Price       int64
	PriceDeal   bool
	ImageUrls   []string
	Isbn        string
	Publisher   string
	Edition     int32
	Description string
}

func NewOfferDto(offer BookOffer, imageUrls []string) *OfferDto {
	return &OfferDto{
		ID:          offer.ID,
		OwnerId:     offer.OwnerId,
		Name:        offer.Name,
		Price:       offer.Price,
		PriceDeal:   offer.PriceDeal,
		ImageUrls:   imageUrls,
		Isbn:        offer.Isbn,
		Publisher:   offer.Publisher,
		Edition:     offer.Edition,
		Description: offer.Description,
	}
}
