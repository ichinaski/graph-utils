package search

import (
    "github.com/ichinaski/graph-utils/container"
    "github.com/ichinaski/graph-utils/graph"
)

const (
    type_bfs = iota
    type_dfs
    type_dijkstra
    type_astar
)

// State is a graph position for which its ancestors have been evaluated.
// Contains the nodeId, its parent int, and the total cost of traversing 
// the graph to reach this position
type State struct {
    nodeId, parentId int
    cost int
}

// OpenSet defines the functions that any container used to keep track of
// non-expanded nodes must implement
type OpenSet interface {
    Push(item interface{}, priority int)
    Pop() interface{}
    Len() int
}

// Wrap Queue to match OpenSet interface signature
type OpenQueue struct { *container.Queue }
func (q *OpenQueue ) Push(item interface{}, priority int) { q.Queue.Push(item) }

// Wrap Stack to match OpenSet interface signature
type OpenStack struct { *container.Stack }
func (s *OpenStack) Push(item interface{}, priority int) { s.Stack.Push(item) }

// Heuristic calculates the estimated cost between two nodes
type Heuristic func(node, goal int) int

func NullHeuristic(node, goal int) int { return 0 }

func BreathFirstSearch(graph graph.Graph, start, goal int) []int {
    return search(graph, start, goal, type_bfs, nil)
}

func DepthFirstSearch(graph graph.Graph, start, goal int) []int {
    return search(graph, start, goal, type_dfs, nil)
}

func Dijkstra(graph graph.Graph, start, goal int) []int {
    return search(graph, start, goal, type_dijkstra, nil)
}

func Astar(graph graph.Graph, start, goal int, heuristic Heuristic) []int {
    return search(graph, start, goal, type_astar, heuristic)
}

func search(graph graph.Graph, start, goal int, search_type uint, heuristic Heuristic) []int {
    if heuristic == nil { heuristic = NullHeuristic }

    // Initialize the open set, according to the type of search passed in
    var openSet OpenSet
    switch search_type {
        case type_bfs:
            openSet = &OpenQueue{ container.NewQueue(), }// FIFO approach to expand nodes
        case type_dfs:
            openSet = &OpenStack{ container.NewStack(), }// LIFO approach to expand nodes
        case type_dijkstra, type_astar:
            openSet = &container.PQueue{}// Priority queue if Dijkstra or A*
    }

    closedSet := make(map[int]int)// Visited nodes, with a reference to their direct ancestor

    state := &State { start, 0, 0}
    openSet.Push(state, 0)

    var path []int
    for openSet.Len() > 0 {
        item := openSet.Pop()
        state = item.(*State)

        // Only consider non expanded nodes (not present in closedSet)
        if _, ok := closedSet[state.nodeId]; !ok {
            // Store this node in the closed list, with a reference to its parent
            closedSet[state.nodeId] = state.parentId

            if state.nodeId == goal {
                path = calculatePath(start, goal, closedSet)
                break
            }

            // Add the nodes not present in the closedSet into the openSet
            edges := graph.GetEdges(state.nodeId)

            // for depth-first search, we have to alter the order in which we add the successors to the stack,
            // to ensure items are expanded as expected (in the order they have been passed in)
            if (search_type == type_dfs) {
                for i, j := 0, len(edges)-1; i < j; i, j = i+1, j-1 {
                    edges[i], edges[j] = edges[j], edges[i]
                }
            }

            for _, edge := range edges {
                if _, ok := closedSet[edge.NodeId]; !ok {
                    nextState := &State { edge.NodeId, state.nodeId, state.cost + edge.Cost}
                    openSet.Push(nextState, nextState.cost + heuristic(edge.NodeId, goal))
                }
            }
        }
    }

    return path
}

func calculatePath(start, goal int, closedSet map[int]int) []int {
    path := make([]int, 0, len(closedSet))
    // fetch all the nodes in a descendant way, from goal to start
    nodeId := goal
    for nodeId != 0 {
        path = append(path, nodeId)
        nodeId = closedSet[nodeId]
    }

    // Reverse the slice
    for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
        path[i], path[j] = path[j], path[i]
    }

    return path
}