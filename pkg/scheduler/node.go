package scheduler

// Node : Scheduler Node
type Node struct {
	Root     []*Node
	UID      string
	Comment  string
	Children []*Node
}
