package shared

import (
	"fmt"
	"strings"
	"sync"

	"github.com/tarunKoyalwar/talosplus/pkg/ioutils"
	"github.com/tarunKoyalwar/talosplus/pkg/mongodb"
)

var SharedVars *Shared = NewShared()

// Shared : This contains data that is shared b/w commands
type Shared struct {
	mutex     *sync.Mutex       // To avoid race conditons while sharing
	variables map[string]string // All variables + Runtime
	explicit  map[string]bool   // Explicitly declared Variable Names
}

func (e *Shared) AcquireLock() {
	e.mutex.Lock()
}

func (e *Shared) ReleaseLock() {
	e.mutex.Unlock()
}

// Exists : Check if Enviornemtn variable already exists
func (e *Shared) Exists(key string) bool {
	_, ok := e.variables[key]
	return ok
}

// IsExplicitVar : Check If Variable Was Explicitly declared
func (e *Shared) IsExplicitVar(key string) bool {
	_, ok := e.explicit[key]
	return ok
}

// Get : Self Explainatory
func (e *Shared) Get(key string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("empty key")
	}

	val := ""

	val = e.variables[key]

	if val == "" {
		return val, fmt.Errorf("empty value for key : %v", key)
	}

	return val, nil
}

// Set : Self Explainatory
func (e *Shared) Set(key string, value string, explicit bool) error {

	if key == "" {
		return fmt.Errorf("empty key while setting")
	}

	// Must Trim Spaces Before setting
	// to avoid inconsistencies
	value = strings.TrimSpace(value)

	if value == "" {
		return fmt.Errorf("empty value while setting key %v", key)
	}

	e.variables[key] = value
	errx := SavetoDB(key, value, explicit)
	if errx != nil {
		ioutils.Cout.PrintWarning("failed to save %v to db\n:%v", key, errx.Error())
	}

	if explicit {
		e.explicit[key] = true
	}

	return nil
}

// GetGlobalVars : Self Explainatory
func (e *Shared) GetGlobalVars() map[string]string {

	z := map[string]string{}

	for k := range e.explicit {
		z[k] = e.variables[k]
	}

	return z

}

func (e *Shared) AddGlobalVarsFromDB() {
	if mongodb.MDB == nil {
		return
	}

	// fmt.Printf("adding vars")

	z, err := LoadAllExplicitVars()
	if err != nil || len(z) == 0 {
		ioutils.Cout.PrintWarning("failed to add global vars %v", err)
		return
	}

	// fmt.Printf("total len %v\n", len(z))
	for k, v := range z {
		e.variables[k] = v
		e.explicit[k] = true
	}

}

func NewShared() *Shared {

	s := Shared{
		mutex:     &sync.Mutex{},
		variables: map[string]string{},
		explicit:  map[string]bool{},
	}

	return &s

}
