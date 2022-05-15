package ioutils

type CSave struct {
	Comment  string `bson:"comment"`
	UID      string `bson:"uid"`
	Output   string `bson:"output"`
	CacheKey string `bson:"cachehash"`
}

type AllSave struct {
	Exports  map[string]string `bson:"exports"`
	Commands []CSave           `bson:"commands"`
}
