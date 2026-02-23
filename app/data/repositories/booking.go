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

type bookingEntity struct {
	col *mongo.Collection
}

type IBooking interface {
	GetBookings(startDate, endDate time.Time) ([]entities.Booking, error)
	GetBookingById(id primitive.ObjectID) (entities.Booking, error)
	CreateBooking(booking entities.Booking) (entities.Booking, error)
	UpdateBookingById(id primitive.ObjectID, booking entities.Booking) error
	UpdateBookingStatus(id primitive.ObjectID, status string) error
	DeleteBookingById(id primitive.ObjectID) error
}

func NewBookingEntity(resource *db.Resource) IBooking {
	col := resource.SnookDb.Collection("bookings")
	return &bookingEntity{col: col}
}

func (entity *bookingEntity) GetBookings(startDate, endDate time.Time) ([]entities.Booking, error) {
	logrus.Info("GetBookings")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{"bookingDate": bson.M{"$gte": startDate, "$lte": endDate}}
	opts := options.Find().SetSort(bson.D{{Key: "bookingDate", Value: -1}})
	cursor, err := entity.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var bookings []entities.Booking
	if err = cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

func (entity *bookingEntity) GetBookingById(id primitive.ObjectID) (entities.Booking, error) {
	logrus.Info("GetBookingById")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var booking entities.Booking
	err := entity.col.FindOne(ctx, bson.M{"_id": id}).Decode(&booking)
	return booking, err
}

func (entity *bookingEntity) CreateBooking(booking entities.Booking) (entities.Booking, error) {
	logrus.Info("CreateBooking")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	booking.Id = primitive.NewObjectID()
	booking.CreatedDate = time.Now()
	booking.UpdatedDate = time.Now()
	_, err := entity.col.InsertOne(ctx, booking)
	return booking, err
}

func (entity *bookingEntity) UpdateBookingById(id primitive.ObjectID, booking entities.Booking) error {
	logrus.Info("UpdateBookingById")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	booking.UpdatedDate = time.Now()
	_, err := entity.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{
		"tableId":       booking.TableId,
		"tableName":     booking.TableName,
		"customerName":  booking.CustomerName,
		"customerPhone": booking.CustomerPhone,
		"bookingDate":   booking.BookingDate,
		"startTime":     booking.StartTime,
		"endTime":       booking.EndTime,
		"note":          booking.Note,
		"updatedBy":     booking.UpdatedBy,
		"updatedDate":   booking.UpdatedDate,
	}})
	return err
}

func (entity *bookingEntity) UpdateBookingStatus(id primitive.ObjectID, status string) error {
	logrus.Info("UpdateBookingStatus")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := entity.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"status": status, "updatedDate": time.Now()}})
	return err
}

func (entity *bookingEntity) DeleteBookingById(id primitive.ObjectID) error {
	logrus.Info("DeleteBookingById")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := entity.col.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
