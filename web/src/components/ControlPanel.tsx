import React, { useState } from 'react';
import { Algorithm } from '../types';
import './ControlPanel.css';

interface ControlPanelProps {
  algorithm: string;
  keyCount: number;
  onAlgorithmChange: (algo: Algorithm) => void;
  onAddNode: () => void;
  onRemoveNode: (nodeId: string) => void;
  onRegenerateKeys: () => void;
  onKeyCountChange: (count: number) => void;
  nodes: string[];
  onShowEducation: () => void;
  chblConfig?: { loadFactor: number; expectedKeys: number; capacityPerNode: number };
  onCHBLConfigChange?: (loadFactor: number, expectedKeys: number) => void;
}

const ControlPanel: React.FC<ControlPanelProps> = ({
  algorithm,
  keyCount,
  onAlgorithmChange,
  onAddNode,
  onRemoveNode,
  onRegenerateKeys,
  onKeyCountChange,
  nodes,
  onShowEducation,
  chblConfig,
  onCHBLConfigChange,
}) => {
  const [localKeyCount, setLocalKeyCount] = useState(keyCount);
  const [localLoadFactor, setLocalLoadFactor] = useState(chblConfig?.loadFactor || 1.25);
  const [localExpectedKeys, setLocalExpectedKeys] = useState(chblConfig?.expectedKeys || 1000);

  const handleKeyCountSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onKeyCountChange(localKeyCount);
  };

  const handleCHBLConfigSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (onCHBLConfigChange) {
      onCHBLConfigChange(localLoadFactor, localExpectedKeys);
    }
  };

  // Update local state when chblConfig changes
  React.useEffect(() => {
    if (chblConfig) {
      setLocalLoadFactor(chblConfig.loadFactor);
      setLocalExpectedKeys(chblConfig.expectedKeys);
    }
  }, [chblConfig]);

  return (
    <div className="control-panel">
      <h2>Controls</h2>

      <div className="control-group">
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '8px' }}>
          <label htmlFor="algorithm">Algorithm:</label>
          <button
            onClick={onShowEducation}
            className="btn-education"
            title="Learn how this algorithm works"
          >
            📚 Learn
          </button>
        </div>
        <select
          id="algorithm"
          value={algorithm}
          onChange={(e) => onAlgorithmChange(e.target.value as Algorithm)}
        >
          <option value="ring">Ring Consistent Hash</option>
          <option value="jump">Jump Consistent Hashing</option>
          <option value="maglev">Maglev</option>
          <option value="chbl">CH-BL (Bounded Loads)</option>
        </select>
      </div>

      <div className="control-group">
        <label>Key Count:</label>
        <form onSubmit={handleKeyCountSubmit} className="key-count-form">
          <input
            type="number"
            min="0"
            max="1000"
            value={localKeyCount}
            onChange={(e) => setLocalKeyCount(Number(e.target.value))}
          />
          <button type="submit">Update</button>
        </form>
      </div>

      <div className="control-group">
        <button onClick={onRegenerateKeys} className="btn-primary">
          Regenerate Keys
        </button>
      </div>

      {algorithm === 'chbl' && chblConfig && onCHBLConfigChange && (
        <div className="control-group chbl-config">
          <h3>CH-BL Configuration</h3>
          <div className="config-explanation">
            <p>
              <strong>Capacity per node:</strong> {chblConfig.capacityPerNode} keys
            </p>
            <p className="config-formula">
              Capacity = ⌈Load Factor × (Expected Keys / Nodes)⌉
            </p>
            <p className="config-hint">
              Adjust these values to see how capacity limits affect key distribution.
              Higher load factor = more capacity per node.
            </p>
          </div>
          <form onSubmit={handleCHBLConfigSubmit} className="chbl-config-form">
            <div className="config-input-group">
              <label htmlFor="loadFactor">Load Factor (c):</label>
              <input
                id="loadFactor"
                type="number"
                min="1.0"
                max="3.0"
                step="0.1"
                value={localLoadFactor}
                onChange={(e) => setLocalLoadFactor(Number(e.target.value))}
              />
              <span className="config-hint-small">Typical: 1.1 - 1.5</span>
            </div>
            <div className="config-input-group">
              <label htmlFor="expectedKeys">Expected Keys:</label>
              <input
                id="expectedKeys"
                type="number"
                min="100"
                max="10000"
                step="100"
                value={localExpectedKeys}
                onChange={(e) => setLocalExpectedKeys(Number(e.target.value))}
              />
              <span className="config-hint-small">Total keys you expect</span>
            </div>
            <button type="submit" className="btn-primary">
              Update CH-BL Config
            </button>
          </form>
        </div>
      )}

      <div className="control-group">
        <h3>Nodes ({nodes.length})</h3>
        <button onClick={onAddNode} className="btn-success">
          + Add Node
        </button>
        <div className="node-list">
          {nodes.map((nodeId) => (
            <div key={nodeId} className="node-item">
              <span>{nodeId}</span>
              <button
                onClick={() => onRemoveNode(nodeId)}
                className="btn-danger btn-small"
                disabled={nodes.length <= 1}
              >
                Remove
              </button>
            </div>
          ))}
        </div>
      </div>

      <div className="info-box">
        <h3>Info</h3>
        <p>
          <strong>Algorithm:</strong> {algorithm}
        </p>
        <p>
          <strong>Nodes:</strong> {nodes.length}
        </p>
        <p>
          <strong>Keys:</strong> {keyCount}
        </p>
      </div>
    </div>
  );
};

export default ControlPanel;

