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

type tableEntity struct {
	col *mongo.Collection
}

type ITable interface {
	GetTables() ([]entities.Table, error)
	GetTableById(id primitive.ObjectID) (entities.Table, error)
	CreateTable(table entities.Table) (entities.Table, error)
	UpdateTableById(id primitive.ObjectID, table entities.Table) error
	DeleteTableById(id primitive.ObjectID) error
	UpdateTableStatus(id primitive.ObjectID, status string) error
}

func NewTableEntity(resource *db.Resource) ITable {
	col := resource.SnookDb.Collection("tables")
	return &tableEntity{col: col}
}

func (entity *tableEntity) GetTables() ([]entities.Table, error) {
	logrus.Info("GetTables")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	opts := options.Find().SetSort(bson.D{{Key: "name", Value: 1}})
	cursor, err := entity.col.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	var tables []entities.Table
	if err = cursor.All(ctx, &tables); err != nil {
		return nil, err
	}
	return tables, nil
}

func (entity *tableEntity) GetTableById(id primitive.ObjectID) (entities.Table, error) {
	logrus.Info("GetTableById")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var table entities.Table
	err := entity.col.FindOne(ctx, bson.M{"_id": id}).Decode(&table)
	return table, err
}

func (entity *tableEntity) CreateTable(table entities.Table) (entities.Table, error) {
	logrus.Info("CreateTable")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	table.Id = primitive.NewObjectID()
	table.CreatedDate = time.Now()
	table.UpdatedDate = time.Now()
	_, err := entity.col.InsertOne(ctx, table)
	return table, err
}

func (entity *tableEntity) UpdateTableById(id primitive.ObjectID, table entities.Table) error {
	logrus.Info("UpdateTableById")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	table.UpdatedDate = time.Now()
	_, err := entity.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{
		"name":        table.Name,
		"type":        table.Type,
		"ratePerHour": table.RatePerHour,
		"description": table.Description,
		"updatedBy":   table.UpdatedBy,
		"updatedDate": table.UpdatedDate,
	}})
	return err
}

func (entity *tableEntity) DeleteTableById(id primitive.ObjectID) error {
	logrus.Info("DeleteTableById")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := entity.col.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (entity *tableEntity) UpdateTableStatus(id primitive.ObjectID, status string) error {
	logrus.Info("UpdateTableStatus")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := entity.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"status": status}})
	return err
}
