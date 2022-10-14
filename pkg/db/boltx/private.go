package boltx

import (
	"encoding/json"

	"go.etcd.io/bbolt"
)

/*
Private/ Internal Methods of Provider
*/

// dbEntry of BBoltDB containing Metadata
type dbEntry struct {
	Value      string `json:"value"`
	IsExplicit bool   `json:"isExplicit"`
}

// getBytes
func (d *dbEntry) getBytes() ([]byte, error) {
	bin, err := json.MarshalIndent(d, "", "\t")

	return bin, err
}

// newdbEntry
func newdbEntry(bin []byte) (*dbEntry, error) {
	var x dbEntry

	err := json.Unmarshal(bin, &x)

	if err != nil {
		return nil, err
	}

	return &x, nil
}

// getkvwhere : get key value where explicit is bool
func (p *Provider) getkvwhere(Explicit bool) (map[string]string, error) {
	vars := map[string]string{}

	if _, er := p.validate(); er != nil {
		return vars, er
	}

	var erx error

	p.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(p.BBoltDB.BucketName))

		erx = bucket.ForEach(func(k, v []byte) error {
			item, err := newdbEntry(v)
			if err != nil {
				erx = err
				return err
			}
			if item.IsExplicit == Explicit {
				vars[string(k)] = item.Value
			}

			return nil
		})

		return nil
	},
	)

	if erx != nil {
		return vars, erx
	}

	return vars, nil

}
