package offer

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
	"gitlab.com/narm-group/offer-service/internal/database"
	"gitlab.com/narm-group/offer-service/internal/models"
	"gitlab.com/narm-group/offer-service/internal/msgqueue/kafka"
	"gitlab.com/narm-group/offer-service/internal/persistence"
	"gitlab.com/narm-group/offer-service/internal/util/mappers"
	offerpb "gitlab.com/narm-group/service-api/api/offerpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	topic = "topic1"
)

type OfferService struct {
	offerpb.UnimplementedOfferServiceServer
}

func RegisterGrpcService(s *grpc.Server) {
	offerpb.RegisterOfferServiceServer(s, &OfferService{})
}

func getUserIdFromCtx(ctx context.Context) (int64, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	userIdStr := md["user_id"][0]

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("user id is invalid %s", userIdStr))
	}
	return int64(userId), nil
}

func (s *OfferService) SubmitOffer(ctx context.Context, req *offerpb.SubmitOfferReq) (*offerpb.IdVersion, error) {
	userId, err := getUserIdFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	offer := &models.BookOffer{
		OwnerId:     int64(userId),
		Name:        req.Name,
		Price:       req.Price,
		PriceDeal:   req.PriceDeal,
		Isbn:        req.Isbn,
		Publisher:   req.Publisher,
		Edition:     req.Edition,
		Description: req.Description,
		Status:      models.Pending,
	}

	db := database.GetDB(ctx)

	insertedId, err := persistence.CreateOffer(db, offer, req.ImageUrls)
	if err != nil {
		logrus.Warn(err)
		return nil, status.Errorf(codes.Internal, "couldn't submit offer")
	}

	kafkaEmitter := kafka.GetKafkaEventEmitter()

	event := mappers.NewOfferSubmittedEvent(offer, req.ImageUrls)
	err = kafkaEmitter.Emit(event, topic)
	if err != nil {
		logrus.Errorf("error emitting event %s of offerId: %d on kafka -> %v\n", event.EventName(), offer.ID, err)
	}

	return &offerpb.IdVersion{
		Id: insertedId,
	}, nil
}
