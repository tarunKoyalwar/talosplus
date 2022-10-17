package db

// DBProvider : DB Provider i.e MongoDB , InMemory or bbolt(embedded) etc
type DBProvider int

const (
	Mongo_DB DBProvider = iota
	BBolt_DB
)

// Provider : Interface for All Providers of Database
type Provider interface {
	// Put variable name and value to db
	Put(key, value string, isExplicit bool) error
	// Get value of a variable
	Get(key string) (string, error)
	// GetAll Variable Names
	GetAllVarNames() (map[string]bool, error)
	// GetAllExplicit Variable Names
	GetAllExplicit() (map[string]string, error)
	// GetAllRuntime Variable Names
	GetAllImplicit() (map[string]string, error)
	// // GetAllGlobal Variables
	// GetAllGlobal() (map[string]string, error)
	// ProviderName
	ProviderName() string
	// Close DB Connection
	Close() error
}
