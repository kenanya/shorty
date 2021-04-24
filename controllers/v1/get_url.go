package v1

import (
	"net/http"
	"net/url"

	"github.com/kenanya/shorty/lib/flowV1"
	"github.com/kenanya/shorty/pkg/logger"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func (c *Controller) GetURLByShortCode(ctx echo.Context) (err error) {

	var (
		// req   flowV1.ParamShortenURLRequest
		shortcode string
		curDB     *mongo.Database
	)

	curDB = c.CurDB
	shortcode = ctx.Param("shortcode")

	logger.Log.Info("## Controller GetURLStatByShortCode shortcode : ", zap.String("reason", shortcode))
	rs, errorCode, err := flowV1.GetURLByShortCode(curDB, shortcode)

	if err != nil {
		res := c.getErrorResponse(ctx, err.Error())
		return ctx.JSON(errorCode, res)
	}

	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	_, errUrl := url.ParseRequestURI(rs.Location)
	if errUrl != nil {
		logger.Log.Error("## GetURLByShortCode Error invalid URL", zap.String("reason", errUrl.Error()))
	} else {
		ctx.Response().Header().Set(echo.HeaderLocation, rs.Location)
	}
	ctx.Response().WriteHeader(http.StatusFound)

	return ctx.JSON(http.StatusFound, rs)

}
