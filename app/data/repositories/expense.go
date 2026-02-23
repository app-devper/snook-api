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

type expenseEntity struct {
	col *mongo.Collection
}

type IExpense interface {
	GetExpenses(startDate, endDate time.Time) ([]entities.Expense, error)
	GetExpenseById(id primitive.ObjectID) (entities.Expense, error)
	CreateExpense(expense entities.Expense) (entities.Expense, error)
	UpdateExpenseById(id primitive.ObjectID, expense entities.Expense) error
	DeleteExpenseById(id primitive.ObjectID) error
}

func NewExpenseEntity(resource *db.Resource) IExpense {
	col := resource.SnookDb.Collection("expenses")
	return &expenseEntity{col: col}
}

func (entity *expenseEntity) GetExpenses(startDate, endDate time.Time) ([]entities.Expense, error) {
	logrus.Info("GetExpenses")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{"date": bson.M{"$gte": startDate, "$lte": endDate}}
	opts := options.Find().SetSort(bson.D{{Key: "date", Value: -1}})
	cursor, err := entity.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var expenses []entities.Expense
	if err = cursor.All(ctx, &expenses); err != nil {
		return nil, err
	}
	return expenses, nil
}

func (entity *expenseEntity) GetExpenseById(id primitive.ObjectID) (entities.Expense, error) {
	logrus.Info("GetExpenseById")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var expense entities.Expense
	err := entity.col.FindOne(ctx, bson.M{"_id": id}).Decode(&expense)
	return expense, err
}

func (entity *expenseEntity) CreateExpense(expense entities.Expense) (entities.Expense, error) {
	logrus.Info("CreateExpense")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	expense.Id = primitive.NewObjectID()
	expense.CreatedDate = time.Now()
	expense.UpdatedDate = time.Now()
	_, err := entity.col.InsertOne(ctx, expense)
	return expense, err
}

func (entity *expenseEntity) UpdateExpenseById(id primitive.ObjectID, expense entities.Expense) error {
	logrus.Info("UpdateExpenseById")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	expense.UpdatedDate = time.Now()
	_, err := entity.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{
		"category":    expense.Category,
		"description": expense.Description,
		"amount":      expense.Amount,
		"date":        expense.Date,
		"updatedBy":   expense.UpdatedBy,
		"updatedDate": expense.UpdatedDate,
	}})
	return err
}

func (entity *expenseEntity) DeleteExpenseById(id primitive.ObjectID) error {
	logrus.Info("DeleteExpenseById")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := entity.col.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
