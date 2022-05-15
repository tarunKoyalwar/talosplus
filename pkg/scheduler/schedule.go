package scheduler

import (
	"github.com/tarunKoyalwar/talosplus/pkg/ioutils"
	"github.com/tarunKoyalwar/talosplus/pkg/shared"
)

type Scheduler struct {
	AllNodes    map[string]*Node
	ExecPyramid [][]*Node
	BlackListed map[string]bool
}

func (s *Scheduler) AddNode(uid string, comment string) {
	x := Node{
		Root:     []*Node{},
		UID:      uid,
		Comment:  comment,
		Children: []*Node{},
	}
	s.AllNodes[uid] = &x
}

func (s *Scheduler) Run() {
	//remove old assigned data or verify if parent has changed

	for k, v := range shared.DefaultRegistry.Dependents {
		addrs := shared.DefaultRegistry.VarAddress[k]
		requirednodes := []*Node{}

		for _, uid := range addrs {
			requirednodes = append(requirednodes, s.AllNodes[uid])
		}

		for _, reqnode := range requirednodes {
			for _, b := range v {
				//check if uid is blacklisted

				if reqnode != nil {
					// fmt.Println(b)
					childnode := s.AllNodes[b]

					// fmt.Printf("childnode is %v\n", childnode)

					childnode.Root = append(childnode.Root, reqnode)
					reqnode.Children = append(reqnode.Children, childnode)
				}

			}
		}

		// reqnode := s.AllNodes[id]
		// fmt.Printf("%v requires %v\n", v, k)

	}

	// temporary patch

	// Identify Root Nodes
	roots := []*Node{}
	for _, v := range s.AllNodes {
		if len(v.Root) == 0 {
			roots = append(roots, v)
		}
	}

	s.ExecPyramid = append(s.ExecPyramid, roots)

	s.createExecutionPyramid(roots)

	// fix inconsistencies Ex: A-B-C if b dropped then A-C (All complex cases)

	// new pyramid
	npy := [][]*Node{}

	for _, v := range s.ExecPyramid {
		arr := []*Node{}
		for _, n := range v {
			if !s.BlackListed[n.UID] {
				arr = append(arr, n)
			}
		}

		if len(arr) > 0 {
			npy = append(npy, arr)
		}
	}

	//use npy as execpyramid
	s.ExecPyramid = npy

	count := 0

	ioutils.Cout.Printf("[*] Execution Pyramid by levels (top->bottom)\n")

	for _, v := range s.ExecPyramid {
		ioutils.Cout.Printf("Level %v : ", count)
		for _, b := range v {
			ioutils.Cout.Printf("%v", b.Comment)
		}
		count += 1
		ioutils.Cout.Printf("")
	}

	ioutils.Cout.DrawLine(30)

}

func (s *Scheduler) createExecutionPyramid(CurrNodes []*Node) {
	cnodes := CurrNodes

	for {
		// if there are no nodes just return
		if len(cnodes) < 1 {
			break
		}

		// If children have more than 1 root
		// Using SOmething Similar to DFS and give preference to farthest one
		// also autobalance if any node is blacklisted
		for _, v := range cnodes {
			for _, x2 := range v.Children {
				if len(x2.Root) > 1 {

					// remove parent root if it is blacklisted
					// add its children to level zero

					newroutes := []*Node{}

					for _, v := range x2.Root {
						//skip if it is blacklisted
						if !s.BlackListed[v.UID] {
							newroutes = append(newroutes, v)
						}
					}

					//after filtering blacklisted nodes
					if len(newroutes) == 0 {
						//change status to roots and add to top
						x2.Root = []*Node{}
						s.ExecPyramid[0] = append(s.ExecPyramid[0], x2)
					} else {
						//if not use same remaining as parents
						x2.Root = newroutes
						if len(newroutes) > 1 {
							rnode, _ := FarthestNodethru(x2, 0)
							x2.Root = []*Node{rnode}
						}
					}

				}
			}
		}

		tmpchildren := []*Node{}
		for _, v := range cnodes {
			for _, z := range v.Children {
				if len(z.Root) == 1 {
					if z.Root[0] == v {
						tmpchildren = append(tmpchildren, z)
					}
				}

			}
		}
		// levelarr = append(levelarr, tmp)
		if len(tmpchildren) > 0 {
			s.ExecPyramid = append(s.ExecPyramid, tmpchildren)
		}

		cnodes = tmpchildren

		if len(cnodes) < 1 {
			break
		}

	}

	// return levelarr
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		AllNodes:    map[string]*Node{},
		ExecPyramid: [][]*Node{},
		BlackListed: map[string]bool{},
	}
}
