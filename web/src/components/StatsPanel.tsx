import React from 'react';
import { Statistics } from '../types';
import './StatsPanel.css';

interface StatsPanelProps {
  stats: Statistics | null;
  onClose: () => void;
}

const StatsPanel: React.FC<StatsPanelProps> = ({ stats, onClose }) => {
  if (!stats) return null;

  const getOperationLabel = (op: string) => {
    const labels: Record<string, string> = {
      'add-node': 'Node Added',
      'remove-node': 'Node Removed',
      'regenerate-keys': 'Keys Regenerated',
      'set-algorithm': 'Algorithm Changed',
      'set-key-count': 'Key Count Changed',
    };
    return labels[op] || op;
  };

  const getOperationIcon = (op: string) => {
    if (op === 'add-node') return '➕';
    if (op === 'remove-node') return '➖';
    if (op === 'regenerate-keys') return '🔄';
    if (op === 'set-algorithm') return '⚙️';
    return '📊';
  };

  return (
    <div className="stats-panel-overlay" onClick={onClose}>
      <div className="stats-panel" onClick={(e) => e.stopPropagation()}>
        <div className="stats-panel-header">
          <div className="stats-panel-title">
            <span className="stats-icon">{getOperationIcon(stats.operation)}</span>
            <h2>{getOperationLabel(stats.operation)}</h2>
          </div>
          <button className="stats-close-btn" onClick={onClose} aria-label="Close">
            ✕
          </button>
        </div>

        <div className="stats-content">
          {/* Summary Cards */}
          <div className="stats-summary">
            <div className="stat-card">
              <div className="stat-label">Total Keys</div>
              <div className="stat-value">{stats.totalKeys.toLocaleString()}</div>
            </div>
            <div className="stat-card">
              <div className="stat-label">Keys Moved</div>
              <div className="stat-value highlight">{stats.keysMoved.toLocaleString()}</div>
            </div>
            <div className="stat-card">
              <div className="stat-label">Movement %</div>
              <div className="stat-value highlight">
                {stats.keysMovedPercent.toFixed(2)}%
              </div>
            </div>
          </div>

          {/* Distribution Comparison */}
          <div className="stats-section">
            <h3>Key Distribution</h3>
            <div className="distribution-list">
              {Object.entries(stats.distribution).map(([nodeId, count]) => {
                const prevCount = stats.previousDist[nodeId] || 0;
                const change = count - prevCount;
                const changePercent =
                  prevCount > 0 ? ((change / prevCount) * 100).toFixed(1) : '∞';

                return (
                  <div key={nodeId} className="distribution-item">
                    <div className="distribution-node">{nodeId}</div>
                    <div className="distribution-values">
                      <div className="distribution-count">
                        <span className="count-label">Current:</span>
                        <span className="count-value">{count}</span>
                      </div>
                      <div className="distribution-count">
                        <span className="count-label">Previous:</span>
                        <span className="count-value previous">{prevCount}</span>
                      </div>
                      <div className="distribution-change">
                        {change !== 0 && (
                          <span className={change > 0 ? 'change positive' : 'change negative'}>
                            {change > 0 ? '+' : ''}
                            {change} ({changePercent}%)
                          </span>
                        )}
                        {change === 0 && (
                          <span className="change neutral">No change</span>
                        )}
                      </div>
                    </div>
                  </div>
                );
              })}
            </div>
          </div>

          {/* Movement by Node */}
          {stats.keysMoved > 0 && (
            <div className="stats-section">
              <h3>Keys Received by Node</h3>
              <div className="movement-list">
                {Object.entries(stats.movementByNode)
                  .sort(([, a], [, b]) => b - a)
                  .map(([nodeId, count]) => (
                    <div key={nodeId} className="movement-item">
                      <div className="movement-node">{nodeId}</div>
                      <div className="movement-bar-container">
                        <div
                          className="movement-bar"
                          style={{
                            width: `${(count / stats.keysMoved) * 100}%`,
                          }}
                        >
                          {count}
                        </div>
                      </div>
                    </div>
                  ))}
              </div>
            </div>
          )}

          {/* Key Movements Detail */}
          {stats.keyMovements.length > 0 && (
            <div className="stats-section">
              <h3>
                Key Movements{' '}
                {stats.keyMovements.length < stats.keysMoved && (
                  <span className="stats-subtitle">
                    (showing first {stats.keyMovements.length} of {stats.keysMoved})
                  </span>
                )}
              </h3>
              <div className="key-movements-list">
                {stats.keyMovements.map((movement, idx) => (
                  <div key={idx} className="key-movement-item">
                    <div className="key-movement-key">{movement.keyId}</div>
                    <div className="key-movement-arrow">→</div>
                    <div className="key-movement-from">{movement.fromNode}</div>
                    <div className="key-movement-arrow">→</div>
                    <div className="key-movement-to">{movement.toNode}</div>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Capacity Information for CH-BL */}
          {stats.capacityInfo && stats.capacityInfo.unassignedKeys > 0 && (
            <div className="stats-section capacity-warning">
              <h3>⚠️ Capacity Warning</h3>
              <div className="capacity-alert">
                <div className="alert-message">
                  <strong>{stats.capacityInfo.unassignedKeys} keys</strong> could not be assigned
                  because all nodes are at capacity.
                </div>
                <div className="alert-suggestion">
                  💡 <strong>Solution:</strong> Increase <code>ExpectedKeys</code> or{' '}
                  <code>LoadFactor</code> in CH-BL configuration
                </div>
              </div>
            </div>
          )}

          {stats.capacityInfo && (
            <div className="stats-section">
              <h3>CH-BL Capacity Status</h3>
              {stats.capacityInfo.nodesAtCapacity.length > 0 && (
                <div className="capacity-nodes-warning">
                  <strong>Nodes at Capacity:</strong>{' '}
                  {stats.capacityInfo.nodesAtCapacity.join(', ')}
                </div>
              )}
              <div className="capacity-list">
                {Object.entries(stats.capacityInfo.capacityPerNode).map(([nodeId, capacity]) => {
                  const load = stats.capacityInfo!.currentLoad[nodeId] || 0;
                  const percentage = stats.capacityInfo!.loadPercentage[nodeId] || 0;
                  const isAtCapacity = stats.capacityInfo!.nodesAtCapacity.includes(nodeId);

                  return (
                    <div
                      key={nodeId}
                      className={`capacity-item ${isAtCapacity ? 'at-capacity' : ''}`}
                    >
                      <div className="capacity-node">{nodeId}</div>
                      <div className="capacity-details">
                        <div className="capacity-bar-container">
                          <div
                            className="capacity-bar"
                            style={{
                              width: `${Math.min(percentage, 100)}%`,
                              backgroundColor:
                                percentage >= 90
                                  ? '#ea4335'
                                  : percentage >= 70
                                  ? '#fbbc04'
                                  : '#34a853',
                            }}
                          >
                            {load}/{capacity} ({percentage.toFixed(1)}%)
                          </div>
                        </div>
                        {isAtCapacity && (
                          <span className="capacity-badge">AT CAPACITY</span>
                        )}
                      </div>
                    </div>
                  );
                })}
              </div>
            </div>
          )}

          {/* Algorithm-specific insights */}
          {stats.operation === 'set-algorithm' && (
            <div className="stats-section">
              <h3>Algorithm Comparison</h3>
              <div className="insight-box">
                <p>
                  Changing algorithms redistributes keys based on each algorithm's
                  assignment strategy. The movement percentage indicates how much
                  churn occurred during the transition.
                </p>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default StatsPanel;

