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

type settingEntity struct {
	col *mongo.Collection
}

type ISetting interface {
	GetSetting() (entities.Setting, error)
	UpsertSetting(setting entities.Setting) error
}

func NewSettingEntity(resource *db.Resource) ISetting {
	col := resource.SnookDb.Collection("settings")
	return &settingEntity{col: col}
}

func (entity *settingEntity) GetSetting() (entities.Setting, error) {
	logrus.Info("GetSetting")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var setting entities.Setting
	err := entity.col.FindOne(ctx, bson.M{}).Decode(&setting)
	return setting, err
}

func (entity *settingEntity) UpsertSetting(setting entities.Setting) error {
	logrus.Info("UpsertSetting")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	setting.UpdatedDate = time.Now()
	opts := options.Update().SetUpsert(true)
	_, err := entity.col.UpdateOne(ctx, bson.M{}, bson.M{"$set": bson.M{
		"companyName":    setting.CompanyName,
		"companyAddress": setting.CompanyAddress,
		"companyPhone":   setting.CompanyPhone,
		"companyTaxId":   setting.CompanyTaxId,
		"receiptFooter":  setting.ReceiptFooter,
		"promptPayId":    setting.PromptPayId,
		"updatedBy":      setting.UpdatedBy,
		"updatedDate":    setting.UpdatedDate,
	}, "$setOnInsert": bson.M{"_id": primitive.NewObjectID()}}, opts)
	return err
}
