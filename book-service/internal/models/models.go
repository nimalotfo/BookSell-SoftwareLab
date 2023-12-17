package models

import "time"

type BookOffer struct {
	ID          int64        `json:"id"`
	OwnerId     int64        `json:"owner_id"`
	Isbn        string       `json:"isbn"`
	Name        string       `json:"name"`
	Price       int64        `json:"price"`
	PriceDeal   bool         `json:"price_deal"`
	Publisher   string       `json:"publisher"`
	Edition     int32        `json:"edition"`
	Description string       `json:"description"`
	OfferImages []OfferImage `gorm:"foreignKey:OfferID" json:"offer_images"`
	Status      OfferStatus  `json:"status"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (*BookOffer) TableName() string {
	return "offers"
}

type OfferStatus int32

const (
	UNKNOWN OfferStatus = iota
	ACTIVE
	DEACITVE
)

type OfferImage struct {
	ID       int64
	OfferID  int64
	ImageUrl string
}
