package scheduler

// Distance of Node From Root
func FarthestNodethru(x *Node, depth int) (*Node, int) {
	if len(x.Root) == 0 {
		return x, depth
	} else {
		if len(x.Root) == 1 {
			_, val := FarthestNodethru(x.Root[0], depth+1)
			return x, val
		}
		rvals := map[*Node]int{}
		for _, r := range x.Root {
			_, rvals[r] = FarthestNodethru(r, depth+1)
		}
		return MaxofArr(rvals)

	}
}

func MaxofArr(z map[*Node]int) (*Node, int) {
	max := -1
	var reqnode *Node
	for k, v := range z {
		// fmt.Printf("%v with %v\n", k.Comment, v)
		if v > max {
			max = v
			reqnode = k
		}
	}

	return reqnode, max
}
