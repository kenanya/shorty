package helpertest

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	// log "github.com/sirupsen/logrus"

	cm "github.com/kenanya/shorty/common"
)

// func InitRepository() *as.Client {
// 	hosts := make([]*as.Host, 0)

// 	for _, host := range cm.Config.Databases {
// 		params := strings.Split(host, ":")
// 		if port, e := strconv.Atoi(params[1]); e != nil {
// 			continue
// 		} else {
// 			hosts = append(hosts, as.NewHost(params[0], port))
// 		}
// 	}

// 	repo, e := as.NewClientWithPolicyAndHost(nil, hosts...)
// 	if e != nil {
// 		log.Errorf("Unable to initialize repository %v", e)
// 		os.Exit(1)
// 	}

// 	return repo
// }

// func InitDb() *sqlx.DB {
// 	db, err := sqlx.Connect("mysql", cm.Config.DbConnectUrl)
// 	if err != nil {
// 		log.Error(err)
// 		os.Exit(1)
// 	}
// 	return db
// }

func ConnectToDB() (*mongo.Database, error) {

	// reg := codecs.Register(bson.NewRegistryBuilder()).Build()
	log.Printf("connecting to MongoDB...")
	// logger.Log.Info("connecting to MongoDB...")

	cm.InitConfig("unit_test")
	fmt.Println(cm.ConfEnv)

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
		// logger.Log.Fatal("failed to create new MongoDB client", zap.String("reason", err.Error()))
		log.Fatalf("failed to create new MongoDB client", err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// connect client
	if err = client.Connect(ctx); err != nil {
		log.Fatalf("failed to connect to MongoDB", err.Error())
	}
	log.Println("connected successfully")

	db := client.Database(cm.ConfEnv.DatastoreDBSchema)
	return db, err
}
