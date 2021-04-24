package flowV1

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kenanya/shorty/lib/helper"
	"github.com/kenanya/shorty/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/status"
)

// "bytes"
// "errors"
// "fmt"
// "reflect"
// "database/sql"
// as "github.com/aerospike/aerospike-client-go"
// "github.com/astaxie/beego/validation"

// cm "registry/common"
// "registry/lib/helper"
// // "registry/lib/security"
// "registry/lib/constanta"
// "registry/models"
// "strings"
// "time"

const collName = "url_bank"

type ParamShortenURLRequest struct {
	Url       string `json:"url"`
	ShortCode string `json:"shortcode"`
}

type ResponseShortenURL struct {
	ShortCode string `json:"shortcode"`
}

type URLModel struct {
	URL           string `json:"url"`
	ShortCode     string `json:"shortcode"`
	StartDate     string `json:"startDate,omitempty"`
	LastSeenDate  string `json:"lastSeenDate,omitempty"`
	RedirectCount int    `json:"redirectCount,omitempty"`
}

func CreateShortenURL(CurDB *mongo.Database, req ParamShortenURLRequest) (rs ResponseShortenURL, errorCode int, err error) {

	// textTEST := "TEST dulu ya TG"
	// fmt.Println(textTEST)
	// // helper.SendSms("628811218706", textTEST)
	// return rs, 0, nil

	// var (
	// 	ErrMessage, tempMail, tempMdn, acctId string
	// 	valid                                 validation.Validation
	// 	b                                     bool
	// )

	// if req.Data == "" {
	// 	return rs, cm.ErrParameterInvalid, errors.New("Missing/invalid parameter(s)")
	// }

	// return rs, cm.ErrAccountNotFound, errors.New("Account Not Found")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	coll := CurDB.Collection(collName)
	ts := helper.GetNowTime().Format("2006-01-02T15:04:05-0700")

	tempShort := "tes12345"
	if req.ShortCode != "" {
		tempShort = req.ShortCode
	}

	in := URLModel{
		URL:           req.Url,
		ShortCode:     tempShort,
		StartDate:     ts,
		LastSeenDate:  "",
		RedirectCount: 0,
	}
	// in.CreatedDate = ts

	res, err := coll.InsertOne(ctx, &in)
	if err != nil {
		logger.Log.Error("insert data into collection <"+collName+">", zap.String("reason", err.Error()))
		return rs, 0, errors.New("error insert data into collection")
	}
	// logger.Log.Info("res.InsertedID : " + res.InsertedID.(primitive.ObjectID))
	fmt.Printf("## res.InsertedID : <%+v> \n", res.InsertedID)
	// id := pmongo.NewObjectId(res.InsertedID.(primitive.ObjectID))

	return rs, 0, nil
}

// func (s *initServiceServer) CreateAnalyticTelemetry(ctx context.Context, req *v1.AnalyticTelemetryRequest) (*v1.AnalyticTelemetryResponse, error) {
// 	// check if the API version requested by client is supported by server
// 	if err := s.checkAPI(req.Api); err != nil {
// 		return nil, err
// 	}

// 	curDB := s.db
// 	coll := curDB.Collection(collName)

// 	loc, err := time.LoadLocation("Asia/Jakarta")
// 	if err != nil {
// 		logger.Log.Error("failed to load location", zap.String("reason", err.Error()))
// 		return nil, status.Errorf(codes.Internal, "failed to load location '%s'", err.Error())
// 	}

// 	t := time.Now().In(loc)
// 	ts, err := ptypes.TimestampProto(t)
// 	if err != nil {
// 		logger.Log.Error("failed to convert golang Time to protobuf Timestamp", zap.String("reason", err.Error()))
// 		return nil, status.Error(codes.Internal, "failed to convert golang Time to protobuf Timestamp")
// 	}

// 	in := req.AnalyticData
// 	in.CreatedDate = ts

// 	res, err := coll.InsertOne(ctx, &in)
// 	if err != nil {
// 		logger.Log.Error("insert data into collection <"+collName+">", zap.String("reason", err.Error()))
// 		return nil, status.Errorf(codes.Internal, "insert data into collection <"+collName+"> '%s'", err.Error())
// 	}
// 	id := pmongo.NewObjectId(res.InsertedID.(primitive.ObjectID))

// 	return &v1.AnalyticTelemetryResponse{
// 		Api: apiVersion,
// 		Id:  id,
// 	}, nil
// }
