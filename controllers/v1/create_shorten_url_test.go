// +build unit_create_shorten_url

package v1

import (

	// "fmt"
	"encoding/json"
	"flag"
	"net/http"
	"net/http/httptest"

	"strings"
	"testing"

	cm "github.com/kenanya/shorty/common"
	"github.com/kenanya/shorty/lib/flowV1"
	"github.com/kenanya/shorty/lib/helpertest"

	// "github.com/kenanya/shorty/lib/helpertest"
	"github.com/kenanya/shorty/pkg/logger"
	"go.uber.org/zap"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type ExtraConfig struct {
	// Log parameters section
	// LogLevel is global log level: Debug(-1), Info(0), Warn(1), Error(2), DPanic(3), Panic(4), Fatal(5)
	LogLevel int
	// LogTimeFormat is print time format for logger e.g. 2006-01-02T15:04:05Z07:00
	LogTimeFormat string
}

func TestCreateShortenURL(t *testing.T) {

	var cfg ExtraConfig
	flag.IntVar(&cfg.LogLevel, "log-level", 0, "Global log level")
	flag.StringVar(&cfg.LogTimeFormat, "log-time-format", "2006-01-02T15:04:05Z07:00",
		"Print time format for logger e.g. 2006-01-02T15:04:05Z07:00")
	flag.Parse()

	// initialize logger
	if err := logger.Init(cfg.LogLevel, cfg.LogTimeFormat); err != nil {
		t.Errorf("failed to initialize logger: %v", err)
	}

	t.Run("positive case create shorten url", func(t *testing.T) {

		var (
			e          = echo.New()
			resPayload = flowV1.ResponseCreateShortenURL{}
		)

		cm.InitConfig("unit_test")

		db, err := helpertest.ConnectToDB()
		if err != nil {
			logger.Log.Fatal("failed initialize MongoDB connection", zap.String("reason", err.Error()))
		}

		body := map[string]interface{}{
			// "url": "https://app.diagrams.net/",
			"url": "https://dillinger.io/",
		}

		bodyJson, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/shorty/shorten", strings.NewReader(string(bodyJson)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		handler := &Controller{CurDB: db}

		// Assertions
		if assert.NoError(t, handler.CreateShortenURL(ctx)) {
			//test http code response
			assert.Equal(t, http.StatusCreated, rec.Code)

			assert.NotNil(t, resPayload)
		}
	})

	t.Run("negative case create shorten url", func(t *testing.T) {

		var (
			e          = echo.New()
			resPayload = ErrorResponse{}
		)

		cm.InitConfig("unit_test")

		db, err := helpertest.ConnectToDB()
		if err != nil {
			logger.Log.Fatal("failed initialize MongoDB connection", zap.String("reason", err.Error()))
		}

		body := map[string]interface{}{
			"url":       "https://app.diagrams.net/",
			"shortcode": "x1x1x2x",
		}

		bodyJson, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/shorty/shorten", strings.NewReader(string(bodyJson)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		handler := &Controller{CurDB: db}

		// Assertions
		if assert.NoError(t, handler.CreateShortenURL(ctx)) {
			//test http code response
			assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

			//test message response
			json.NewDecoder(rec.Body).Decode(&resPayload)
			assert.NotEmpty(t, "The shortcode fails to meet the following regexp: ^[0-9a-zA-Z_]{6}$.", resPayload.Description)
		}
	})
}
