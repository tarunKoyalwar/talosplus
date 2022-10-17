package db

import (
	"github.com/tarunKoyalwar/talosplus/pkg/db/boltx"
	"github.com/tarunKoyalwar/talosplus/pkg/db/mongox"
)

// DB
var DB Provider

// UserMongoDB as db backend
func UseMongoDB(url, dbname, collname string) error {

	options := mongox.Options{
		URL:            url,
		DBName:         dbname,
		CollectionName: collname,
	}

	var err error
	DB, err = mongox.New(&options)

	return err
}

// UseBBoltDB as db backend
func UseBBoltDB(dir, filename, bucketname string) error {
	options := boltx.Options{
		Directory:  dir,
		Filename:   filename,
		BucketName: bucketname,
	}

	var err error
	DB, err = boltx.New(&options)

	return err
}
