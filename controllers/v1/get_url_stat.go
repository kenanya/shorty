package v1

import (
	"net/http"

	"github.com/kenanya/shorty/lib/flowV1"
	"github.com/kenanya/shorty/pkg/logger"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func (c *Controller) GetURLStatByShortCode(ctx echo.Context) (err error) {

	var (
		// req   flowV1.ParamShortenURLRequest
		shortcode string
		curDB     *mongo.Database
	)

	curDB = c.CurDB
	shortcode = ctx.Param("shortcode")
	logger.Log.Info("## Controller GetURLStatByShortCode shortcode : ", zap.String("reason", shortcode))

	rs, errorCode, err := flowV1.GetURLStatByShortCode(curDB, shortcode)

	if err != nil {
		res := c.getErrorResponse(ctx, err.Error())
		return ctx.JSON(errorCode, res)
	}

	return ctx.JSON(http.StatusOK, rs)

}
