package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MDB : MongoDB Struct Instance
var MDB *MongoDB

// MongoDB : Wrapper around basic api
type MongoDB struct {
	client     *mongo.Client     // Actual Client
	URL        string            // MongoDB Connection String (default Local)
	db         *mongo.Database   // Pointer to MongoDB Database
	collection *mongo.Collection // Pointer to MongoDB Collection
}

// NewMongoDB : New Instance
func NewMongoDB(URL string) (*MongoDB, error) {

	z := &MongoDB{}
	z.URL = URL

	er := z.Connect()

	return z, er

}

// GetDatabase : Connect to Database
func (m *MongoDB) GetDatabase(name string) {
	m.db = m.client.Database(name)
}

// GetCollection : Get COllection
func (m *MongoDB) GetCollection(name string) {
	m.collection = m.db.Collection(name)
}

// Isconnected : Check If Connected to Database
func (m *MongoDB) Isconnected() bool {
	if m.db == nil {
		return false
	} else {
		return true
	}
}

// Connect : Connect to database
func (m *MongoDB) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if m.URL == "" {
		m.URL = "mongodb://localhost:27017"
	}
	var err error
	m.client, err = mongo.Connect(ctx, options.Client().ApplyURI(m.URL))
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
func (m MongoDB) PingTest() error {
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
func (m *MongoDB) Disconnect() {
	er := m.PingTest()
	if er != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := m.client.Disconnect(ctx); err != nil {
		panic(err)
	}
}
