package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	cm "github.com/kenanya/shorty/common"
	"github.com/kenanya/shorty/pkg/logger"
	"github.com/kenanya/shorty/router"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// ExtraConfig is additional configuration for Server
type ExtraConfig struct {
	// Log parameters section
	// LogLevel is global log level: Debug(-1), Info(0), Warn(1), Error(2), DPanic(3), Panic(4), Fatal(5)
	LogLevel int
	// LogTimeFormat is print time format for logger e.g. 2006-01-02T15:04:05Z07:00
	LogTimeFormat string
}

func ConnectToDB() (*mongo.Database, error) {

	// reg := codecs.Register(bson.NewRegistryBuilder()).Build()
	// log.Printf("connecting to MongoDB...")
	logger.Log.Info("connecting to MongoDB...")

	uri := fmt.Sprintf(`mongodb://%s:%s@%s/%s`,
		cm.ConfEnv.DatastoreDBUser,
		cm.ConfEnv.DatastoreDBPassword,
		cm.ConfEnv.DatastoreDBHost,
		cm.ConfEnv.DatastoreDBSchema,
	)
	// uri := fmt.Sprintf(`mongodb://%s`,
	// 	cm.ConfEnv.DatastoreDBHost,
	// )

	client, err := mongo.NewClient(
		options.Client().ApplyURI(uri),
		// &options.ClientOptions{
		// 	Registry: reg,
		// }
	)

	if err != nil {
		logger.Log.Fatal("failed to create new MongoDB client", zap.String("reason", err.Error()))
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// connect client
	if err = client.Connect(ctx); err != nil {
		logger.Log.Fatal("failed to connect to MongoDB", zap.String("reason", err.Error()))
	}
	logger.Log.Info("connected successfully")

	db := client.Database(cm.ConfEnv.DatastoreDBSchema)
	return db, err
}

func main() {
	var cfg ExtraConfig

	flag.IntVar(&cfg.LogLevel, "log-level", 0, "Global log level")
	flag.StringVar(&cfg.LogTimeFormat, "log-time-format", "",
		"Print time format for logger e.g. 2006-01-02T15:04:05Z07:00")
	flag.Parse()

	// initialize logger
	if err := logger.Init(cfg.LogLevel, cfg.LogTimeFormat); err != nil {
		fmt.Printf("failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	cm.InitConfig()

	e := echo.New()

	//initilize storage connections
	db, err := ConnectToDB()
	if err != nil {
		// log.Fatalf("failed initialize MongoDB connection: %#v", err)
		logger.Log.Fatal("failed initialize MongoDB connection", zap.String("reason", err.Error()))
	}

	e.Server.Addr = cm.ConfEnv.RestPort

	//assign routing
	router.AssignRouting(e, db)
}
