package visualizer

import (
	"fmt"
	"math"
	"sync"

	"github.com/bhusalashish/consistent-hashing-bounded-loads.git/pkg/hash"
	"github.com/bhusalashish/consistent-hashing-bounded-loads.git/pkg/router"
	"github.com/bhusalashish/consistent-hashing-bounded-loads.git/pkg/router/chbl"
	"github.com/bhusalashish/consistent-hashing-bounded-loads.git/pkg/routercore"
)

// State represents the current state of the visualization.
type State struct {
	Nodes         []string            `json:"nodes"`
	Keys          []string            `json:"keys"`
	Positions     map[string]float64  `json:"positions"`   // key/node → position (0..1 around circle)
	Assignments   map[string]string   `json:"assignments"` // key → node
	NodeAngles    map[string]float64  `json:"nodeAngles"`  // node → angle in radians
	Algorithm     string              `json:"algorithm"`
	Stats         *Statistics         `json:"stats,omitempty"` // Statistics for the last operation
	CHBLConfig    *CHBLConfig         `json:"chblConfig,omitempty"` // CH-BL specific config
}

// CHBLConfig contains CH-BL algorithm configuration.
type CHBLConfig struct {
	LoadFactor    float64 `json:"loadFactor"`
	ExpectedKeys  int     `json:"expectedKeys"`
	CapacityPerNode int   `json:"capacityPerNode"`
}

// Statistics tracks key movement and distribution changes.
type Statistics struct {
	Operation        string            `json:"operation"`        // "add-node", "remove-node", "regenerate-keys", "set-algorithm"
	TotalKeys        int               `json:"totalKeys"`
	KeysMoved        int               `json:"keysMoved"`
	KeysMovedPercent float64           `json:"keysMovedPercent"`
	MovementByNode   map[string]int    `json:"movementByNode"`   // node → count of keys moved to/from
	Distribution     map[string]int    `json:"distribution"`      // node → current key count
	PreviousDist     map[string]int    `json:"previousDist"`     // node → previous key count
	KeyMovements     []KeyMovement     `json:"keyMovements"`     // Detailed movements (limited to first 20)
	// Capacity information for CH-BL
	CapacityInfo     *CapacityInfo     `json:"capacityInfo,omitempty"` // CH-BL capacity details
}

// CapacityInfo provides information about node capacity status for CH-BL.
type CapacityInfo struct {
	NodesAtCapacity  []string          `json:"nodesAtCapacity"`  // Nodes that are at capacity
	UnassignedKeys   int               `json:"unassignedKeys"`  // Number of keys that couldn't be assigned
	CapacityPerNode  map[string]int    `json:"capacityPerNode"` // node → capacity limit
	CurrentLoad       map[string]int    `json:"currentLoad"`     // node → current load
	LoadPercentage    map[string]float64 `json:"loadPercentage"`  // node → load percentage (0-100)
}

// KeyMovement represents a single key moving from one node to another.
type KeyMovement struct {
	KeyID    string `json:"keyId"`
	FromNode string `json:"fromNode"`
	ToNode   string `json:"toNode"`
}

// Manager manages the visualizer state and router.
type Manager struct {
	mu            sync.RWMutex
	mapper        routercore.Mapper
	nodes         []string
	keys          []string
	algo          routercore.Algo
	opts          routercore.Options
	keyGen        int // counter for generating unique keys
	prevAssignments map[string]string // previous assignments for computing stats
}

// NewManager creates a new visualizer manager.
func NewManager() *Manager {
	m := &Manager{
		nodes:          make([]string, 0),
		keys:           make([]string, 0),
		algo:           routercore.AlgoRing,
		keyGen:         0,
		prevAssignments: make(map[string]string),
		opts: routercore.Options{
			TableSize:     65537,
			LoadFactor:    1.25,
			Vnodes:        100,
			WalkThreshold: 8,
			HashSeed:      42,
			ExpectedKeys:  1000,
		},
	}
	// Initialize with 3 nodes
	for i := 0; i < 3; i++ {
		m.nodes = append(m.nodes, m.generateNodeID())
	}
	// Build initial mapper
	m.rebuild()
	// Initialize with 50 keys
	m.RegenerateKeys(50)
	// Initialize previous assignments
	for _, key := range m.keys {
		m.prevAssignments[key] = m.mapper.Pick([]byte(key))
	}
	return m
}

