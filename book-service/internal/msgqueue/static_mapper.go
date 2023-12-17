package msgqueue

import (
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"gitlab.com/narm-group/service-api/events"
)

type StaticMapper struct{}

func (mapper *StaticMapper) MapEvent(name string, payload interface{}) (Event, error) {
	var event Event

	switch name {
	case "offer_approved":
		event = &events.OfferApproved{}
	default:
		return nil, fmt.Errorf("event %s is unknown\n", name)
	}

	config := &mapstructure.DecoderConfig{
		TagName: "json",
	}
	config.Result = &event
	decoder, _ := mapstructure.NewDecoder(config)
	err := decoder.Decode(payload)
	if err != nil {
		return nil, err
	}

	switch t := event.(type) {
	case *events.OfferApproved:
		fmt.Println("received event description: ", t.Description)
		return t, nil
	default:
		msg := fmt.Sprintf("event %s paylaod is of invalid type\n", name)
		logrus.Error(msg)
		return nil, errors.New(msg)
	}
}
