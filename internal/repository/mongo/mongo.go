package mongorep

import (
	"context"
	"rental-server/internal/domain"
	"rental-server/internal/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBRepository struct {
	client   *mongo.Client
	Database string
}

func NewMongoDBRepository(uri string, database string) (*MongoDBRepository, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	return &MongoDBRepository{
		client:   client,
		Database: database,
	}, err
}

func (r *MongoDBRepository) Add(userId int64, object domain.RentObject) error {
	coll := r.client.Database(r.Database).Collection("objects")

	_, err := r.GetByName(userId, object.Name)
	if err == nil {
		return repository.ObjectAlreadyExists
	}

	if object.Records == nil {
		object.Records = []domain.Record{}
	}

	var data = struct {
		UserID            int64 `bson:"user_id"`
		domain.RentObject `bson:"rent_object"`
	}{
		userId,
		object,
	}

	_, err = coll.InsertOne(context.TODO(), data)
	return err
}

func (r *MongoDBRepository) Delete(userId int64, objectName string) error {
	_, err := r.GetByName(userId, objectName)
	if err == repository.ObjectNotFoundError {
		return repository.ObjectNotFoundError
	}

	coll := r.client.Database(r.Database).Collection("objects")

	filter := bson.D{
		{Key: "user_id", Value: userId},
		{Key: "rent_object.name", Value: objectName},
	}
	_, err = coll.DeleteOne(context.TODO(), filter)
	return err

}

func (r *MongoDBRepository) Update(userId int64, objectName string, input domain.UpdateRentObjectInput) error {
	obj, err := r.GetByName(userId, objectName)

	if err == repository.ObjectNotFoundError {
		return repository.ObjectNotFoundError
	}

	updated := obj.Update(input)
	coll := r.client.Database(r.Database).Collection("objects")

	filter := bson.D{
		{Key: "user_id", Value: userId},
		{Key: "rent_object.name", Value: objectName},
	}
	_, err = coll.UpdateOne(context.TODO(), filter, bson.D{{Key: "$set", Value: bson.D{{Key: "rent_object", Value: updated}}}})
	return err
}

func (r *MongoDBRepository) GetByName(userId int64, objectName string) (domain.RentObject, error) {
	coll := r.client.Database(r.Database).Collection("objects")

	var data struct {
		UserID            int64 `bson:"user_id"`
		domain.RentObject `bson:"rent_object"`
	}
	filter := bson.D{
		{Key: "user_id", Value: userId},
		{Key: "rent_object.name", Value: objectName},
	}
	err := coll.FindOne(context.TODO(), filter).Decode(&data)
	if err == mongo.ErrNoDocuments {
		return data.RentObject, repository.ObjectNotFoundError
	}

	return data.RentObject, err
}

func (r *MongoDBRepository) GetAll(userId int64) ([]domain.RentObject, error) {
	coll := r.client.Database(r.Database).Collection("objects")

	filter := bson.D{
		{Key: "user_id", Value: userId},
	}

	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var results []struct {
		UserID            int64 `bson:"user_id"`
		domain.RentObject `bson:"rent_object"`
	}
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}
	var objects = []domain.RentObject{}
	for _, res := range results {
		objects = append(objects, res.RentObject)
	}
	return objects, nil

}

func (r *MongoDBRepository) AddRecord(userID int64, objectName string, record domain.Record) (int, error) {
	obj, err := r.GetByName(userID, objectName)

	if err == repository.ObjectNotFoundError {
		return 0, repository.ObjectNotFoundError
	}

	coll := r.client.Database(r.Database).Collection("objects")
	filter := bson.D{
		{Key: "user_id", Value: userID},
		{Key: "rent_object.name", Value: objectName},
	}

	obj.AddRecord(record)

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "rent_object", Value: obj},
		},
		}}

	_, err = coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return 0, err
	}
	return len(obj.Records) - 1, nil
}

func (r *MongoDBRepository) DeleteRecord(userID int64, objectName string, recordIndex int) error {
	obj, err := r.GetByName(userID, objectName)
	if err == repository.ObjectNotFoundError {
		return repository.ObjectNotFoundError
	}

	err = obj.DeleteRecord(recordIndex)
	if err != nil {
		return err
	}

	coll := r.client.Database(r.Database).Collection("objects")

	filter := bson.D{
		{Key: "user_id", Value: userID},
		{Key: "rent_object.name", Value: objectName},
	}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "rent_object", Value: obj}}}}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	return err
}

func (r *MongoDBRepository) UpdateRecord(userID int64, objectName string, recordIndex int, input domain.UpdateRecordInput) error {
	obj, err := r.GetByName(userID, objectName)
	if err == repository.ObjectNotFoundError {
		return repository.ObjectNotFoundError
	}

	err = obj.UpdateRecord(recordIndex, input)
	if err != nil {
		return err
	}

	coll := r.client.Database(r.Database).Collection("objects")

	filter := bson.D{
		{Key: "user_id", Value: userID},
		{Key: "rent_object.name", Value: objectName},
	}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "rent_object", Value: obj}}}}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	return err
}

func (r *MongoDBRepository) GetRecordByIndex(userID int64, objectName string, recordIndex int) (domain.Record, error) {
	obj, err := r.GetByName(userID, objectName)
	if err == repository.ObjectNotFoundError {
		return domain.Record{}, repository.ObjectNotFoundError
	}
	return obj.GetRecordByIndex(recordIndex)
}

func (r *MongoDBRepository) GetAllRecords(userID int64, objectName string) ([]domain.Record, error) {
	obj, err := r.GetByName(userID, objectName)
	if err == repository.ObjectNotFoundError {
		return nil, repository.ObjectNotFoundError
	}
	return obj.GetAllRecords(), nil
}

func (r *MongoDBRepository) Clear() {
	r.client.Database(r.Database).Drop(context.TODO())
}
