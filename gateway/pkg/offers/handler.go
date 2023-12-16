package offers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gitlab.com/narm-group/gateway/internal/security"
	"gitlab.com/narm-group/service-api/api/offerpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type OfferHandler struct{}

var offerClient offerpb.OfferServiceClient

func RegisterGrpcClient(url string) {
	cc, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		logrus.Fatal(err)
	}

	offerClient = offerpb.NewOfferServiceClient(cc)

	logrus.Infof("connection to offer service grpc: %s\n", cc.GetState().String())
}

func (h *OfferHandler) SubmitOffer(c *gin.Context) {
	var req offerpb.SubmitOfferReq
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

	res, err := offerClient.SubmitOffer(ctx, &req)
	if err != nil {
		errStatus, ok := status.FromError(err)
		if !ok {
			c.String(http.StatusInternalServerError, err.Error())
		} else {
			c.String(http.StatusInternalServerError, errStatus.Message())
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": res.Id,
	})

}
