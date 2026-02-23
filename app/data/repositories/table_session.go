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

type tableSessionEntity struct {
	col *mongo.Collection
}

type ITableSession interface {
	GetTableSessions(startDate, endDate time.Time) ([]entities.TableSession, error)
	GetTableSessionById(id primitive.ObjectID) (entities.TableSession, error)
	GetActiveSessionByTableId(tableId primitive.ObjectID) (entities.TableSession, error)
	CreateTableSession(session entities.TableSession) (entities.TableSession, error)
	UpdateTableSession(id primitive.ObjectID, session entities.TableSession) error
	GetSessionSummary(startDate, endDate time.Time) (entities.SessionSummary, error)
	GetSessionDailyChart(startDate, endDate time.Time) ([]entities.SessionDailyChart, error)
	GetSessionsByTableId(tableId primitive.ObjectID, startDate, endDate time.Time) ([]entities.TableSession, error)
}

func NewTableSessionEntity(resource *db.Resource) ITableSession {
	col := resource.SnookDb.Collection("table_sessions")
	return &tableSessionEntity{col: col}
}

func (entity *tableSessionEntity) GetTableSessions(startDate, endDate time.Time) ([]entities.TableSession, error) {
	logrus.Info("GetTableSessions")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{"startTime": bson.M{"$gte": startDate, "$lte": endDate}}
	opts := options.Find().SetSort(bson.D{{Key: "startTime", Value: -1}})
	cursor, err := entity.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var sessions []entities.TableSession
	if err = cursor.All(ctx, &sessions); err != nil {
		return nil, err
	}
	return sessions, nil
}

func (entity *tableSessionEntity) GetTableSessionById(id primitive.ObjectID) (entities.TableSession, error) {
	logrus.Info("GetTableSessionById")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var session entities.TableSession
	err := entity.col.FindOne(ctx, bson.M{"_id": id}).Decode(&session)
	return session, err
}

func (entity *tableSessionEntity) GetActiveSessionByTableId(tableId primitive.ObjectID) (entities.TableSession, error) {
	logrus.Info("GetActiveSessionByTableId")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var session entities.TableSession
	err := entity.col.FindOne(ctx, bson.M{"tableId": tableId, "status": bson.M{"$in": []string{"ACTIVE", "PAUSED"}}}).Decode(&session)
	return session, err
}

func (entity *tableSessionEntity) CreateTableSession(session entities.TableSession) (entities.TableSession, error) {
	logrus.Info("CreateTableSession")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	session.Id = primitive.NewObjectID()
	session.CreatedDate = time.Now()
	session.UpdatedDate = time.Now()
	_, err := entity.col.InsertOne(ctx, session)
	return session, err
}

func (entity *tableSessionEntity) UpdateTableSession(id primitive.ObjectID, session entities.TableSession) error {
	logrus.Info("UpdateTableSession")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	session.UpdatedDate = time.Now()
	_, err := entity.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{
		"status":            session.Status,
		"endTime":           session.EndTime,
		"pausedAt":          session.PausedAt,
		"totalPausedMins":   session.TotalPausedMins,
		"durationMins":      session.DurationMins,
		"tableCharge":       session.TableCharge,
		"foodTotal":         session.FoodTotal,
		"discount":          session.Discount,
		"promotionId":       session.PromotionId,
		"promotionName":     session.PromotionName,
		"promotionDiscount": session.PromotionDiscount,
		"grandTotal":        session.GrandTotal,
		"note":              session.Note,
		"tableId":           session.TableId,
		"tableName":         session.TableName,
		"updatedBy":         session.UpdatedBy,
		"updatedDate":       session.UpdatedDate,
	}})
	return err
}

func (entity *tableSessionEntity) GetSessionSummary(startDate, endDate time.Time) (entities.SessionSummary, error) {
	logrus.Info("GetSessionSummary")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"startTime": bson.M{"$gte": startDate, "$lte": endDate}, "status": "CLOSED"}}},
		{{Key: "$group", Value: bson.M{
			"_id":           nil,
			"totalSessions": bson.M{"$sum": 1},
			"totalRevenue":  bson.M{"$sum": "$grandTotal"},
			"totalTable":    bson.M{"$sum": "$tableCharge"},
			"totalFood":     bson.M{"$sum": "$foodTotal"},
		}}},
	}
	cursor, err := entity.col.Aggregate(ctx, pipeline)
	if err != nil {
		return entities.SessionSummary{}, err
	}
	var results []entities.SessionSummary
	if err = cursor.All(ctx, &results); err != nil {
		return entities.SessionSummary{}, err
	}
	if len(results) == 0 {
		return entities.SessionSummary{}, nil
	}
	return results[0], nil
}

func (entity *tableSessionEntity) GetSessionDailyChart(startDate, endDate time.Time) ([]entities.SessionDailyChart, error) {
	logrus.Info("GetSessionDailyChart")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"startTime": bson.M{"$gte": startDate, "$lte": endDate}, "status": "CLOSED"}}},
		{{Key: "$group", Value: bson.M{
			"_id":           bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$startTime"}},
			"totalSessions": bson.M{"$sum": 1},
			"totalRevenue":  bson.M{"$sum": "$grandTotal"},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "_id", Value: 1}}}},
	}
	cursor, err := entity.col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	var results []entities.SessionDailyChart
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (entity *tableSessionEntity) GetSessionsByTableId(tableId primitive.ObjectID, startDate, endDate time.Time) ([]entities.TableSession, error) {
	logrus.Info("GetSessionsByTableId")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{"tableId": tableId, "startTime": bson.M{"$gte": startDate, "$lte": endDate}}
	opts := options.Find().SetSort(bson.D{{Key: "startTime", Value: -1}})
	cursor, err := entity.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var sessions []entities.TableSession
	if err = cursor.All(ctx, &sessions); err != nil {
		return nil, err
	}
	return sessions, nil
}
