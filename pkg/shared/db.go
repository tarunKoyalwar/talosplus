package shared

import (
	"github.com/tarunKoyalwar/talosplus/pkg/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

type DBEntry struct {
	Varname  string `bson:"varname"`
	Value    string `bson:"value"`
	Explicit bool   `bson:"explicit"`
}

func SavetoDB(key string, value string, explicit bool) error {
	if mongodb.MDB == nil {
		return nil
	}

	x := DBEntry{
		Varname:  key,
		Value:    value,
		Explicit: explicit,
	}

	filter := bson.M{"varname": key}
	data := bson.M{"$set": x}

	_, err := mongodb.MDB.UpdateDocument(filter, data)

	return err

}

func GetFromDB(key string) (string, error) {

	if mongodb.MDB == nil {
		return "", nil
	}

	var x DBEntry
	filter := bson.M{"varname": key}

	err := mongodb.MDB.FindOne(filter, &x)

	return x.Value, err
}

func LoadAllExplicitVars() (map[string]string, error) {

	if mongodb.MDB == nil {
		return map[string]string{}, nil
	}

	filter := bson.M{"explicit": true}

	resp, err := mongodb.MDB.FindWhere(filter)

	if err != nil {
		return map[string]string{}, err
	}

	res := map[string]string{}

	for _, v := range resp {

		var zx DBEntry
		bin, _ := bson.Marshal(v)

		bson.Unmarshal(bin, &zx)

		res[zx.Varname] = zx.Value

	}

	return res, nil

}

func LoadAllRuntimeVars() (map[string]string, error) {
	if mongodb.MDB == nil {
		return map[string]string{}, nil
	}

	filter := bson.M{"explicit": false}

	resp, err := mongodb.MDB.FindWhere(filter)

	if err != nil {
		return map[string]string{}, err
	}

	res := map[string]string{}

	for _, v := range resp {

		var zx DBEntry
		bin, _ := bson.Marshal(v)

		bson.Unmarshal(bin, &zx)

		res[zx.Varname] = zx.Value

	}

	return res, nil
}

func GetAllVarNames() (map[string]bool, error) {

	if mongodb.MDB == nil {
		return map[string]bool{}, nil
	}

	resp, err := mongodb.MDB.FindWhere(bson.M{})
	if err != nil {
		return map[string]bool{}, err
	}

	res := map[string]bool{}

	for _, v := range resp {

		var zx DBEntry
		bin, _ := bson.Marshal(v)

		bson.Unmarshal(bin, &zx)

		res[zx.Varname] = zx.Explicit

	}

	return res, nil
}
