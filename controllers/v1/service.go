package v1

import (
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

type ErrorResponse struct {
	Description string `json:"description"`
}

type Controller struct {
	CurDB *mongo.Database
}

/**
Get Error Response Payload
*/
func (this *Controller) getErrorResponse(ctx echo.Context, errorMsg string) ErrorResponse {
	s := ErrorResponse{}
	s.Description = errorMsg

	return s
}
