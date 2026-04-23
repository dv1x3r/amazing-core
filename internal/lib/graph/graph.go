package graph

type AdjacencyMap map[int]map[int]struct{}

func (nodes AdjacencyMap) Add(nodeID int, childID int) {
	if nodes[nodeID] == nil {
		nodes[nodeID] = map[int]struct{}{}
	}
	nodes[nodeID][childID] = struct{}{}
}

// Breadth-First Search for circular dependencies
func (nodes AdjacencyMap) HasCycleFrom(startingNode int) bool {
	visited := map[int]struct{}{}
	queue := []int{}

	// fmt.Printf("startingNode: %d\n", startingNode)
	// for nodeID, edges := range nodes {
	// 	fmt.Printf("%d: %v\n", nodeID, edges)
	// }

	// seed the queue with direct neighbors
	for neighborID := range nodes[startingNode] {
		queue = append(queue, neighborID)
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current == startingNode {
			return true
		}

		if _, seen := visited[current]; seen {
			continue
		}
		visited[current] = struct{}{}

		for neighborID := range nodes[current] {
			queue = append(queue, neighborID)
		}
	}

	return false
}
