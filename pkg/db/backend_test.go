package db_test

import (
	"testing"

	"github.com/tarunKoyalwar/talosplus/pkg/db"
	"github.com/tarunKoyalwar/talosplus/pkg/db/boltx"
	"github.com/tarunKoyalwar/talosplus/pkg/db/mongox"
)

func Test_MongoDBCompatibility(t *testing.T) {
	var _ db.Provider = (*mongox.Provider)(nil)
}

func Test_BBoltDBCompatibility(t *testing.T) {
	var _ db.Provider = (*boltx.Provider)(nil)
}
