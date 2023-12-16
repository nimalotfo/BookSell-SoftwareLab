package review

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gitlab.com/narm-group/gateway/internal/security"
	"gitlab.com/narm-group/service-api/api/reviewpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ReviewHandler struct{}

var reviewClient reviewpb.ReviewServiceClient

func RegisterGrpcClient(url string) {
	cc, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		logrus.Fatal(err)
	}

	reviewClient = reviewpb.NewReviewServiceClient(cc)

	logrus.Infof("connection to review service grpc: %s\n", cc.GetState().String())
}

func (h *ReviewHandler) GetUserOffers(c *gin.Context) {
	var req reviewpb.UserOfferReq
	if err := c.BindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, "provide valid body")
		return
	}

	ctx := c.Request.Context()

	userId, err := security.GetCurrentUserId(ctx)
	if err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(
		"user_id", fmt.Sprintf("%d", userId),
	))

	res, err := reviewClient.GetUserOffers(ctx, &req)
	fmt.Printf("res -> %#v\n", res)
	fmt.Printf("err -> %v\n", err)
	if err != nil {
		c.Error(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"offers": res.Offers,
	})

}

func (h *ReviewHandler) GetPendingOffers(c *gin.Context) {
	var req reviewpb.PendingOffersReq
	if err := c.BindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, "provide valid body")
		return
	}

	ctx := c.Request.Context()

	userId, err := security.GetCurrentUserId(ctx)
	if err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized")
		return
	}

	//ctx = context.WithValue(ctx, "userId", userId)
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(
		"user_id", fmt.Sprintf("%d", userId),
	))

	res, err := reviewClient.GetPendingOffers(ctx, &req)
	if err != nil {
		errStatus, ok := status.FromError(err)
		if !ok {
			c.String(http.StatusInternalServerError, err.Error())
		} else {
			c.String(int(errStatus.Code()), errStatus.Message())
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"offers": res.Offers,
	})
}
func (h *ReviewHandler) SubmitReviewResult(c *gin.Context) {
	var req reviewpb.OfferReview
	if err := c.BindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, "provide valid body")
		return
	}

	ctx := c.Request.Context()

	userId, err := security.GetCurrentUserId(ctx)
	if err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized")
		return
	}

	//ctx = context.WithValue(ctx, "userId", userId)
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(
		"user_id", fmt.Sprintf("%d", userId),
	))

	if req.Status == reviewpb.OfferReview_UNKNOWN {
		c.String(http.StatusBadRequest, "status is invalid")
		return
	}

	_, err = reviewClient.SubmitReviewResult(ctx, &req)
	if err != nil {
		errStatus, ok := status.FromError(err)
		if !ok {
			c.String(http.StatusInternalServerError, err.Error())
		} else {
			c.String(int(errStatus.Code()), errStatus.Message())
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok": "true",
	})
}
