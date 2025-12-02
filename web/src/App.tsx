import React, { useState, useEffect, useCallback } from 'react';
import Ring from './components/Ring';
import ControlPanel from './components/ControlPanel';
import StatsPanel from './components/StatsPanel';
import EducationPanel from './components/EducationPanel';
import WelcomeGuide from './components/WelcomeGuide';
import ComparisonPanel from './components/ComparisonPanel';
import DynamicComparisonPanel from './components/DynamicComparisonPanel';
import { VisualizerState, Algorithm } from './types';
import * as api from './api';
import './App.css';

const App: React.FC = () => {
  const [state, setState] = useState<VisualizerState | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [hoveredNode, setHoveredNode] = useState<string | null>(null);
  const [hoveredKey, setHoveredKey] = useState<string | null>(null);
  const [showStats, setShowStats] = useState(false);
  const [showEducation, setShowEducation] = useState(false);
  const [showComparison, setShowComparison] = useState(false);
  const [showDynamicComparison, setShowDynamicComparison] = useState(false);
  const [dynamicComparisonOp, setDynamicComparisonOp] = useState<'add-node' | 'remove-node' | 'regenerate-keys' | null>(null);
  const [dynamicComparisonNodeId, setDynamicComparisonNodeId] = useState<string | undefined>(undefined);
  const [showWelcome, setShowWelcome] = useState(() => {
    const hasSeenWelcome = localStorage.getItem('hasSeenWelcome');
    return !hasSeenWelcome;
  });
  const [dimensions, setDimensions] = useState({
    width: window.innerWidth,
    height: window.innerHeight,
  });

  const loadState = useCallback(async () => {
    try {
      setError(null);
      const newState = await api.fetchState();
      setState(newState);
      setLoading(false);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load state');
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadState();
  }, [loadState]);

  // Handle window resize
  useEffect(() => {
    const handleResize = () => {
      setDimensions({
        width: window.innerWidth,
        height: window.innerHeight,
      });
    };

    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, []);

  const handleAddNode = async () => {
    try {
      const newState = await api.addNode();
      setState(newState);
      if (newState.stats) {
        setShowStats(true);
      }
      // Show dynamic comparison
      setDynamicComparisonOp('add-node');
      setShowDynamicComparison(true);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to add node');
    }
  };

  const handleRemoveNode = async (nodeId: string) => {
    try {
      const newState = await api.removeNode(nodeId);
      setState(newState);
      if (newState.stats) {
        setShowStats(true);
      }
      // Show dynamic comparison
      setDynamicComparisonOp('remove-node');
      setDynamicComparisonNodeId(nodeId);
      setShowDynamicComparison(true);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to remove node');
    }
  };

  const handleRegenerateKeys = async () => {
    if (!state) return;
    try {
      const newState = await api.regenerateKeys(state.keys.length);
      setState(newState);
      if (newState.stats) {
        setShowStats(true);
      }
      // Show dynamic comparison
      setDynamicComparisonOp('regenerate-keys');
      setShowDynamicComparison(true);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to regenerate keys');
    }
  };

  const handleAlgorithmChange = async (algo: Algorithm) => {
    try {
      const newState = await api.setAlgorithm(algo);
      setState(newState);
      if (newState.stats) {
        setShowStats(true);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to change algorithm');
    }
  };

  const handleKeyCountChange = async (count: number) => {
    try {
      const newState = await api.setKeyCount(count);
      setState(newState);
      if (newState.stats) {
        setShowStats(true);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to change key count');
    }
  };

  const handleCHBLConfigChange = async (loadFactor: number, expectedKeys: number) => {
    try {
      const newState = await api.setCHBLConfig(loadFactor, expectedKeys);
      setState(newState);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to update CH-BL config');
    }
  };

  if (loading) {
    return (
      <div className="app-loading">
        <div className="spinner"></div>
        <p>Loading visualizer...</p>
      </div>
    );
  }

  if (error) {
    const isConnectionError = error.includes('Failed to fetch') || error.includes('NetworkError');
    return (
      <div className="app-error">
        <h2>Error</h2>
        <p>{error}</p>
        {isConnectionError && (
          <div className="error-help">
            <p><strong>Connection Error:</strong> Make sure the Go backend server is running:</p>
            <code>go run ./cmd/visualizer</code>
            <p>The server should be running on <code>http://localhost:8080</code></p>
          </div>
        )}
        <button onClick={loadState}>Retry</button>
      </div>
    );
  }

  if (!state) {
    return <div className="app-error">No state available</div>;
  }

  // Calculate responsive dimensions
  const controlPanelWidth = dimensions.width > 768 ? 320 : dimensions.width;
  const ringWidth = Math.max(400, dimensions.width - (dimensions.width > 768 ? controlPanelWidth : 0));
  const ringHeight = dimensions.height - 80; // Account for header

  return (
    <div className="app">
      <header className="app-header">
        <div>
          <h1>Consistent Hashing Visualizer</h1>
          <p>Interactive visualization of distributed hashing algorithms</p>
        </div>
        <button
          onClick={() => setShowComparison(true)}
          className="header-compare-btn"
          title="Compare algorithms side-by-side"
        >
          📊 Compare Algorithms
        </button>
      </header>

      <div className="app-content">
        <div className="ring-container">
          <Ring
            state={state}
            width={ringWidth}
            height={ringHeight}
            onNodeHover={setHoveredNode}
            onKeyHover={setHoveredKey}
            highlightedKey={hoveredKey}
          />
          {(hoveredNode || hoveredKey) && (
            <div className="tooltip">
              {hoveredNode && (
                <div>
                  <strong>Node:</strong> {hoveredNode}
                  <br />
                  <strong>Keys assigned:</strong>{' '}
                  {state.keys.filter((k) => state.assignments[k] === hoveredNode).length}
                  <br />
                  <br />
                  <div className="tooltip-explanation">
                    <strong>Why keys go here:</strong>
                    <br />
                    {state.algorithm === 'ring' &&
                      'Keys hash to positions on the ring. This node owns all keys between its position and the next node clockwise.'}
                    {state.algorithm === 'jump' &&
                      'The jump function calculates which node this key belongs to based on a mathematical formula, ensuring minimal movement when nodes change.'}
                    {state.algorithm === 'maglev' &&
                      'Keys hash to a lookup table slot. This node was assigned to that slot during table construction using a permutation algorithm.'}
                    {state.algorithm === 'chbl' &&
                      'Keys hash to the ring, then walk clockwise to find this node. If nodes are at capacity, keys continue walking to find available capacity.'}
                  </div>
                </div>
              )}
              {hoveredKey && (
                <div>
                  <strong>Key:</strong> {hoveredKey}
                  <br />
                  <strong>Assigned to:</strong> {state.assignments[hoveredKey]}
                  <br />
                  <strong>Position:</strong> {(state.positions[hoveredKey] * 100).toFixed(2)}% around ring
                  <br />
                  <div className="tooltip-position-explanation">
                    💡 <em>Position {state.positions[hoveredKey] * 100}% means the key landed at {state.positions[hoveredKey] * 360}° around the circular ring (0° = top, 90° = right, 180° = bottom, 270° = left)</em>
                  </div>
                  <br />
                  <div className="tooltip-explanation">
                    <strong>Why this assignment:</strong>
                    <br />
                    {state.algorithm === 'ring' &&
                      `This key hashed to position ${(state.positions[hoveredKey] * 100).toFixed(2)}% on the ring. Moving clockwise from this position, ${state.assignments[hoveredKey]} is the first node encountered (its vnodes are positioned on the ring).`}
                    {state.algorithm === 'jump' &&
                      `The jump consistent hash function directly calculated that this key belongs to ${state.assignments[hoveredKey]} using the formula: bucket = jump(hash(key), numNodes). No ring traversal needed!`}
                    {state.algorithm === 'maglev' &&
                      `This key hashed to a lookup table slot (position ${(state.positions[hoveredKey] * 100).toFixed(2)}% maps to a table index). That slot was assigned to ${state.assignments[hoveredKey]} during the Maglev permutation table construction.`}
                    {state.algorithm === 'chbl' && (
                      <>
                        This key hashed to position {(state.positions[hoveredKey] * 100).toFixed(2)}% on the ring.
                        <br />
                        <br />
                        <strong>Walking process:</strong>
                        <br />
                        1. Key starts at hash position on ring
                        <br />
                        2. Finds first node clockwise (${state.assignments[hoveredKey]})
                        <br />
                        3. Checks if node has capacity (current: {state.keys.filter(k => state.assignments[k] === state.assignments[hoveredKey]).length}/{state.chblConfig?.capacityPerNode || 'N/A'})
                        <br />
                        4. If at capacity, continues walking clockwise to next node
                        <br />
                        5. Assigns to first node with available capacity
                        <br />
                        <br />
                        <em>This ensures no node exceeds {state.chblConfig?.capacityPerNode || 'N/A'} keys (Load Factor: {state.chblConfig?.loadFactor || 'N/A'} × Expected Keys: {state.chblConfig?.expectedKeys || 'N/A'})</em>
                      </>
                    )}
                  </div>
                </div>
              )}
            </div>
          )}
        </div>

        <ControlPanel
          algorithm={state.algorithm}
          keyCount={state.keys.length}
          onAlgorithmChange={handleAlgorithmChange}
          onAddNode={handleAddNode}
          onRemoveNode={handleRemoveNode}
          onRegenerateKeys={handleRegenerateKeys}
          onKeyCountChange={handleKeyCountChange}
          nodes={state.nodes}
          onShowEducation={() => setShowEducation(true)}
          chblConfig={state.chblConfig}
          onCHBLConfigChange={handleCHBLConfigChange}
        />
      </div>

      {showStats && state.stats && (
        <StatsPanel
          stats={state.stats}
          onClose={() => setShowStats(false)}
        />
      )}

      {showEducation && (
        <EducationPanel
          algorithm={state.algorithm as Algorithm}
          onClose={() => setShowEducation(false)}
        />
      )}

      {showWelcome && (
        <WelcomeGuide
          onClose={() => {
            setShowWelcome(false);
            localStorage.setItem('hasSeenWelcome', 'true');
          }}
        />
      )}

      {showComparison && (
        <ComparisonPanel
          onClose={() => setShowComparison(false)}
        />
      )}

      {showDynamicComparison && dynamicComparisonOp && (
        <DynamicComparisonPanel
          operation={dynamicComparisonOp}
          nodeId={dynamicComparisonNodeId}
          onClose={() => {
            setShowDynamicComparison(false);
            setDynamicComparisonOp(null);
            setDynamicComparisonNodeId(undefined);
          }}
        />
      )}
    </div>
  );
};

export default App;