// GetState returns the current visualization state.
func (m *Manager) GetState() (*State, error) {
	return m.getStateWithStats("", nil)
}

// getStateWithStats returns the state with computed statistics.
func (m *Manager) getStateWithStats(operation string, stats *Statistics) (*State, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.mapper == nil {
		m.rebuild()
	}

	state := &State{
		Nodes:       make([]string, len(m.nodes)),
		Keys:        make([]string, len(m.keys)),
		Positions:   make(map[string]float64),
		Assignments: make(map[string]string),
		NodeAngles:  make(map[string]float64),
		Algorithm:   string(m.algo),
		Stats:       stats,
	}

	// Add CH-BL config if using CH-BL
	if m.algo == routercore.AlgoCHBL && len(m.nodes) > 0 {
		avg := float64(m.opts.ExpectedKeys) / float64(len(m.nodes))
		capacityPerNode := int(math.Ceil(m.opts.LoadFactor * avg))
		state.CHBLConfig = &CHBLConfig{
			LoadFactor:     m.opts.LoadFactor,
			ExpectedKeys:   m.opts.ExpectedKeys,
			CapacityPerNode: capacityPerNode,
		}
	}

	copy(state.Nodes, m.nodes)
	copy(state.Keys, m.keys)

	// Compute node positions (evenly spaced around circle)
	for i, node := range m.nodes {
		angle := 2 * math.Pi * float64(i) / float64(len(m.nodes))
		state.NodeAngles[node] = angle
		state.Positions[node] = angle / (2 * math.Pi) // normalize to 0..1
	}

	// Compute key positions and assignments
	if m.mapper != nil && len(m.nodes) > 0 {
		for _, key := range m.keys {
			// Use recover to handle potential panics from Pick
			func() {
				defer func() {
					if r := recover(); r != nil {
						// If Pick panics, skip this key
					}
				}()
				node := m.mapper.Pick([]byte(key))
				state.Assignments[key] = node

				// Compute key position based on hash
				h := hash.XXH64([]byte(key), m.opts.HashSeed)
				// Normalize hash to 0..1
				position := float64(h) / float64(math.MaxUint64)
				state.Positions[key] = position
			}()
		}
	}

	return state, nil
}

