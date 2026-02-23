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

type menuCategoryEntity struct {
	col *mongo.Collection
}

type IMenuCategory interface {
	GetMenuCategories() ([]entities.MenuCategory, error)
	GetMenuCategoryById(id primitive.ObjectID) (entities.MenuCategory, error)
	CreateMenuCategory(cat entities.MenuCategory) (entities.MenuCategory, error)
	UpdateMenuCategoryById(id primitive.ObjectID, cat entities.MenuCategory) error
	DeleteMenuCategoryById(id primitive.ObjectID) error
}

func NewMenuCategoryEntity(resource *db.Resource) IMenuCategory {
	col := resource.SnookDb.Collection("menu_categories")
	return &menuCategoryEntity{col: col}
}

func (entity *menuCategoryEntity) GetMenuCategories() ([]entities.MenuCategory, error) {
	logrus.Info("GetMenuCategories")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	opts := options.Find().SetSort(bson.D{{Key: "sortOrder", Value: 1}})
	cursor, err := entity.col.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	var cats []entities.MenuCategory
	if err = cursor.All(ctx, &cats); err != nil {
		return nil, err
	}
	return cats, nil
}

func (entity *menuCategoryEntity) GetMenuCategoryById(id primitive.ObjectID) (entities.MenuCategory, error) {
	logrus.Info("GetMenuCategoryById")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var cat entities.MenuCategory
	err := entity.col.FindOne(ctx, bson.M{"_id": id}).Decode(&cat)
	return cat, err
}

func (entity *menuCategoryEntity) CreateMenuCategory(cat entities.MenuCategory) (entities.MenuCategory, error) {
	logrus.Info("CreateMenuCategory")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cat.Id = primitive.NewObjectID()
	cat.CreatedDate = time.Now()
	cat.UpdatedDate = time.Now()
	_, err := entity.col.InsertOne(ctx, cat)
	return cat, err
}

func (entity *menuCategoryEntity) UpdateMenuCategoryById(id primitive.ObjectID, cat entities.MenuCategory) error {
	logrus.Info("UpdateMenuCategoryById")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cat.UpdatedDate = time.Now()
	_, err := entity.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{
		"name":        cat.Name,
		"sortOrder":   cat.SortOrder,
		"updatedBy":   cat.UpdatedBy,
		"updatedDate": cat.UpdatedDate,
	}})
	return err
}

func (entity *menuCategoryEntity) DeleteMenuCategoryById(id primitive.ObjectID) error {
	logrus.Info("DeleteMenuCategoryById")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := entity.col.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
