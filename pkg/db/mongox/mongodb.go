package mongox

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Options : MongoDB Options
type Options struct {
	URL            string // MongoDB Connection URL
	DBName         string
	CollectionName string
}

// Provider : MongoDB Provider
type Provider struct {
	MongoDB    *Options
	client     *mongo.Client     // MongoDB Connection Client
	db         *mongo.Database   // MongoDB Database
	collection *mongo.Collection // MongoDB Collection
}

// New Provider Instance
func New(opts *Options) (*Provider, error) {
	p := &Provider{
		MongoDB: opts,
	}

	err := p.Open()

	return p, err
}

// validate DB Connection
func (m *Provider) validate() (string, error) {
	if m.MongoDB == nil {
		return "No DB Options", fmt.Errorf("mongodb options missing")
	}
	if m.db == nil {
		return "No DB Connection", fmt.Errorf("mongodb no connection")
	}
	return "", nil
}

// Open DB Connection
func (m *Provider) Open() error {

	if err := m.Connect(); err != nil {
		return err
	}

	if m.MongoDB == nil {
		return fmt.Errorf("mongodb options missing")
	}

	if m.MongoDB.DBName == "" {
		return fmt.Errorf("db name missing")
	} else {
		m.GetDatabase(m.MongoDB.DBName)
	}

	if m.MongoDB.CollectionName == "" {
		return fmt.Errorf("mongoDB collection name missing")
	}

	return nil

}

// Close DB Connection
func (m *Provider) Close() error {
	return m.Disconnect()
}

// Get variable value
func (m *Provider) Get(key string) (string, error) {

	if msg, err := m.validate(); err != nil {
		return msg, err
	}

	var x dbEntry
	filter := bson.M{"varname": key}

	err := m.FindOne(filter, &x)

	return x.Value, err
}

// Put Variable to DB
func (m *Provider) Put(key, value string, isExplicit bool) error {

	if _, err := m.validate(); err != nil {
		return err
	}

	x := dbEntry{
		Varname:  key,
		Value:    value,
		Explicit: isExplicit,
	}

	filter := bson.M{"varname": key}
	data := bson.M{"$set": x}

	_, err := m.UpdateDocument(filter, data)

	return err
}

// GetAllVarNames
func (m *Provider) GetAllVarNames() (map[string]bool, error) {

	if _, err := m.validate(); err != nil {
		return map[string]bool{}, err
	}

	resp, err := m.FindWhere(bson.M{})
	if err != nil {
		return map[string]bool{}, err
	}

	res := map[string]bool{}

	for _, v := range resp {

		var zx dbEntry
		bin, _ := bson.Marshal(v)

		bson.Unmarshal(bin, &zx)

		res[zx.Varname] = zx.Explicit

	}

	return res, nil
}

// GetAllExplicit Variables and their Values
func (m *Provider) GetAllExplicit() (map[string]string, error) {

	if _, err := m.validate(); err != nil {
		return map[string]string{}, err
	}

	filter := bson.M{"explicit": true}

	resp, err := m.FindWhere(filter)

	if err != nil {
		return map[string]string{}, err
	}

	res := map[string]string{}

	for _, v := range resp {

		var zx dbEntry
		bin, _ := bson.Marshal(v)

		bson.Unmarshal(bin, &zx)

		res[zx.Varname] = zx.Value

	}

	return res, nil
}

// GetAllImplicit Variables and Values
func (m *Provider) GetAllImplicit() (map[string]string, error) {

	if _, err := m.validate(); err != nil {
		return map[string]string{}, err
	}

	filter := bson.M{"explicit": false}

	resp, err := m.FindWhere(filter)

	if err != nil {
		return map[string]string{}, err
	}

	res := map[string]string{}

	for _, v := range resp {

		var zx dbEntry
		bin, _ := bson.Marshal(v)

		bson.Unmarshal(bin, &zx)

		res[zx.Varname] = zx.Value

	}

	return res, nil
}

// ProviderName i.e DB backend here (mongodb)
func (m *Provider) ProviderName() string {
	return "MongoDB"
}
