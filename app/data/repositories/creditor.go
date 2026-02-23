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

type creditorEntity struct {
	col    *mongo.Collection
	payCol *mongo.Collection
}

type ICreditor interface {
	GetCreditors(status string) ([]entities.Creditor, error)
	GetCreditorById(id primitive.ObjectID) (entities.Creditor, error)
	CreateCreditor(creditor entities.Creditor) (entities.Creditor, error)
	UpdateCreditor(id primitive.ObjectID, creditor entities.Creditor) error
	GetCreditorPayments(creditorId primitive.ObjectID) ([]entities.CreditorPayment, error)
	CreateCreditorPayment(payment entities.CreditorPayment) (entities.CreditorPayment, error)
}

func NewCreditorEntity(resource *db.Resource) ICreditor {
	col := resource.SnookDb.Collection("creditors")
	payCol := resource.SnookDb.Collection("creditor_payments")
	return &creditorEntity{col: col, payCol: payCol}
}

func (entity *creditorEntity) GetCreditors(status string) ([]entities.Creditor, error) {
	logrus.Info("GetCreditors")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{}
	if status != "" {
		filter["status"] = status
	}
	opts := options.Find().SetSort(bson.D{{Key: "createdDate", Value: -1}})
	cursor, err := entity.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var creditors []entities.Creditor
	if err = cursor.All(ctx, &creditors); err != nil {
		return nil, err
	}
	return creditors, nil
}

func (entity *creditorEntity) GetCreditorById(id primitive.ObjectID) (entities.Creditor, error) {
	logrus.Info("GetCreditorById")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var creditor entities.Creditor
	err := entity.col.FindOne(ctx, bson.M{"_id": id}).Decode(&creditor)
	return creditor, err
}

func (entity *creditorEntity) CreateCreditor(creditor entities.Creditor) (entities.Creditor, error) {
	logrus.Info("CreateCreditor")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	creditor.Id = primitive.NewObjectID()
	creditor.CreatedDate = time.Now()
	creditor.UpdatedDate = time.Now()
	_, err := entity.col.InsertOne(ctx, creditor)
	return creditor, err
}

func (entity *creditorEntity) UpdateCreditor(id primitive.ObjectID, creditor entities.Creditor) error {
	logrus.Info("UpdateCreditor")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	creditor.UpdatedDate = time.Now()
	_, err := entity.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{
		"paidAmount":  creditor.PaidAmount,
		"remaining":   creditor.Remaining,
		"status":      creditor.Status,
		"note":        creditor.Note,
		"updatedBy":   creditor.UpdatedBy,
		"updatedDate": creditor.UpdatedDate,
	}})
	return err
}

func (entity *creditorEntity) GetCreditorPayments(creditorId primitive.ObjectID) ([]entities.CreditorPayment, error) {
	logrus.Info("GetCreditorPayments")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	opts := options.Find().SetSort(bson.D{{Key: "createdDate", Value: -1}})
	cursor, err := entity.payCol.Find(ctx, bson.M{"creditorId": creditorId}, opts)
	if err != nil {
		return nil, err
	}
	var payments []entities.CreditorPayment
	if err = cursor.All(ctx, &payments); err != nil {
		return nil, err
	}
	return payments, nil
}

func (entity *creditorEntity) CreateCreditorPayment(payment entities.CreditorPayment) (entities.CreditorPayment, error) {
	logrus.Info("CreateCreditorPayment")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	payment.Id = primitive.NewObjectID()
	payment.CreatedDate = time.Now()
	_, err := entity.payCol.InsertOne(ctx, payment)
	return payment, err
}
