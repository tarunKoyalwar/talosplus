package shared

var DefaultRegistry *Registry = NewRegistry()

// Registry : Struct Containing all info that is shared b/w cmds
// This includes Local Vars, Exports etc and a lot of other stuff
type Registry struct {
	VarAddress map[string][]string // Address of commands that export data
	FoundVars  map[string]bool     // All Implicit Declarations
	Dependents map[string][]string // CMD's that depend on a particular variable
	// Dependents[Variable][cmd uid's]
}

// Registerdep : Registers uid and variable to registry
func (r *Registry) Registerdep(uid string, variable string) {

	//check if variable was explicitly declared at start of script
	//if explicit no need to mark as dep

	if SharedVars.IsExplicitVar(variable) {
		return
	}

	if r.Dependents[variable] == nil {
		r.Dependents[variable] = []string{}
	}

	// and uid to  dependents list
	r.Dependents[variable] = append(r.Dependents[variable], uid)

	//check if variable was declared already
	if !SharedVars.Exists(variable) {
		//if ok is false then variable is not declared
		if !r.FoundVars[variable] {
			//Also not a implicit / runtime declaration
			r.FoundVars[variable] = false // mark this as not declared
			// if found later will be marked as true then
		}
	}

}

// Registerexp : Registers uid as exporter of given variable
func (r *Registry) RegisterExport(uid string, variable string) {

	// add address of variable here address is uid
	if r.VarAddress[variable] == nil {
		r.VarAddress[variable] = []string{uid}
	} else {
		r.VarAddress[variable] = append(r.VarAddress[variable], uid)
	}

	//also mark variable as implicitly declared
	r.FoundVars[variable] = true
}

func NewRegistry() *Registry {
	r := &Registry{
		VarAddress: map[string][]string{},
		FoundVars:  map[string]bool{},
		Dependents: map[string][]string{},
	}

	return r

}
