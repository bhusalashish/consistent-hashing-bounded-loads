import React, { useState, useEffect } from 'react';
import { AlgorithmComparison } from '../types';
import * as api from '../api';
import './DynamicComparisonPanel.css';

interface DynamicComparisonPanelProps {
  operation: 'add-node' | 'remove-node' | 'regenerate-keys';
  nodeId?: string;
  onClose: () => void;
}

const DynamicComparisonPanel: React.FC<DynamicComparisonPanelProps> = ({ operation, nodeId, onClose }) => {
  const [comparisons, setComparisons] = useState<AlgorithmComparison[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchComparison = async () => {
      try {
        setLoading(true);
        setError(null);
        const results = await api.compareOperation(operation, nodeId);
        setComparisons(results);
        setLoading(false);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to compare operations');
        setLoading(false);
      }
    };

    fetchComparison();
  }, [operation, nodeId]);

  const getOperationTitle = () => {
    switch (operation) {
      case 'add-node':
        return 'Adding a New Node';
      case 'remove-node':
        return 'Removing a Node';
      case 'regenerate-keys':
        return 'Regenerating Keys';
      default:
        return 'Operation Comparison';
    }
  };

  const algorithmNames: Record<string, string> = {
    ring: 'Ring Consistent Hash',
    jump: 'Jump Consistent Hashing',
    maglev: 'Maglev',
    chbl: 'CH-BL (Bounded Loads)',
  };

  if (loading) {
    return (
      <div className="dynamic-comparison-overlay" onClick={onClose}>
        <div className="dynamic-comparison-panel" onClick={(e) => e.stopPropagation()}>
          <div className="comparison-loading">
            <div className="spinner"></div>
            <p>Running operation on all algorithms...</p>
          </div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="dynamic-comparison-overlay" onClick={onClose}>
        <div className="dynamic-comparison-panel" onClick={(e) => e.stopPropagation()}>
          <div className="comparison-error">
            <h3>Error</h3>
            <p>{error}</p>
            <button onClick={onClose}>Close</button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="dynamic-comparison-overlay" onClick={onClose}>
      <div className="dynamic-comparison-panel" onClick={(e) => e.stopPropagation()}>
        <div className="dynamic-comparison-header">
          <div>
            <h2>{getOperationTitle()}</h2>
            <p>Real-time comparison across all algorithms</p>
          </div>
          <button className="dynamic-comparison-close-btn" onClick={onClose} aria-label="Close">
            ✕
          </button>
        </div>

        <div className="dynamic-comparison-content">
          <div className="comparison-grid">
            {comparisons.map((comp) => (
              <div key={comp.algorithm} className="comparison-algorithm-card">
                <div className="comparison-card-header">
                  <h3>{algorithmNames[comp.algorithm] || comp.algorithm}</h3>
                  <span className={`algorithm-badge ${comp.algorithm}`}>
                    {comp.algorithm.toUpperCase()}
                  </span>
                </div>

                {comp.stats && (
                  <div className="comparison-stats">
                    <div className="stat-row">
                      <span className="stat-label">Keys Moved:</span>
                      <span className="stat-value highlight">
                        {comp.stats.keysMoved} ({comp.stats.keysMovedPercent.toFixed(2)}%)
                      </span>
                    </div>
                    <div className="stat-row">
                      <span className="stat-label">Total Keys:</span>
                      <span className="stat-value">{comp.stats.totalKeys}</span>
                    </div>
                  </div>
                )}

                <div className="comparison-distribution">
                  <h4>Key Distribution</h4>
                  <div className="distribution-list">
                    {Object.entries(comp.stats?.distribution || {}).map(([nodeId, count]) => {
                      const prevCount = comp.stats?.previousDist[nodeId] || 0;
                      const change = count - prevCount;
                      return (
                        <div key={nodeId} className="distribution-item">
                          <div className="dist-node">{nodeId}</div>
                          <div className="dist-values">
                            <span className="dist-current">{count}</span>
                            {change !== 0 && (
                              <span className={`dist-change ${change > 0 ? 'positive' : 'negative'}`}>
                                {change > 0 ? '+' : ''}{change}
                              </span>
                            )}
                          </div>
                        </div>
                      );
                    })}
                  </div>
                </div>

                {comp.stats && comp.stats.movementByNode && Object.keys(comp.stats.movementByNode).length > 0 && (
                  <div className="comparison-movements">
                    <h4>Keys Received by Node</h4>
                    <div className="movement-bars">
                      {Object.entries(comp.stats.movementByNode)
                        .sort(([, a], [, b]) => b - a)
                        .map(([nodeId, count]) => (
                          <div key={nodeId} className="movement-bar-item">
                            <span className="movement-node">{nodeId}</span>
                            <div className="movement-bar-container">
                              <div
                                className="movement-bar-fill"
                                style={{ width: `${(count / (comp.stats?.keysMoved || 1)) * 100}%` }}
                              >
                                {count}
                              </div>
                            </div>
                          </div>
                        ))}
                    </div>
                  </div>
                )}

                {comp.state.chblConfig && (
                  <div className="chbl-info">
                    <h4>CH-BL Configuration</h4>
                    <div className="chbl-details">
                      <div>Capacity: {comp.state.chblConfig.capacityPerNode} keys/node</div>
                      <div>Load Factor: {comp.state.chblConfig.loadFactor}</div>
                      <div>Expected Keys: {comp.state.chblConfig.expectedKeys}</div>
                    </div>
                  </div>
                )}
              </div>
            ))}
          </div>

          <div className="comparison-summary-box">
            <h3>Key Insights</h3>
            <ul>
              {comparisons.map((comp) => {
                const churn = comp.stats?.keysMovedPercent || 0;
                return (
                  <li key={comp.algorithm}>
                    <strong>{algorithmNames[comp.algorithm]}:</strong>{' '}
                    {churn.toFixed(2)}% churn ({comp.stats?.keysMoved || 0} keys moved)
                  </li>
                );
              })}
            </ul>
            <div className="insight-note">
              💡 Lower churn means fewer keys need to move, which is better for system stability.
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default DynamicComparisonPanel;

