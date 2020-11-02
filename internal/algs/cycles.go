package algs

type BaseNode interface {
	GetAdjacent() []int
}

type NodeState int

const (
	notVisited        NodeState = 0
	partiallyVisited            = 1
	completelyVisited           = 2
)

type Node struct {
	ind int

	adjacent []*Node
	parent   *Node
	state    NodeState
}

func dfs(node, prev *Node, cycles *[][]*Node) {
	if node.state == completelyVisited {
		return
	}

	if node.state == partiallyVisited {
		cycle := []*Node{}

		cur := prev
		cycle = append(cycle, cur)

		for cur != node {
			cur = cur.parent
			cycle = append(cycle, cur)
		}

		*cycles = append(*cycles, cycle)

		return
	}

	node.parent = prev
	node.state = partiallyVisited

	for _, v := range node.adjacent {
		if node.parent == v {
			continue
		}

		dfs(v, node, cycles)
	}

	node.state = completelyVisited
}

func cleanNodes(nodes []*Node) {
	for _, node := range nodes {
		node.parent = nil
		node.state = notVisited
	}
}

func DetectCycles(baseNodes []BaseNode) (res [][]BaseNode) {
	// TODO: rewrite to return all induced cycles

	nodes := make([]*Node, len(baseNodes))
	for i := range nodes {
		nodes[i] = &Node{ind: i}
	}

	for i, baseNode := range baseNodes {
		baseAdjacent := baseNode.GetAdjacent()
		// Pre-allocate memory for adjacent nodes
		nodes[i].adjacent = make([]*Node, len(baseAdjacent))

		for j, adjInd := range baseAdjacent {
			nodes[i].adjacent[j] = nodes[adjInd]
		}
	}

	var cycles [][]*Node

	dfs(nodes[0], nil, &cycles)
	/*
		for _, node := range nodes {
			dfs(node, nil, &cycles)
			cleanNodes(nodes)
		}
	*/

	for _, cycle := range cycles {
		baseCycle := []BaseNode{}
		for _, node := range cycle {
			baseCycle = append(baseCycle, baseNodes[node.ind])
		}

		res = append(res, baseCycle)
	}

	return
}
