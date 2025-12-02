package visualizer

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// API provides HTTP handlers for the visualizer.
type API struct {
	manager *Manager
}

// NewAPI creates a new API instance.
func NewAPI() *API {
	return &API{
		manager: NewManager(),
	}
}

// StateResponse is the response for GET /state.
type StateResponse struct {
	State *State `json:"state"`
}

// AddNodeRequest is the request for POST /add-node.
type AddNodeRequest struct {
	// Empty for now, could add node ID later
}

// AddNodeResponse is the response for POST /add-node.
type AddNodeResponse struct {
	NodeID string `json:"nodeId"`
	State  *State `json:"state"`
}

// RemoveNodeRequest is the request for POST /remove-node.
type RemoveNodeRequest struct {
	NodeID string `json:"nodeId"`
}

// RemoveNodeResponse is the response for POST /remove-node.
type RemoveNodeResponse struct {
	State *State `json:"state"`
}

// RegenerateKeysRequest is the request for POST /regenerate-keys.
type RegenerateKeysRequest struct {
	Count int `json:"count"`
}

// RegenerateKeysResponse is the response for POST /regenerate-keys.
type RegenerateKeysResponse struct {
	State *State `json:"state"`
}

// SetAlgorithmRequest is the request for POST /set-algorithm.
type SetAlgorithmRequest struct {
	Algorithm string `json:"algorithm"`
}

// SetAlgorithmResponse is the response for POST /set-algorithm.
type SetAlgorithmResponse struct {
	State *State `json:"state"`
}

// SetCHBLConfigRequest is the request for POST /set-chbl-config.
type SetCHBLConfigRequest struct {
	LoadFactor   float64 `json:"loadFactor"`
	ExpectedKeys int     `json:"expectedKeys"`
}

// SetCHBLConfigResponse is the response for POST /set-chbl-config.
type SetCHBLConfigResponse struct {
	State *State `json:"state"`
}

// HandleState returns the current state.
func (a *API) HandleState(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		handleCORS(w)
		return
	}
	if r.Method != http.MethodGet {
		handleCORS(w)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	state, err := a.manager.GetState()
	if err != nil {
		handleCORS(w)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, StateResponse{State: state})
}

// HandleCompareOperation runs an operation on all algorithms and returns comparison data.
func (a *API) HandleCompareOperation(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		handleCORS(w)
		return
	}
	if r.Method != http.MethodPost {
		handleCORS(w)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Operation string `json:"operation"` // "add-node", "remove-node", "regenerate-keys"
		NodeID    string `json:"nodeId,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleCORS(w)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	comparison, err := a.manager.CompareOperation(req.Operation, req.NodeID)
	if err != nil {
		handleCORS(w)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, map[string]interface{}{"comparison": comparison})
}

// HandleAddNode adds a new node.
func (a *API) HandleAddNode(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		handleCORS(w)
		return
	}
	if r.Method != http.MethodPost {
		handleCORS(w)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AddNodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Empty body is OK
	}

	nodeID, stats := a.manager.AddNode()
	state, err := a.manager.getStateWithStats("add-node", stats)
	if err != nil {
		handleCORS(w)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, AddNodeResponse{NodeID: nodeID, State: state})
}

// HandleRemoveNode removes a node.
func (a *API) HandleRemoveNode(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		handleCORS(w)
		return
	}
	if r.Method != http.MethodPost {
		handleCORS(w)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RemoveNodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleCORS(w)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.NodeID == "" {
		handleCORS(w)
		http.Error(w, "nodeId is required", http.StatusBadRequest)
		return
	}

	stats, err := a.manager.RemoveNode(req.NodeID)
	if err != nil {
		handleCORS(w)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	state, err := a.manager.getStateWithStats("remove-node", stats)
	if err != nil {
		handleCORS(w)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, RemoveNodeResponse{State: state})
}

// HandleRegenerateKeys regenerates keys.
func (a *API) HandleRegenerateKeys(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		handleCORS(w)
		return
	}
	if r.Method != http.MethodPost {
		handleCORS(w)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegenerateKeysRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		req.Count = 50 // default
	}

	if req.Count <= 0 {
		req.Count = 50
	}
	if req.Count > 1000 {
		req.Count = 1000 // cap at 1000 for performance
	}

	stats := a.manager.RegenerateKeys(req.Count)
	state, err := a.manager.getStateWithStats("regenerate-keys", stats)
	if err != nil {
		handleCORS(w)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, RegenerateKeysResponse{State: state})
}

// HandleSetAlgorithm changes the algorithm.
func (a *API) HandleSetAlgorithm(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		handleCORS(w)
		return
	}
	if r.Method != http.MethodPost {
		handleCORS(w)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SetAlgorithmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleCORS(w)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	stats, err := a.manager.SetAlgorithm(req.Algorithm)
	if err != nil {
		handleCORS(w)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	state, err := a.manager.getStateWithStats("set-algorithm", stats)
	if err != nil {
		handleCORS(w)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, SetAlgorithmResponse{State: state})
}

// HandleSetKeyCount adjusts the number of keys.
func (a *API) HandleSetKeyCount(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		handleCORS(w)
		return
	}
	if r.Method != http.MethodPost {
		handleCORS(w)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	countStr := r.URL.Query().Get("count")
	count, err := strconv.Atoi(countStr)
	if err != nil || count < 0 {
		handleCORS(w)
		http.Error(w, "Invalid count parameter", http.StatusBadRequest)
		return
	}

	if count > 1000 {
		count = 1000
	}

	stats := a.manager.SetKeyCount(count)
	state, err := a.manager.getStateWithStats("set-key-count", stats)
	if err != nil {
		handleCORS(w)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, StateResponse{State: state})
}

// HandleSetCHBLConfig updates CH-BL algorithm parameters.
func (a *API) HandleSetCHBLConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		handleCORS(w)
		return
	}
	if r.Method != http.MethodPost {
		handleCORS(w)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SetCHBLConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleCORS(w)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := a.manager.SetCHBLConfig(req.LoadFactor, req.ExpectedKeys); err != nil {
		handleCORS(w)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	state, err := a.manager.GetState()
	if err != nil {
		handleCORS(w)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, SetCHBLConfigResponse{State: state})
}

// handleCORS sets CORS headers for preflight and regular requests.
func handleCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Max-Age", "3600")
}

// respondJSON sends a JSON response.
func respondJSON(w http.ResponseWriter, data interface{}) {
	handleCORS(w)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		handleCORS(w)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

