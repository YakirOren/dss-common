package db

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/yakiroren/dss-common/models"

	"go.mongodb.org/mongo-driver/mongo/readpref"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDataStore struct {
	FilesCollection *mongo.Collection
	Client          *mongo.Client
}

func GetConnectionString(username string, password string, address string) string {
	return fmt.Sprintf("mongodb://%s:%s@%s", username, password, address)
}

type MongoConfig struct {
	MongoUsername       string `env:",required,notEmpty"`
	MongoPassword       string `env:",required,notEmpty"`
	MongoURL            string `env:"MONGO_URL,required,notEmpty"`
	MongoFileCollection string `env:",required,notEmpty"`
	MongoDBName         string `env:",required,notEmpty"`
}

func NewMongoDataStore(config *MongoConfig) (*MongoDataStore, error) {
	connection := GetConnectionString(config.MongoUsername, config.MongoPassword, config.MongoURL)

	client, err := mongo.NewClient(options.Client().ApplyURI(connection))
	if err != nil {
		return nil, err
	}

	if err := client.Connect(context.TODO()); err != nil {
		return nil, err
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, err
	}

	mod := mongo.IndexModel{
		Keys: bson.M{
			"path": 1, // index in ascending order
		}, Options: nil,
	}

	FilesCollection := client.Database(config.MongoDBName).Collection(config.MongoFileCollection)

	if _, err := FilesCollection.Indexes().CreateOne(context.TODO(), mod); err != nil {
		return nil, err
	}

	return &MongoDataStore{
		FilesCollection: FilesCollection,
		Client:          client,
	}, nil
}

func (db *MongoDataStore) WriteFile(ctx context.Context, file models.FileMetadata) (string, error) {
	result, err := db.FilesCollection.InsertOne(ctx, file)
	if err != nil {
		return "", err
	}
	fileID := result.InsertedID.(primitive.ObjectID).Hex()

	return fileID, nil
}

func (db *MongoDataStore) AppendFragment(ctx context.Context, id string, fragment models.Fragment) error {
	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = db.FilesCollection.UpdateOne(
		ctx,
		bson.M{"_id": hex},
		bson.M{"$push": bson.M{"fragments": fragment}, "$inc": bson.M{"size": fragment.Size}},
	)

	return err
}

func (db *MongoDataStore) GetMetadataByID(ctx context.Context, id interface{}) (*models.FileMetadata, bool) {
	output := models.FileMetadata{}
	filter := bson.D{{"_id", id}}

	if err := db.FilesCollection.FindOne(ctx, filter).Decode(&output); err != nil {
		return nil, false
	}

	return &output, true
}

func (db *MongoDataStore) GetMetadataByPath(ctx context.Context, path string) (*models.FileMetadata, bool) {
	output := models.FileMetadata{}

	basePath := filepath.Dir(path)
	name := filepath.Base(path)

	filter := bson.D{{"path", basePath}, {"name", name}}

	if err := db.FilesCollection.FindOne(ctx, filter).Decode(&output); err != nil {
		return nil, false
	}

	return &output, true
}

func (db *MongoDataStore) UpdateField(ctx context.Context, id string, field string, value interface{}) error {
	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.D{{"_id", hex}}

	update := bson.D{{"$set", bson.D{{field, value}}}}

	_, err = db.FilesCollection.UpdateOne(ctx, filter, update)
	return err
}

func (db *MongoDataStore) ListFiles(ctx context.Context, path string) ([]models.FileMetadata, error) {
	var output []models.FileMetadata

	cursor, err := db.FilesCollection.Find(ctx, bson.D{{"path", path}})
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &output); err != nil {
		return nil, err
	}

	return output, nil
}
