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

type menuItemEntity struct {
	col *mongo.Collection
}

type IMenuItem interface {
	GetMenuItems(category string) ([]entities.MenuItem, error)
	GetMenuItemById(id primitive.ObjectID) (entities.MenuItem, error)
	CreateMenuItem(item entities.MenuItem) (entities.MenuItem, error)
	UpdateMenuItemById(id primitive.ObjectID, item entities.MenuItem) error
	DeleteMenuItemById(id primitive.ObjectID) error
	UpdateMenuItemQuantity(id primitive.ObjectID, quantity int) error
	GetLowStockMenuItems(threshold int) ([]entities.LowStockMenuItem, error)
}

func NewMenuItemEntity(resource *db.Resource) IMenuItem {
	col := resource.SnookDb.Collection("menu_items")
	return &menuItemEntity{col: col}
}

func (entity *menuItemEntity) GetMenuItems(category string) ([]entities.MenuItem, error) {
	logrus.Info("GetMenuItems")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{}
	if category != "" {
		filter["category"] = category
	}
	opts := options.Find().SetSort(bson.D{{Key: "name", Value: 1}})
	cursor, err := entity.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var items []entities.MenuItem
	if err = cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (entity *menuItemEntity) GetMenuItemById(id primitive.ObjectID) (entities.MenuItem, error) {
	logrus.Info("GetMenuItemById")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var item entities.MenuItem
	err := entity.col.FindOne(ctx, bson.M{"_id": id}).Decode(&item)
	return item, err
}

func (entity *menuItemEntity) CreateMenuItem(item entities.MenuItem) (entities.MenuItem, error) {
	logrus.Info("CreateMenuItem")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	item.Id = primitive.NewObjectID()
	item.CreatedDate = time.Now()
	item.UpdatedDate = time.Now()
	_, err := entity.col.InsertOne(ctx, item)
	return item, err
}

func (entity *menuItemEntity) UpdateMenuItemById(id primitive.ObjectID, item entities.MenuItem) error {
	logrus.Info("UpdateMenuItemById")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	item.UpdatedDate = time.Now()
	_, err := entity.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{
		"name":        item.Name,
		"category":    item.Category,
		"price":       item.Price,
		"costPrice":   item.CostPrice,
		"quantity":    item.Quantity,
		"unit":        item.Unit,
		"status":      item.Status,
		"imageUrl":    item.ImageUrl,
		"updatedBy":   item.UpdatedBy,
		"updatedDate": item.UpdatedDate,
	}})
	return err
}

func (entity *menuItemEntity) DeleteMenuItemById(id primitive.ObjectID) error {
	logrus.Info("DeleteMenuItemById")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := entity.col.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (entity *menuItemEntity) UpdateMenuItemQuantity(id primitive.ObjectID, quantity int) error {
	logrus.Info("UpdateMenuItemQuantity")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := entity.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$inc": bson.M{"quantity": quantity}})
	return err
}

func (entity *menuItemEntity) GetLowStockMenuItems(threshold int) ([]entities.LowStockMenuItem, error) {
	logrus.Info("GetLowStockMenuItems")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{"quantity": bson.M{"$lte": threshold}, "status": "ACTIVE"}
	opts := options.Find().SetSort(bson.D{{Key: "quantity", Value: 1}})
	cursor, err := entity.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var items []entities.LowStockMenuItem
	if err = cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}
