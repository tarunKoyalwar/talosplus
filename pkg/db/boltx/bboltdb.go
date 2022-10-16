package boltx

import (
	"fmt"
	"path"

	"go.etcd.io/bbolt"
)

// Options : BBolt DB
type Options struct {
	Directory  string
	Filename   string
	BucketName string
}

// Provider of BBolt DB
type Provider struct {
	BBoltDB *Options  /// DB Options
	db      *bbolt.DB // DB pointer
}

// New BBoltDB Instance
func New(options *Options) (*Provider, error) {
	p := &Provider{
		BBoltDB: options,
	}

	er := p.open()

	return p, er
}

// validate connection
func (p *Provider) validate() (string, error) {
	if p.BBoltDB.BucketName == "" {
		return "DB Options Missing", fmt.Errorf("bboltdb bucket name not provider")
	}

	if p.db == nil {
		return "DB Instance Missing", fmt.Errorf("bbolt db no connection")
	}

	return "", nil
}

// Open DB Connection
func (p *Provider) open() error {
	var er error
	var dbpath string

	if p.BBoltDB.Directory == "." || p.BBoltDB.Directory == "" {
		dbpath = p.BBoltDB.Filename
	} else {
		dbpath = path.Join(p.BBoltDB.Directory, p.BBoltDB.Filename)
	}

	p.db, er = bbolt.Open(dbpath, 0644, nil)

	if er != nil {
		return er
	}

	p.db.Update(
		func(tx *bbolt.Tx) error {
			_, er := tx.CreateBucketIfNotExists([]byte(p.BBoltDB.BucketName))
			if er != nil {
				return er
			}
			return nil
		},
	)

	return nil
}

// Close DB Connection
func (p *Provider) Close() error {
	return p.db.Close()
}

// Get Variable Value
func (p *Provider) Get(key string) (string, error) {
	if _, er := p.validate(); er != nil {
		return "", er
	}

	var item *dbEntry
	var erx error

	p.db.View(
		func(tx *bbolt.Tx) error {
			bucket := tx.Bucket([]byte(p.BBoltDB.BucketName))

			val := bucket.Get([]byte(key))

			if val == nil {
				erx = fmt.Errorf("value of variable is missing")
				return nil
			}

			item, erx = newdbEntry(val)

			return nil
		},
	)

	if erx != nil || item == nil {
		return "", fmt.Errorf("bbolt db get failed value of key missing")
	}

	return item.Value, nil
}

// Put Variable to DB
func (p *Provider) Put(key, value string, isExplicit bool) error {
	item := dbEntry{
		Value:      value,
		IsExplicit: isExplicit,
	}

	itembin, er := item.getBytes()
	if er != nil {
		return er
	}

	if _, er := p.validate(); er != nil {
		return er
	}

	var erx error

	p.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(p.BBoltDB.BucketName))
		if err != nil {
			erx = err
			return err
		}

		erx = bucket.Put([]byte(key), itembin)
		return nil
	},
	)

	return erx

}

// GetAllVarNames
func (p *Provider) GetAllVarNames() (map[string]bool, error) {
	allvars := map[string]bool{}

	if _, er := p.validate(); er != nil {
		return allvars, er
	}

	var erx error

	p.db.View(
		func(tx *bbolt.Tx) error {
			bucket := tx.Bucket([]byte(p.BBoltDB.BucketName))

			erx = bucket.ForEach(func(k, v []byte) error {
				allvars[string(k)] = true
				return nil
			})
			return nil
		},
	)

	if erx != nil {
		return allvars, fmt.Errorf("bbolt db get failed value of key missing")
	}

	return allvars, nil
}

// GetAllExplicit Variables and their Values
func (p *Provider) GetAllExplicit() (map[string]string, error) {
	return p.getkvwhere(true)
}

// GetAllImplicit Variables and Values
func (p *Provider) GetAllImplicit() (map[string]string, error) {
	return p.getkvwhere(false)
}

// ProviderName i.e DB backend here (mongodb)
func (m *Provider) ProviderName() string {
	return "BBoltDB"
}