// computeStatistics computes statistics by comparing current and previous assignments.
func (m *Manager) computeStatistics(operation string) *Statistics {
	stats := &Statistics{
		Operation:      operation,
		TotalKeys:      len(m.keys),
		MovementByNode: make(map[string]int),
		Distribution:   make(map[string]int),
		PreviousDist:   make(map[string]int),
		KeyMovements:   make([]KeyMovement, 0),
	}

	// Compute previous distribution - include all nodes that had keys
	allPreviousNodes := make(map[string]bool)
	for _, prevNode := range m.prevAssignments {
		allPreviousNodes[prevNode] = true
	}
	for node := range allPreviousNodes {
		stats.PreviousDist[node] = 0
	}
	for _, node := range m.nodes {
		stats.PreviousDist[node] = 0
	}
	for _, key := range m.keys {
		if prevNode, exists := m.prevAssignments[key]; exists {
			stats.PreviousDist[prevNode]++
		}
	}

	// For CH-BL, we need to handle stateful load tracking differently
	// CH-BL's Pick increments load, so we need to assign all keys first to build up
	// the correct load state, then compute statistics from those assignments
	if m.algo == routercore.AlgoCHBL && m.mapper != nil {
		// Rebuild to reset load state
		m.rebuild()
		
		// Ensure ExpectedKeys is sufficient before assigning
		// Capacity per node = ceil(LoadFactor * ExpectedKeys / numNodes)
		// Total capacity = numNodes * ceil(LoadFactor * ExpectedKeys / numNodes)
		// We need total capacity >= totalKeys
		// Simplifying: ExpectedKeys >= totalKeys / LoadFactor (approximately)
		// Add 20% buffer to be safe
		if len(m.nodes) > 0 && len(m.keys) > 0 {
			minExpectedKeys := int(float64(len(m.keys)) / m.opts.LoadFactor * 1.2)
			if minExpectedKeys < len(m.keys) {
				minExpectedKeys = len(m.keys) // At least equal to key count
			}
			if m.opts.ExpectedKeys < minExpectedKeys {
				m.opts.ExpectedKeys = minExpectedKeys
				m.rebuild()
			}
		}
		
		// Pre-assign all keys to build up load state correctly
		// This ensures capacity is respected as keys are assigned
		currentAssignments := make(map[string]string)
		unassignedKeys := 0
		for _, key := range m.keys {
			if m.mapper != nil {
				assignedNode := m.mapper.Pick([]byte(key))
				if assignedNode == "" {
					// Key couldn't be assigned - all nodes at capacity
					unassignedKeys++
				} else {
					currentAssignments[key] = assignedNode
				}
			}
		}
		
		// Get capacity information from CH-BL mapper
		if chblMapper, ok := m.mapper.(chbl.CHBLMapper); ok {
			capacityStatus := chblMapper.GetCapacityStatus()
			stats.CapacityInfo = &CapacityInfo{
				NodesAtCapacity: capacityStatus.NodesAtCapacity,
				UnassignedKeys:  unassignedKeys,
				CapacityPerNode: capacityStatus.CapacityPerNode,
				CurrentLoad:     capacityStatus.CurrentLoad,
				LoadPercentage:  capacityStatus.LoadPercentage,
			}
		} else if unassignedKeys > 0 {
			// If we have unassigned keys but can't get capacity status, create basic info
			stats.CapacityInfo = &CapacityInfo{
				UnassignedKeys: unassignedKeys,
				NodesAtCapacity: make([]string, 0),
				CapacityPerNode: make(map[string]int),
				CurrentLoad:     make(map[string]int),
				LoadPercentage: make(map[string]float64),
			}
		}
		
		// Now compute statistics from the assignments
		for _, key := range m.keys {
			currentNode, exists := currentAssignments[key]
			if !exists || currentNode == "" {
				continue
			}
			
			stats.Distribution[currentNode]++
			prevNode, hadPrevious := m.prevAssignments[key]
			if hadPrevious && prevNode != currentNode {
				// Key moved
				stats.KeysMoved++
				stats.MovementByNode[currentNode]++
				if len(stats.KeyMovements) < 20 {
					stats.KeyMovements = append(stats.KeyMovements, KeyMovement{
						KeyID:    key,
						FromNode: prevNode,
						ToNode:   currentNode,
					})
				}
			} else if !hadPrevious {
				// New key
				stats.MovementByNode[currentNode]++
			}
		}
		
		// Update previous assignments using the assignments we computed (don't call Pick again!)
		m.prevAssignments = currentAssignments
	} else {
		// For non-CH-BL algorithms, Pick is stateless, so we can call it directly
		for _, key := range m.keys {
			if m.mapper != nil {
				currentNode := m.mapper.Pick([]byte(key))
				stats.Distribution[currentNode]++

				prevNode, hadPrevious := m.prevAssignments[key]
				if hadPrevious && prevNode != currentNode {
					// Key moved
					stats.KeysMoved++
					stats.MovementByNode[currentNode]++
					if len(stats.KeyMovements) < 20 {
						stats.KeyMovements = append(stats.KeyMovements, KeyMovement{
							KeyID:    key,
							FromNode: prevNode,
							ToNode:   currentNode,
						})
					}
				} else if !hadPrevious {
					// New key
					stats.MovementByNode[currentNode]++
				}
			}
		}
	}

	if stats.TotalKeys > 0 {
		stats.KeysMovedPercent = float64(stats.KeysMoved) / float64(stats.TotalKeys) * 100
	}

	// For non-CH-BL algorithms, update previous assignments
	// (CH-BL assignments are already updated in computeStatistics above)
	if m.algo != routercore.AlgoCHBL {
		m.prevAssignments = make(map[string]string)
		for _, key := range m.keys {
			if m.mapper != nil {
				m.prevAssignments[key] = m.mapper.Pick([]byte(key))
			}
		}
	}

	return stats
}

// SetAlgorithm changes the routing algorithm.
func (m *Manager) SetAlgorithm(algo string) (*Statistics, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	switch algo {
	case "ring":
		m.algo = routercore.AlgoRing
	case "jump":
		m.algo = routercore.AlgoJump
	case "maglev":
		m.algo = routercore.AlgoMaglev
	case "chbl":
		m.algo = routercore.AlgoCHBL
	default:
		return nil, routercore.ErrUnknownAlgo
	}

	m.rebuild()
	stats := m.computeStatistics("set-algorithm")
	return stats, nil
}

