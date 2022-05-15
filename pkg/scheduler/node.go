package scheduler

type Node struct {
	Root     []*Node
	UID      string
	Comment  string
	Children []*Node
}
