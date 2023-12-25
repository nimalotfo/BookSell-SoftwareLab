package events

import "time"

type ReviewSubmitted struct {
	OfferId     int64     `json:"offer_id"`
	OfferStatus int32     `json:"offer_status"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func (event *ReviewSubmitted) EventName() string {
	return "review_submitted"
}
