package mongox

/*
Private or Internal Methods for CRUD Operations
*/

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type dbEntry struct {
	Varname  string `bson:"varname"`
	Value    string `bson:"value"`
	Explicit bool   `bson:"explicit"`
}

//lint:file-ignore U1000 Ignore all unused code, it's for future use only

// GetDatabase : Connect to Database
func (m *Provider) GetDatabase(name string) {
	m.db = m.client.Database(name)
}

// GetCollection : Get COllection
func (m *Provider) GetCollection(name string) {
	m.collection = m.db.Collection(name)
}

// Isconnected : Check If Connected to Database
func (m *Provider) Isconnected() bool {
	if m.db == nil {
		return false
	} else {
		return true
	}
}

// Connect : Connect to database
func (m *Provider) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if m.MongoDB.URL == "" {
		m.MongoDB.URL = "mongodb://localhost:27017"
	}
	var err error
	m.client, err = mongo.Connect(ctx, options.Client().ApplyURI(m.MongoDB.URL))
	if err != nil {
		return err
	}

	//add ping here
	err = m.client.Ping(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

// PingTest : Test Successful Connection
func (m *Provider) PingTest() error {
	if m.client == nil {
		er := m.Connect()
		if er != nil {
			return er
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := m.client.Ping(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

// Disconnect : throws error when fails
func (m *Provider) Disconnect() error {
	er := m.PingTest()
	if er != nil {
		return er
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := m.client.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}

// ListDatabases : List all Databases
func (m *Provider) ListDatabases() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	results, err := m.client.ListDatabaseNames(ctx, bson.D{})
	if err != nil {
		return []string{}, err
	}
	return results, nil
}

// ListDBCollections : List All Collections of Current Database
func (m *Provider) ListDBCollections() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	results, err := m.db.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return []string{}, err
	}
	return results, nil
}

// CreateCollection : Creates New Collection
func (m *Provider) CreateCollection(collname string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := m.db.CreateCollection(ctx, collname)
	return err
}

// UpdateDocument : Update Existing Document or create New One
func (m *Provider) UpdateDocument(filter interface{}, dat interface{}) (*mongo.UpdateResult, error) {
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
func (m *Provider) FindAll() ([]bson.D, error) {
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
func (m *Provider) FindWhere(filter interface{}) ([]bson.D, error) {
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
func (m *Provider) InsertOne(dat interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := m.collection.InsertOne(ctx, dat)

	return err
}

// FindOne : Find One Document using Filter
func (m *Provider) FindOne(filter interface{}, data interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := m.collection.FindOne(ctx, filter).Decode(data); err != nil {
		return err
	}

	return nil

}
