package db

import "github.com/tarunKoyalwar/talosplus/pkg/db/mongox"

// DB
var DB Provider

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
