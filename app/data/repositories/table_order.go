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

type tableOrderEntity struct {
	col *mongo.Collection
}

type ITableOrder interface {
	GetOrdersBySessionId(sessionId primitive.ObjectID) ([]entities.TableOrder, error)
	GetTableOrderById(id primitive.ObjectID) (entities.TableOrder, error)
	CreateTableOrder(order entities.TableOrder) (entities.TableOrder, error)
	UpdateTableOrder(id primitive.ObjectID, order entities.TableOrder) error
	DeleteTableOrder(id primitive.ObjectID) error
	GetOrdersByDateRange(startDate, endDate time.Time) ([]entities.TableOrder, error)
}

func NewTableOrderEntity(resource *db.Resource) ITableOrder {
	col := resource.SnookDb.Collection("table_orders")
	return &tableOrderEntity{col: col}
}

func (entity *tableOrderEntity) GetOrdersBySessionId(sessionId primitive.ObjectID) ([]entities.TableOrder, error) {
	logrus.Info("GetOrdersBySessionId")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	opts := options.Find().SetSort(bson.D{{Key: "createdDate", Value: -1}})
	cursor, err := entity.col.Find(ctx, bson.M{"sessionId": sessionId}, opts)
	if err != nil {
		return nil, err
	}
	var orders []entities.TableOrder
	if err = cursor.All(ctx, &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

func (entity *tableOrderEntity) GetTableOrderById(id primitive.ObjectID) (entities.TableOrder, error) {
	logrus.Info("GetTableOrderById")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var order entities.TableOrder
	err := entity.col.FindOne(ctx, bson.M{"_id": id}).Decode(&order)
	return order, err
}

func (entity *tableOrderEntity) CreateTableOrder(order entities.TableOrder) (entities.TableOrder, error) {
	logrus.Info("CreateTableOrder")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	order.Id = primitive.NewObjectID()
	order.CreatedDate = time.Now()
	_, err := entity.col.InsertOne(ctx, order)
	return order, err
}

func (entity *tableOrderEntity) UpdateTableOrder(id primitive.ObjectID, order entities.TableOrder) error {
	logrus.Info("UpdateTableOrder")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := entity.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{
		"quantity": order.Quantity,
		"discount": order.Discount,
		"total":    order.Total,
	}})
	return err
}

func (entity *tableOrderEntity) DeleteTableOrder(id primitive.ObjectID) error {
	logrus.Info("DeleteTableOrder")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := entity.col.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (entity *tableOrderEntity) GetOrdersByDateRange(startDate, endDate time.Time) ([]entities.TableOrder, error) {
	logrus.Info("GetOrdersByDateRange")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{"createdDate": bson.M{"$gte": startDate, "$lte": endDate}}
	opts := options.Find().SetSort(bson.D{{Key: "createdDate", Value: -1}})
	cursor, err := entity.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var orders []entities.TableOrder
	if err = cursor.All(ctx, &orders); err != nil {
		return nil, err
	}
	return orders, nil
}