// SetCHBLConfig updates CH-BL algorithm parameters.
func (m *Manager) SetCHBLConfig(loadFactor float64, expectedKeys int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.algo != routercore.AlgoCHBL {
		return fmt.Errorf("can only set CH-BL config when using CH-BL algorithm")
	}

	if loadFactor <= 0 {
		return fmt.Errorf("loadFactor must be > 0")
	}
	if expectedKeys <= 0 {
		return fmt.Errorf("expectedKeys must be > 0")
	}

	m.opts.LoadFactor = loadFactor
	m.opts.ExpectedKeys = expectedKeys
	m.rebuild()
	return nil
}

// AlgorithmComparison represents the result of running an operation on one algorithm.
type AlgorithmComparison struct {
	Algorithm   string            `json:"algorithm"`
	State       *State            `json:"state"`
	Stats       *Statistics       `json:"stats"`
}

// CompareOperation runs the same operation on all algorithms and returns comparison data.
func (m *Manager) CompareOperation(operation string, nodeID string) ([]AlgorithmComparison, error) {
	m.mu.Lock()
	
	// Save current state (we don't modify the original manager)
	currentNodes := make([]string, len(m.nodes))
	copy(currentNodes, m.nodes)
	currentKeys := make([]string, len(m.keys))
	copy(currentKeys, m.keys)
	currentPrevAssignments := make(map[string]string)
	for k, v := range m.prevAssignments {
		currentPrevAssignments[k] = v
	}
	currentKeyGen := m.keyGen
	currentOpts := m.opts

	m.mu.Unlock()

	algorithms := []routercore.Algo{
		routercore.AlgoRing,
		routercore.AlgoJump,
		routercore.AlgoMaglev,
		routercore.AlgoCHBL,
	}

	results := make([]AlgorithmComparison, 0, len(algorithms))

	for _, algo := range algorithms {
		// Create a temporary manager for this algorithm
		tempManager := &Manager{
			nodes:          make([]string, len(currentNodes)),
			keys:           make([]string, len(currentKeys)),
			algo:           algo,
			keyGen:         currentKeyGen,
			prevAssignments: make(map[string]string),
			opts:           currentOpts,
		}
		copy(tempManager.nodes, currentNodes)
		copy(tempManager.keys, currentKeys)
		for k, v := range currentPrevAssignments {
			tempManager.prevAssignments[k] = v
		}

		tempManager.rebuild()

		// For CH-BL, ensure ExpectedKeys is sufficient (with buffer for load factor)
		if tempManager.algo == routercore.AlgoCHBL {
			requiredExpectedKeys := int(float64(len(currentKeys)) * 1.5)
			if tempManager.opts.ExpectedKeys < requiredExpectedKeys {
				tempManager.opts.ExpectedKeys = requiredExpectedKeys
				tempManager.rebuild()
			}
		}

		var stats *Statistics

		// Perform operation
		switch operation {
		case "add-node":
			newNodeID := tempManager.generateNodeID()
			tempManager.nodes = append(tempManager.nodes, newNodeID)
			tempManager.rebuild()
			// Update ExpectedKeys for CH-BL if needed (with buffer)
			if tempManager.algo == routercore.AlgoCHBL {
				requiredExpectedKeys := int(float64(len(tempManager.keys)) * 1.5)
				if tempManager.opts.ExpectedKeys < requiredExpectedKeys {
					tempManager.opts.ExpectedKeys = requiredExpectedKeys
					tempManager.rebuild()
				}
			}
			stats = tempManager.computeStatistics("add-node")
		case "remove-node":
			if nodeID == "" {
				continue
			}
			var newNodes []string
			for _, n := range tempManager.nodes {
				if n != nodeID {
					newNodes = append(newNodes, n)
				}
			}
			if len(newNodes) == 0 {
				continue
			}
			tempManager.nodes = newNodes
			if len(tempManager.nodes) == 0 {
				continue
			}
			tempManager.rebuild()
			// Update ExpectedKeys for CH-BL if needed (with buffer)
			if tempManager.algo == routercore.AlgoCHBL {
				requiredExpectedKeys := int(float64(len(tempManager.keys)) * 1.5)
				if tempManager.opts.ExpectedKeys < requiredExpectedKeys {
					tempManager.opts.ExpectedKeys = requiredExpectedKeys
					tempManager.rebuild()
				}
			}
			stats = tempManager.computeStatistics("remove-node")
		case "regenerate-keys":
			tempManager.keys = make([]string, len(currentKeys))
			for i := 0; i < len(currentKeys); i++ {
				tempManager.keys[i] = tempManager.generateKey()
			}
			// Update ExpectedKeys for CH-BL if needed (with buffer)
			if tempManager.algo == routercore.AlgoCHBL {
				requiredExpectedKeys := int(float64(len(tempManager.keys)) * 1.5)
				if tempManager.opts.ExpectedKeys < requiredExpectedKeys {
					tempManager.opts.ExpectedKeys = requiredExpectedKeys
					tempManager.rebuild()
				}
			}
			stats = tempManager.computeStatistics("regenerate-keys")
		default:
			continue
		}

		// Get state for this algorithm
		state, err := tempManager.getStateWithStats(operation, stats)
		if err != nil {
			continue
		}

		results = append(results, AlgorithmComparison{
			Algorithm: string(algo),
			State:     state,
			Stats:     stats,
		})
	}

	return results, nil
}

