package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Resource struct {
	SnookDb     *mongo.Database
	RdDb        *redis.Client
	mongoClient *mongo.Client
}

// Close use this method to close database connection
func (r *Resource) Close() {
	logrus.Warning("Closing all db connections")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if r.mongoClient != nil {
		if err := r.mongoClient.Disconnect(ctx); err != nil {
			logrus.Error("failed to disconnect mongo: ", err)
		}
	}
	if r.RdDb != nil {
		if err := r.RdDb.Close(); err != nil {
			logrus.Error("failed to close redis: ", err)
		}
	}
}

func InitResource() (*Resource, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Print(err)
	}

	host := os.Getenv("MONGO_HOST")
	snookDbName := os.Getenv("MONGO_SNOOK_DB_NAME")
	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(host))
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = mongoClient.Connect(ctx)
	if err != nil {
		return nil, err
	}

	redisHost := os.Getenv("REDIS_HOST")
	redisOp, err := redis.ParseURL(redisHost)
	if err != nil {
		return nil, err
	}
	rdb := redis.NewClient(redisOp)

	return &Resource{
		SnookDb:     mongoClient.Database(snookDbName),
		RdDb:        rdb,
		mongoClient: mongoClient,
	}, nil
}
