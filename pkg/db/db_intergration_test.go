package db_test

import (
	"math/rand"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/tarunKoyalwar/talosplus/pkg/db"
)

func Test_DBProvider(t *testing.T) {
	// Check if all methods of DB Provider are working correctly
	mongodb := os.Getenv("USE_MONGODB")

	dbname := randomString(6)

	if mongodb == "" {
		// Use bbolt embedded db
		db.UseBBoltDB(os.TempDir(), dbname, "talosplus")
		t.Logf("Using  BBoltDB as DB")
	}

	testcases := []struct {
		key      string
		value    string
		explicit bool
	}{
		{"pd", "subfinder", true},
		{"owasp", "amass", false},
		{"portswigger", "burpsuite", true},
	}

	// Test Database Provider call `Put`

	t.Logf("Loading testcases for Put Call")

	for _, tc := range testcases {
		er1 := db.DB.Put(tc.key, tc.value, tc.explicit)
		if er1 != nil {
			t.Errorf("failed to save variable to db %v for %v when %v", er1, db.DB.ProviderName(), tc)
		}

	}

	t.Logf("Test for Put Method Successful")

	// get test

	val, err := db.DB.Get("pd")

	if val != "subfinder" || err != nil {
		t.Errorf("failed to get value of pd expected `subfinder` but got %v with error %v", val, err)
	}

	val2, err2 := db.DB.Get("owasp")

	if val2 != "amass" || err2 != nil {
		t.Errorf("failed to get value of owasp expected `amass` but fot %v with error %v", val2, err2)
	}

	// Test for GetAllVarNames

	t.Logf("Loading Test for GetAllVarNames Method")

	expected_allvars := map[string]bool{
		"pd":          true,
		"owasp":       true,
		"portswigger": true,
	}

	allvars, erx := db.DB.GetAllVarNames()

	if !reflect.DeepEqual(allvars, expected_allvars) || erx != nil {
		t.Errorf("failed to getallvarNames expected %v but got %v with error %v", expected_allvars, allvars, erx)
	}

	t.Logf("Test for GetAllVarNames Successful")

	// Test for GetAllImplicit
	t.Logf("Loading Test for GetAllImplicit Method")

	expected_implicitvars := map[string]string{
		"owasp": "amass",
	}

	implicitvars, erx2 := db.DB.GetAllImplicit()

	if !reflect.DeepEqual(implicitvars, expected_implicitvars) || erx2 != nil {
		t.Errorf("failed to getimplicitvars expected %v but got %v with error %v", expected_implicitvars, implicitvars, erx2)
	}

	t.Logf("Test for GetAllImplicit Successful")

	// Test for GetAllExplicit

	t.Logf("Loading Test for GetAllExplicit Method")

	expected_explicitvars := map[string]string{
		"pd":          "subfinder",
		"portswigger": "burpsuite",
	}

	explicitvars, erx3 := db.DB.GetAllExplicit()

	if !reflect.DeepEqual(explicitvars, expected_explicitvars) || erx3 != nil {
		t.Errorf("failed to getallexplicitvars expected %v but got %v with error %v", expected_explicitvars, explicitvars, erx3)
	}

	t.Logf("Test for GetAllExplicit Successful")

	err = db.DB.Close()

	if err != nil {
		t.Errorf("failed to properly close db connection")
	}

}

// randomString : Generates random strings of given length with Uppercase charset
func randomString(size int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	var b strings.Builder
	for i := 0; i < size; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	str := b.String()

	return str
}
