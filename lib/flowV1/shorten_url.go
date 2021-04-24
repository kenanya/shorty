package flowV1

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/kenanya/shorty/lib/helper"
	"github.com/kenanya/shorty/pkg/logger"
	"github.com/lucasjones/reggen"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

const (
	collName    = "url_bank"
	desc409     = "The the desired shortcode is already in use. Shortcodes are case-sensitive."
	desc422     = "The shortcode fails to meet the following regexp: ^[0-9a-zA-Z_]{6}$."
	desc404     = "The shortcode cannot be found in the system"
	scodePatern = "^[0-9a-zA-Z_]{6}$"
)

type ParamShortenURLRequest struct {
	Url       string `json:"url"`
	ShortCode string `json:"shortcode"`
}

type ResponseCreateShortenURL struct {
	ShortCode string `json:"shortcode"`
}

type ResponseGetURL struct {
	Location string `json:"location"`
}

type ResponseGetURLStat struct {
	StartDate     string `json:"startDate"`
	LastSeenDate  string `json:"lastSeenDate,omitempty"`
	RedirectCount int    `json:"redirectCount"`
}

type URLModel struct {
	ID            primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	URL           string             `json:"url" valid:"Required"`
	ShortCode     string             `json:"shortcode"`
	StartDate     string             `json:"startDate,omitempty"`
	LastSeenDate  string             `json:"lastSeenDate,omitempty"`
	RedirectCount int                `json:"redirectCount,omitempty"`
}

func CreateShortenURL(CurDB *mongo.Database, req ParamShortenURLRequest) (rs ResponseCreateShortenURL, errorCode int, err error) {

	if req.Url == "" {
		logger.Log.Error("url is not present")
		return rs, http.StatusBadRequest, errors.New("url is not present")
	} else {
		_, errUrl := url.ParseRequestURI(req.Url)
		if errUrl != nil {
			logger.Log.Error("## CreateShortenURL Error invalid URL", zap.String("reason", errUrl.Error()))
			return rs, http.StatusBadRequest, errors.New("Invalid URL")
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	coll := CurDB.Collection(collName)
	ts := helper.GetNowTime().Format("2006-01-02T15:04:05-0700")

	tempShort := ""
	if req.ShortCode != "" {
		tempShort = req.ShortCode
		if isExistShortcode(ctx, CurDB, tempShort) {
			logger.Log.Error(desc409)
			return rs, http.StatusConflict, errors.New(desc409)
		}
		if !isMatchPattern(tempShort) {
			logger.Log.Error(desc422)
			return rs, http.StatusUnprocessableEntity, errors.New(desc422)
		}
	} else {
		// isExistShortcode
		tempShort = generateShortCode()
		for isExistShortcode(ctx, CurDB, tempShort) || tempShort == "" {
			tempShort = generateShortCode()
			logger.Log.Info("## new shortcode : " + tempShort)
		}
	}

	in := URLModel{
		ID:            primitive.NewObjectID(),
		URL:           req.Url,
		ShortCode:     tempShort,
		StartDate:     ts,
		LastSeenDate:  "",
		RedirectCount: 0,
	}

	fmt.Printf("## CreateShortenURL in : <%+v>", in)
	res, err := coll.InsertOne(ctx, &in)
	if err != nil {
		logger.Log.Error("insert data into collection <"+collName+">", zap.String("reason", err.Error()))
		return rs, 0, errors.New("error insert data into collection")
	}
	fmt.Printf("## res.InsertedID : <%+v> \n", res.InsertedID)

	rs.ShortCode = in.ShortCode
	return rs, 0, nil
}

func GetURLByShortCode(CurDB *mongo.Database, shortcode string) (rs ResponseGetURL, errorCode int, err error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	coll := CurDB.Collection(collName)
	ts := helper.GetNowTime().Format("2006-01-02T15:04:05-0700")

	var out, up *URLModel
	var filter = bson.D{}
	if shortcode == "" {
		logger.Log.Error("shortcode is not present")
		return rs, http.StatusBadRequest, errors.New("shortcode is not present")
	} else {
		filter = bson.D{
			{Key: "shortcode", Value: shortcode},
		}
	}
	fmt.Printf("%+v\n\n\n", filter)

	// // Read data from collection
	err = coll.FindOne(ctx, filter).Decode(&out)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			logger.Log.Info("No data found")
			return rs, http.StatusNotFound, errors.New(desc404)
		} else {
			logger.Log.Error("Fail to read data from collection <"+collName+">", zap.String("reason", err.Error()))
			return rs, http.StatusInternalServerError, errors.New("Fail to read data from collection <" + collName + ">")
		}
	}
	fmt.Printf("%#v\n\n\n", out)

	if out != nil {
		up = &URLModel{
			RedirectCount: out.RedirectCount + 1,
			LastSeenDate:  ts,
			ID:            out.ID,
		}
		err = updateURLData(ctx, CurDB, up)
		if err != nil {
			logger.Log.Error("Fail to update RedirectCount and LastSeenDate", zap.String("reason", err.Error()))
			return rs, http.StatusInternalServerError, errors.New("Fail to update RedirectCount and LastSeenDate")
		}

	}

	rs.Location = out.URL
	return rs, 0, nil
}

