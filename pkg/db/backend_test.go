package db_test

import (
	"testing"

	"github.com/tarunKoyalwar/talosplus/pkg/db"
	"github.com/tarunKoyalwar/talosplus/pkg/db/mongox"
)

func Test_MongoDBCompatibility(t *testing.T) {
	var _ db.Provider = (*mongox.Provider)(nil)
}
