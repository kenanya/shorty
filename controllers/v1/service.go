package v1

import (
	"time"

	"github.com/kenanya/shorty/lib/helper"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	// "errors"
	// log "github.com/sirupsen/logrus"
	// "math"
	// "net/http"
	// "registry/lib/constanta"
	// "registry/lib/helper"
	// "registry/lib/security"
	// "strconv"
	// "strings"
	// "time"
	// cm "registry/common"
	// "registry/models"
	// as "github.com/aerospike/aerospike-client-go"
	// "github.com/jmoiron/sqlx"
	// "github.com/labstack/echo/v4"
)

// type RegistryMessage struct {
// 	//OTP related params
// 	MDN string `json:"mdn,omitempty" validate:"omitempty,mdn"`
// 	OTP string `json:"otp,omitempty" validate:"omitempty,max=6,min=4"`

// 	//Signup related params
// 	Email    string `json:"email,omitempty" validate:"omitempty,email"`
// 	Password string `json:"password,omitempty" validate:"omitempty,max=50,min=6"`
// 	Name     string `json:"name,omitempty" validate:"omitempty,max=50,min=2"`
// 	VerifyBy string `json:"verifyBy,omitempty" validate:"omitempty,max=5,min=2"`

// 	//JWT related params
// 	GrantType    []models.JwtScope `json:"grant_type,omitempty" validate:"-"`
// 	AccessToken  string            `json:"token,omitempty" validate:"-"`
// 	RefreshToken string            `json:"refresh,omitempty" validate:"-"`

// 	//sign in related params
// 	Members []Member `json:"members,omitempty" validate:"-"`

// 	//default
// 	Error   int    `json:"error"`
// 	Message string `json:"message"`
// }

// type Member struct {
// 	MDN    string `json:"mdn,omitempty" validate:"omitempty,mdn"`
// 	Name   string `json:"name,omitempty" validate:"omitempty,mdn"`
// 	Status string `json:"status,omitempty" validate:"omitempty,mdn"`
// }

type ResponsePayload struct {
	Error     int         `json:"error"`
	Message   string      `json:"message"`
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Total     int         `json:"total"`
	NextPage  string      `json:"next_page"`
	FirstPage string      `json:"first_page"`
	LastPage  string      `json:"last_page"`
	Timestamp time.Time   `json:"timestamp"`
}

type Controller struct {
	CurDB *mongo.Database
}

/**
Get Success Response Payload
*/
func (this *Controller) getSuccessResponse(ctx echo.Context, errorCode int, data interface{}, total int, params []string) ResponsePayload {
	s := ResponsePayload{}
	s.Error = errorCode
	s.Message = "Success"
	s.Type = "object"
	s.Data = data
	s.Total = total
	s.NextPage = ""
	s.FirstPage = ""
	s.LastPage = ""

	s.Timestamp = helper.GetNowTime()

	return s
}

/**
Get Error Response Payload
*/
func (this *Controller) getErrorResponse(ctx echo.Context, errorCode int, errorMsg string) ResponsePayload {
	s := ResponsePayload{}
	s.Error = errorCode
	s.Message = errorMsg
	s.Type = "object"
	s.Total = 1
	s.NextPage = ""
	s.FirstPage = ""
	s.LastPage = ""

	s.Timestamp = helper.GetNowTime()

	return s
}