func GetURLStatByShortCode(CurDB *mongo.Database, shortcode string) (rs ResponseGetURLStat, errorCode int, err error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	coll := CurDB.Collection(collName)

	var out *URLModel
	var filter = bson.D{}
	if shortcode == "" {
		logger.Log.Error("shortcode is not present")
		return rs, http.StatusBadRequest, errors.New("shortcode is not present")
	} else {
		filter = bson.D{
			{Key: "shortcode", Value: shortcode},
		}
	}
	fmt.Printf("%+v\n\n\n", filter)

	// // Read data from collection
	err = coll.FindOne(ctx, filter).Decode(&out)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			logger.Log.Info("No data found")
			return rs, http.StatusNotFound, errors.New(desc404)
		} else {
			logger.Log.Error("Fail to read data from collection <"+collName+">", zap.String("reason", err.Error()))
			return rs, http.StatusInternalServerError, errors.New("Fail to read data from collection <" + collName + ">")
		}
	}
	fmt.Printf("%#v\n\n\n", out)

	rs.StartDate = out.StartDate
	rs.RedirectCount = out.RedirectCount
	if rs.RedirectCount > 0 {
		rs.LastSeenDate = out.LastSeenDate
	}

	return rs, 0, nil
}

func isExistShortcode(ctx context.Context, CurDB *mongo.Database, scode string) bool {

	coll := CurDB.Collection(collName)
	isExist := false

	var out *URLModel
	var filter = bson.D{
		{Key: "shortcode", Value: scode},
	}
	fmt.Printf("%+v\n\n\n", filter)

	// // Read data from collection
	err := coll.FindOne(ctx, filter).Decode(&out)
	if err != nil {
		logger.Log.Error("Fail to read data from collection <"+collName+">", zap.String("reason", err.Error()))
	} else {
		if out.ShortCode != "" {
			isExist = true
		}
	}

	fmt.Printf("%#v\n\n\n", out)
	return isExist
}

func updateURLData(ctx context.Context, CurDB *mongo.Database, urlData *URLModel) error {

	coll := CurDB.Collection(collName)

	var filter = bson.D{}
	filter = bson.D{
		{Key: "_id", Value: urlData.ID},
	}

	var setElements bson.D
	setElements = append(setElements, bson.E{Key: "redirectcount", Value: urlData.RedirectCount})
	setElements = append(setElements, bson.E{Key: "lastseendate", Value: urlData.LastSeenDate})

	updateRec := bson.D{
		{Key: "$set", Value: setElements},
	}

	// Update data in the collection
	res, err := coll.UpdateOne(ctx, filter, updateRec)
	if err != nil {
		logger.Log.Error("Update data into collection <"+collName+"> failed : '%s'", zap.String("reason", err.Error()))
		return errors.New("Update data into collection <" + collName + "> failed: '%s'" + err.Error())
	}
	// fmt.Printf("Matched %v documents and updated %v documents.\n", res.MatchedCount, res.ModifiedCount)

	if res.ModifiedCount > 0 {
		logger.Log.Info("Data has been updated successfully")
		return nil
	} else {
		logger.Log.Error("Update data into collection <" + collName + "> failed")
		return errors.New("Update data into collection <" + collName + "> failed")
	}

}

func isMatchPattern(scode string) bool {
	r, _ := regexp.Compile(scodePatern)
	result := r.Match([]byte(scode))

	return result
}

func generateShortCode() string {
	str, err := reggen.Generate(scodePatern, 6)
	if err != nil {
		logger.Log.Error("Error generating shortcode", zap.String("reason", err.Error()))
		str = ""
	}

	return str
}
