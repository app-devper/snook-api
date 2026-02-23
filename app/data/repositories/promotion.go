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

type promotionEntity struct {
	col *mongo.Collection
}

type IPromotion interface {
	GetPromotions() ([]entities.Promotion, error)
	GetPromotionById(id primitive.ObjectID) (entities.Promotion, error)
	GetActivePromotions(tableType string) ([]entities.Promotion, error)
	CreatePromotion(promo entities.Promotion) (entities.Promotion, error)
	UpdatePromotionById(id primitive.ObjectID, promo entities.Promotion) error
	DeletePromotionById(id primitive.ObjectID) error
}

func NewPromotionEntity(resource *db.Resource) IPromotion {
	col := resource.SnookDb.Collection("promotions")
	return &promotionEntity{col: col}
}

func (entity *promotionEntity) GetPromotions() ([]entities.Promotion, error) {
	logrus.Info("GetPromotions")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	opts := options.Find().SetSort(bson.D{{Key: "createdDate", Value: -1}})
	cursor, err := entity.col.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	var promos []entities.Promotion
	if err = cursor.All(ctx, &promos); err != nil {
		return nil, err
	}
	return promos, nil
}

func (entity *promotionEntity) GetPromotionById(id primitive.ObjectID) (entities.Promotion, error) {
	logrus.Info("GetPromotionById")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var promo entities.Promotion
	err := entity.col.FindOne(ctx, bson.M{"_id": id}).Decode(&promo)
	return promo, err
}

func (entity *promotionEntity) GetActivePromotions(tableType string) ([]entities.Promotion, error) {
	logrus.Info("GetActivePromotions")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	now := time.Now()
	filter := bson.M{
		"status":    "ACTIVE",
		"startDate": bson.M{"$lte": now},
		"endDate":   bson.M{"$gte": now},
	}
	if tableType != "" {
		filter["$or"] = []bson.M{
			{"tableTypes": tableType},
			{"tableTypes": bson.M{"$exists": false}},
			{"tableTypes": nil},
			{"tableTypes": bson.A{}},
		}
	}
	cursor, err := entity.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var promos []entities.Promotion
	if err = cursor.All(ctx, &promos); err != nil {
		return nil, err
	}
	return promos, nil
}

func (entity *promotionEntity) CreatePromotion(promo entities.Promotion) (entities.Promotion, error) {
	logrus.Info("CreatePromotion")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	promo.Id = primitive.NewObjectID()
	promo.CreatedDate = time.Now()
	promo.UpdatedDate = time.Now()
	_, err := entity.col.InsertOne(ctx, promo)
	return promo, err
}

func (entity *promotionEntity) UpdatePromotionById(id primitive.ObjectID, promo entities.Promotion) error {
	logrus.Info("UpdatePromotionById")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	promo.UpdatedDate = time.Now()
	_, err := entity.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{
		"name":        promo.Name,
		"description": promo.Description,
		"type":        promo.Type,
		"playHours":   promo.PlayHours,
		"freeHours":   promo.FreeHours,
		"discountPct": promo.DiscountPct,
		"discountAmt": promo.DiscountAmt,
		"tableTypes":  promo.TableTypes,
		"startDate":   promo.StartDate,
		"endDate":     promo.EndDate,
		"status":      promo.Status,
		"updatedBy":   promo.UpdatedBy,
		"updatedDate": promo.UpdatedDate,
	}})
	return err
}

func (entity *promotionEntity) DeletePromotionById(id primitive.ObjectID) error {
	logrus.Info("DeletePromotionById")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := entity.col.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
