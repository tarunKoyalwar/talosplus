package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CRUD Operations for MongoDB Struct

// ListDatabases : List all Databases
func (m *MongoDB) ListDatabases() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	results, err := m.client.ListDatabaseNames(ctx, bson.D{})
	if err != nil {
		return []string{}, err
	}
	return results, nil
}

// ListDBCollections : List All Collections of Current Database
func (m *MongoDB) ListDBCollections() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	results, err := m.db.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return []string{}, err
	}
	return results, nil
}

// CreateCollection : Creates New Collection
func (m *MongoDB) CreateCollection(collname string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := m.db.CreateCollection(ctx, collname)
	return err
}

// UpdateDocument : Update Existing Document or create New One
func (m *MongoDB) UpdateDocument(filter interface{}, dat interface{}) (*mongo.UpdateResult, error) {
	// fmt.Println("Update Document called")
	opts := options.Update().SetUpsert(true)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := m.collection.UpdateOne(ctx, filter, dat, opts)
	if err != nil {
		return result, err
	}
	// fmt.Printf("Matchedcount %v , Modified COunt %v, with upsert id %v\n ", result.MatchedCount, result.ModifiedCount, result.UpsertedID)
	return result, nil
}

// FindAll : Find All Possible Matches
func (m *MongoDB) FindAll() ([]bson.D, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var arr []bson.D

	cursor, err := m.collection.Find(ctx, bson.M{})
	if err != nil {
		return arr, err
	}

	if err = cursor.All(ctx, &arr); err != nil {
		return arr, err
	}

	return arr, nil
}

// FindWhere : Find All Possible Matches Where
func (m *MongoDB) FindWhere(filter interface{}) ([]bson.D, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var arr []bson.D

	cursor, err := m.collection.Find(ctx, filter)
	if err != nil {
		return arr, err
	}

	if err = cursor.All(ctx, &arr); err != nil {
		return arr, err
	}

	return arr, nil
}

// InsertOne : Insert One Document
func (m *MongoDB) InsertOne(dat interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := m.collection.InsertOne(ctx, dat)

	return err
}

// FindOne : Find One Document using Filter
func (m *MongoDB) FindOne(filter interface{}, data interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := m.collection.FindOne(ctx, filter).Decode(data); err != nil {
		return err
	}

	return nil

}
