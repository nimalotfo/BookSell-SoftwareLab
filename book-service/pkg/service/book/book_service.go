package book

import (
	"context"
	"net/http"

	"gitlab.com/narm-group/book-service/internal/persistence"
	"gitlab.com/narm-group/service-api/api/bookpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidPriceFilter = status.Errorf(http.StatusBadRequest, "toPrice can't be less than fromPrice")
	ErrInvalidDateFilter  = status.Errorf(http.StatusBadRequest, "toDate can't be less than fromDate")
)

type BookService struct {
	bookpb.UnimplementedBookServiceServer
}

func RegisterGrpcService(s *grpc.Server) {
	bookpb.RegisterBookServiceServer(s, &BookService{})
}

func (s *BookService) GetBookOffers(ctx context.Context, fp *bookpb.FilterParams) (*bookpb.OfferList, error) {
	if fp.FromPrice > fp.ToPrice {
		return nil, ErrInvalidPriceFilter
	}

	fromDate := fp.FromDate.AsTime()
	toDate := fp.ToDate.AsTime()

	if fromDate.After(toDate) {
		return nil, ErrInvalidDateFilter
	}

	filterParams := &persistence.FilterParams{
		UserId:          fp.UserId,
		FromDate:        fromDate,
		ToDate:          toDate,
		FromPrice:       fp.FromPrice,
		ToPrice:         fp.ToPrice,
		Name:            fp.Name,
		PriceDealStatus: persistence.PriceDealStatus(fp.PriceDeal),
	}

	offers, err := persistence.GetBookOffers(ctx, *filterParams)
	if err != nil {
		return nil, err
	}

	return &bookpb.OfferList{
		Offers: toPbOffers(offers),
	}, nil

}

func (s *BookService) GetBookInfo(ctx context.Context, req *bookpb.IdVersion) (*bookpb.Offer, error) {
	offer, err := persistence.GetBookInfo(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return toPbOffer(offer), nil

}

// func getUserIdFromCtx(ctx context.Context) (int64, error) {
// 	md, _ := metadata.FromIncomingContext(ctx)
// 	userIdStr := md["user_id"][0]

// 	logrus.Infof("fetch userId from ctx -> %s\n", userIdStr)
// 	userId, err := strconv.Atoi(userIdStr)
// 	if err != nil {
// 		return 0, errors.New(fmt.Sprintf("user id is invalid %s", userIdStr))
// 	}
// 	return int64(userId), nil
// }
