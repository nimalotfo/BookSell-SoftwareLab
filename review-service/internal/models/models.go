package models

import "time"

type BookOffer struct {
	ID          int64        `json:"id"`
	OwnerId     int64        `json:"owner_id"`
	Name        string       `json:"name"`
	Price       int64        `json:"price"`
	PriceDeal   bool         `json:"price_deal"`
	Isbn        string       `json:"isbn"`
	Publisher   string       `json:"publisher"`
	Edition     int32        `json:"edition"`
	Description string       `json:"description"`
	Status      OfferStatus  `json:"status"`
	OfferImages []OfferImage `gorm:"foreignKey:OfferID" json:"offer_images"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (*BookOffer) TableName() string {
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

type OfferImage struct {
	ID       int64
	OfferID  int64
	ImageUrl string
}

type Review struct {
	ID          int64       `json:"id"`
	OfferId     int64       `json:"offer_id"`
	OfferStatus OfferStatus `json:"offer_status"`
	Description string      `json:"description"`
	ReviewerId  int64       `json:"reviewer_id"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Review) TableName() string {
	return "reviews"
}