// AddNode adds a new node to the ring.
func (m *Manager) AddNode() (string, *Statistics) {
	m.mu.Lock()
	defer m.mu.Unlock()

	nodeID := m.generateNodeID()
	m.nodes = append(m.nodes, nodeID)
	m.rebuild()
	
	stats := m.computeStatistics("add-node")
	return nodeID, stats
}

// RemoveNode removes a node from the ring.
func (m *Manager) RemoveNode(nodeID string) (*Statistics, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	found := false
	var newNodes []string
	for _, n := range m.nodes {
		if n != nodeID {
			newNodes = append(newNodes, n)
		} else {
			found = true
		}
	}

	if !found {
		return nil, nil // node doesn't exist, no-op
	}

	m.nodes = newNodes
	if len(m.nodes) == 0 {
		m.mapper = nil
		return nil, nil
	}

	m.rebuild()
	stats := m.computeStatistics("remove-node")
	return stats, nil
}

// RegenerateKeys generates a new set of keys.
func (m *Manager) RegenerateKeys(count int) *Statistics {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.keys = make([]string, count)
	for i := 0; i < count; i++ {
		m.keys[i] = m.generateKey()
	}
	
	// For CH-BL, ensure ExpectedKeys is at least as large as the key count
	// Add some buffer (1.5x) to account for load factor
	if m.algo == routercore.AlgoCHBL {
		requiredExpectedKeys := int(float64(count) * 1.5)
		if m.opts.ExpectedKeys < requiredExpectedKeys {
			m.opts.ExpectedKeys = requiredExpectedKeys
			m.rebuild()
		}
	}
	
	stats := m.computeStatistics("regenerate-keys")
	return stats
}

// SetKeyCount adjusts the number of keys.
func (m *Manager) SetKeyCount(count int) *Statistics {
	m.mu.Lock()
	defer m.mu.Unlock()

	if count < len(m.keys) {
		m.keys = m.keys[:count]
	} else {
		for len(m.keys) < count {
			m.keys = append(m.keys, m.generateKey())
		}
	}
	
	// For CH-BL, ensure ExpectedKeys is at least as large as the key count
	// Add some buffer (1.5x) to account for load factor
	if m.algo == routercore.AlgoCHBL {
		requiredExpectedKeys := int(float64(count) * 1.5)
		if m.opts.ExpectedKeys < requiredExpectedKeys {
			m.opts.ExpectedKeys = requiredExpectedKeys
			m.rebuild()
		}
	}
	
	stats := m.computeStatistics("set-key-count")
	return stats
}

// rebuild rebuilds the mapper with current nodes.
func (m *Manager) rebuild() {
	if len(m.nodes) == 0 {
		m.mapper = nil
		return
	}

	var err error
	m.mapper, err = router.New(m.algo, m.opts, m.nodes)
	if err != nil {
		// Fallback to ring if algo fails
		m.mapper, _ = router.New(routercore.AlgoRing, m.opts, m.nodes)
	}
}

func (m *Manager) generateNodeID() string {
	id := m.keyGen
	m.keyGen++
	return fmt.Sprintf("node-%d", id)
}

func (m *Manager) generateKey() string {
	id := m.keyGen
	m.keyGen++
	return fmt.Sprintf("key-%d", id)
}

