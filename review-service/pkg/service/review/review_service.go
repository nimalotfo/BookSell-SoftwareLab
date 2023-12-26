package review

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
	"gitlab.com/narm-group/review-service/internal/database"
	"gitlab.com/narm-group/review-service/internal/models"
	"gitlab.com/narm-group/review-service/internal/msgqueue/kafka"
	"gitlab.com/narm-group/review-service/internal/persistence"
	"gitlab.com/narm-group/review-service/internal/util/mappers"
	"gitlab.com/narm-group/service-api/api/reviewpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	ErrInvalidStatus = errors.New("provided review status is invalid")
)

type ReviewService struct {
	reviewpb.UnimplementedReviewServiceServer
}

func RegisterGrpcService(s *grpc.Server) {
	reviewpb.RegisterReviewServiceServer(s, &ReviewService{})
}

func (s *ReviewService) GetUserOffers(ctx context.Context, req *reviewpb.UserOfferReq) (*reviewpb.OfferList, error) {
	userId, err := getUserIdFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	offers, err := persistence.GetUserOffers(ctx, userId, models.OfferStatus(req.GetStatus()), int(req.GetCount()))
	if err != nil {
		return nil, err
	}

	return models.ToOfferList(offers), nil
}

func (s *ReviewService) GetPendingOffers(ctx context.Context, req *reviewpb.PendingOffersReq) (*reviewpb.OfferList, error) {
	db := database.GetDB(ctx)
	logrus.Info("in pendingoffers review service")
	fmt.Printf("count is : %d\n", req.Count)
	offers, err := persistence.GetPendingOffers(db, int(req.Count))
	if err != nil {
		return nil, err
	}
	fmt.Println("len offers : ", len(offers))

	for _, offer := range offers {
		err = db.Table("offer_images").
			Where("offer_id = ?", offer.ID).
			Select("image_url").
			Find(&offer.OfferImages).Error
		if err != nil {
			return nil, err
		}
	}

	return models.ToOfferList(offers), nil
	//kafkaListener := kafka.GetKafkaListener()
}

func (s *ReviewService) SubmitReviewResult(ctx context.Context, req *reviewpb.OfferReview) (res *emptypb.Empty, err error) {
	userId, err := getUserIdFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	res = &emptypb.Empty{}

	var status models.OfferStatus
	if req.Status == reviewpb.OfferReview_APPROVED {
		status = models.Accepted
	} else if req.Status == reviewpb.OfferReview_REJECTED {
		status = models.Rejected
	} else {
		logrus.Errorf("status %d is invalid\n", req.Status)
		return nil, ErrInvalidStatus
	}

	review := &models.Review{
		OfferId:     req.OfferId,
		ReviewerId:  userId,
		OfferStatus: status,
		Description: req.Description,
	}

	db := database.GetDB(ctx)

	insertedId, err := persistence.SubmitReview(db, review)
	if err != nil {
		return nil, err
	}

	logrus.Infof("submitted review with id: %d on offer id : %d\n", insertedId, review.OfferId)

	kafkaEmitter := kafka.GetKafkaEventEmitter()

	reviewEvent := mappers.NewReviewSubmittedEvent(*review)
	err = kafkaEmitter.Emit(reviewEvent, "review_topic")
	if err != nil {
		logrus.Warnf("error emitting review result event to kafka broker: %v\n", err)
		return
	}

	offer, err := persistence.GetOffer(ctx, review.OfferId)
	if err != nil {
		logrus.Errorf("error while fetching offer id : %d : %v\n", review.OfferId, err)
		return
	}

	if review.OfferStatus == models.Accepted {
		appEvent := mappers.NewApprovedOfferEvent(offer, review)
		err = kafkaEmitter.Emit(appEvent, "offer_approved")
		fmt.Printf("sent offer_approved -> %#v\n", appEvent)
		if err != nil {
			logrus.Warnf("error emitting offer_approved result event to kafka broker: %v\n", err)
			return
		}
	}

	return
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
