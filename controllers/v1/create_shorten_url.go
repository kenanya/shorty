package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kenanya/shorty/lib/flowV1"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	// log "github.com/sirupsen/logrus"
)

func (c *Controller) CreateShortenURL(ctx echo.Context) (err error) {

	var (
		req   flowV1.ParamShortenURLRequest
		res   ResponsePayload
		curDB *mongo.Database
	)

	curDB = c.CurDB

	json.NewDecoder(ctx.Request().Body).Decode(&req)
	// log.WithField("request", req).Info("Change Password request received")
	fmt.Println("## Controller -> CreateShortenURL ")

	rs, errorCode, err := flowV1.CreateShortenURL(curDB, req)

	if err != nil {
		res = c.getErrorResponse(ctx, errorCode, err.Error())
		return ctx.JSON(http.StatusBadRequest, res)
	}

	res = c.getSuccessResponse(ctx, errorCode, rs, 1, []string{})
	return ctx.JSON(http.StatusOK, res)
}
