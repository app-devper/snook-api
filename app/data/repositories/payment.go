package repositories

import (
	"context"
	"snook/app/data/entities"
	"snook/db"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type paymentEntity struct {
	col *mongo.Collection
}

type IPayment interface {
	GetPaymentsBySessionId(sessionId primitive.ObjectID) ([]entities.Payment, error)
	CreatePayment(payment entities.Payment) (entities.Payment, error)
	DeletePayment(id primitive.ObjectID) error
	GetPaymentsByDateRange(startDate, endDate time.Time) ([]entities.Payment, error)
}

func NewPaymentEntity(resource *db.Resource) IPayment {
	col := resource.SnookDb.Collection("payments")
	return &paymentEntity{col: col}
}

func (entity *paymentEntity) GetPaymentsBySessionId(sessionId primitive.ObjectID) ([]entities.Payment, error) {
	logrus.Info("GetPaymentsBySessionId")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	opts := options.Find().SetSort(bson.D{{Key: "createdDate", Value: -1}})
	cursor, err := entity.col.Find(ctx, bson.M{"sessionId": sessionId}, opts)
	if err != nil {
		return nil, err
	}
	var payments []entities.Payment
	if err = cursor.All(ctx, &payments); err != nil {
		return nil, err
	}
	return payments, nil
}

func (entity *paymentEntity) CreatePayment(payment entities.Payment) (entities.Payment, error) {
	logrus.Info("CreatePayment")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	payment.Id = primitive.NewObjectID()
	payment.CreatedDate = time.Now()
	_, err := entity.col.InsertOne(ctx, payment)
	return payment, err
}

func (entity *paymentEntity) DeletePayment(id primitive.ObjectID) error {
	logrus.Info("DeletePayment")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := entity.col.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (entity *paymentEntity) GetPaymentsByDateRange(startDate, endDate time.Time) ([]entities.Payment, error) {
	logrus.Info("GetPaymentsByDateRange")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{"createdDate": bson.M{"$gte": startDate, "$lte": endDate}}
	opts := options.Find().SetSort(bson.D{{Key: "createdDate", Value: -1}})
	cursor, err := entity.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var payments []entities.Payment
	if err = cursor.All(ctx, &payments); err != nil {
		return nil, err
	}
	return payments, nil
}
