package book

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gitlab.com/narm-group/service-api/api/bookpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BookHandler struct{}

var bookClient bookpb.BookServiceClient

func RegisterGrpcClient(url string) {
	cc, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		logrus.Fatal(err)
	}

	bookClient = bookpb.NewBookServiceClient(cc)

	logrus.Infof("connection to book service grpc: %s\n", cc.GetState().String())
}

func (h *BookHandler) GetBookInfo(c *gin.Context) {
	var req bookpb.IdVersion
	if err := c.BindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, "provide valid body")
		return
	}

	ctx := c.Request.Context()
	res, err := bookClient.GetBookInfo(ctx, &req)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *BookHandler) GetBookOffers(c *gin.Context) {
	var fp bookpb.FilterParams
	if err := c.BindJSON(&fp); err != nil {
		c.String(http.StatusBadRequest, "provide valid body")
		return
	}

	fmt.Println(fp.FromDate.AsTime())
	fmt.Println(fp.ToDate.AsTime())

	ctx := c.Request.Context()
	res, err := bookClient.GetBookOffers(ctx, &fp)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func HandleError(c *gin.Context, err error) {
	errStatus, ok := status.FromError(err)
	if !ok {
		c.String(http.StatusInternalServerError, err.Error())
	}

	var status int
	switch errStatus.Code() {
	case codes.PermissionDenied:
		status = http.StatusForbidden
	case codes.Unauthenticated:
		status = http.StatusUnauthorized
	case codes.FailedPrecondition, codes.AlreadyExists, codes.InvalidArgument:
		status = http.StatusBadRequest
	case codes.NotFound:
		status = http.StatusNotFound
	default:
		status = http.StatusInternalServerError
	}

	c.JSON(status, gin.H{
		"message": errStatus.Message(),
	})

}
