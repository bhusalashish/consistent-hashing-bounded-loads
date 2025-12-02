import React, { useState } from 'react';
import './ComparisonPanel.css';

interface ComparisonPanelProps {
  onClose: () => void;
}

const ComparisonPanel: React.FC<ComparisonPanelProps> = ({ onClose }) => {
  const [selectedOperation, setSelectedOperation] = useState<'add-node' | 'remove-node' | 'new-key'>('add-node');

  const operations = {
    'add-node': {
      title: 'Adding a New Node',
      description: 'What happens when a new node joins the cluster?',
    },
    'remove-node': {
      title: 'Removing a Node',
      description: 'What happens when a node leaves the cluster?',
    },
    'new-key': {
      title: 'Adding New Keys',
      description: 'How are new keys distributed across nodes?',
    },
  };

  const algorithms = {
    ring: {
      name: 'Ring Consistent Hash',
      'add-node': {
        behavior: 'Keys between the new node\'s vnodes and previous nodes move to the new node',
        churn: '~1/N keys move (where N = number of nodes)',
        visualization: 'New vnodes are added to the ring. Keys hash to positions, and those between old and new vnodes move.',
        example: 'With 4 nodes, adding a 5th node moves ~20% of keys',
      },
      'remove-node': {
        behavior: 'All keys from the removed node move to the next node clockwise',
        churn: 'All keys from removed node move',
        visualization: 'Node\'s vnodes are removed. Keys that were assigned to those vnodes now go to the next node.',
        example: 'Removing 1 node out of 5 moves all keys from that node (~20% of total)',
      },
      'new-key': {
        behavior: 'Key hashes to a position on the ring, assigned to first node clockwise',
        churn: 'No existing keys move',
        visualization: 'Key position determined by hash. Walk clockwise to find first node.',
        example: 'New key at 60% position goes to the node whose vnode is first after 60%',
      },
    },
    jump: {
      name: 'Jump Consistent Hashing',
      'add-node': {
        behavior: 'Jump function recalculates assignments. Only keys that need to move do so.',
        churn: '~1/(N+1) keys move (minimal churn!)',
        visualization: 'Mathematical function directly maps key to node index. Adding a node only affects keys that hash to the new node.',
        example: 'With 4 nodes, adding a 5th node moves only ~20% of keys (1/5)',
      },
      'remove-node': {
        behavior: 'Keys from removed node redistribute evenly to remaining nodes',
        churn: 'All keys from removed node move',
        visualization: 'Jump function ensures even redistribution. Keys are reassigned based on new node count.',
        example: 'Removing 1 node out of 5 redistributes its keys evenly to remaining 4 nodes',
      },
      'new-key': {
        behavior: 'Jump function calculates which node the key belongs to directly',
        churn: 'No existing keys move',
        visualization: 'No ring structure. Formula: bucket = jump(hash, numBuckets)',
        example: 'Key hash + number of nodes → direct node assignment',
      },
    },
    maglev: {
      name: 'Maglev Load Balancing',
      'add-node': {
        behavior: 'Lookup table is rebuilt. Most keys stay with same node due to stable permutations.',
        churn: 'Low - table rebuild is efficient',
        visualization: 'Each node gets a permutation of table slots. Table is filled by iterating permutations. Most assignments stay stable.',
        example: 'Adding a node rebuilds the table, but ~80-90% of keys stay with same node',
      },
      'remove-node': {
        behavior: 'Table is rebuilt. Keys are reassigned but distribution remains uniform.',
        churn: 'All keys from removed node move',
        visualization: 'Permutation table ensures uniform distribution even after rebuild.',
        example: 'Removing a node redistributes its keys uniformly across remaining nodes',
      },
      'new-key': {
        behavior: 'Key hashes to a table slot, which points to a pre-assigned node',
        churn: 'No existing keys move',
        visualization: 'Key → hash → table slot → node. O(1) lookup.',
        example: 'Key hashes to slot 12345, which was assigned to node-2 during table construction',
      },
    },
    chbl: {
      name: 'CH-BL (Bounded Loads)',
      'add-node': {
        behavior: 'Keys redistribute, but no node exceeds capacity. New node gets keys up to its capacity limit.',
        churn: 'Similar to Ring CH, but respects capacity bounds',
        visualization: 'Keys hash to ring, walk clockwise. If nodes are at capacity, continue walking to find available node.',
        example: 'New node gets keys until it reaches capacity (c × avg). Remaining keys stay with other nodes.',
      },
      'remove-node': {
        behavior: 'Keys from removed node redistribute, but no node exceeds its capacity limit',
        churn: 'All keys from removed node move',
        visualization: 'Keys walk the ring clockwise, skipping nodes at capacity until finding available space.',
        example: 'Removed node\'s keys walk to next available nodes, respecting each node\'s capacity',
      },
      'new-key': {
        behavior: 'Key hashes to ring, walks clockwise until finding a node with available capacity',
        churn: 'No existing keys move',
        visualization: 'Hash → ring position → walk clockwise → check capacity → assign or continue walking',
        example: 'Key at 70% position. Node at 75% is full, so key goes to node at 80% which has capacity',
      },
    },
  };

  const currentOp = operations[selectedOperation];

  return (
    <div className="comparison-overlay" onClick={onClose}>
      <div className="comparison-panel" onClick={(e) => e.stopPropagation()}>
        <div className="comparison-header">
          <div>
            <h2>Algorithm Comparison</h2>
            <p>Side-by-side comparison of how algorithms handle operations</p>
          </div>
          <button className="comparison-close-btn" onClick={onClose} aria-label="Close">
            ✕
          </button>
        </div>

        <div className="comparison-tabs">
          {Object.entries(operations).map(([key, op]) => (
            <button
              key={key}
              className={`comparison-tab ${selectedOperation === key ? 'active' : ''}`}
              onClick={() => setSelectedOperation(key as any)}
            >
              {op.title}
            </button>
          ))}
        </div>

        <div className="comparison-content">
          <div className="operation-description">
            <h3>{currentOp.title}</h3>
            <p>{currentOp.description}</p>
          </div>

          <div className="algorithm-comparisons">
            {Object.entries(algorithms).map(([algoKey, algo]) => {
              const behavior = algo[selectedOperation];
              return (
                <div key={algoKey} className="algorithm-card">
                  <div className="algorithm-card-header">
                    <h4>{algo.name}</h4>
                    <span className={`algorithm-badge ${algoKey}`}>{algoKey.toUpperCase()}</span>
                  </div>
                  <div className="algorithm-details">
                    <div className="detail-section">
                      <strong>Behavior:</strong>
                      <p>{behavior.behavior}</p>
                    </div>
                    <div className="detail-section">
                      <strong>Churn:</strong>
                      <p className="churn-info">{behavior.churn}</p>
                    </div>
                    <div className="detail-section">
                      <strong>How it works:</strong>
                      <p>{behavior.visualization}</p>
                    </div>
                    <div className="detail-section">
                      <strong>Example:</strong>
                      <p className="example-text">{behavior.example}</p>
                    </div>
                  </div>
                </div>
              );
            })}
          </div>

          <div className="comparison-summary">
            <h4>Key Takeaways</h4>
            <ul>
              {selectedOperation === 'add-node' && (
                <>
                  <li><strong>Jump</strong> has the lowest churn (~1/(N+1))</li>
                  <li><strong>Ring</strong> and <strong>CH-BL</strong> have similar churn (~1/N)</li>
                  <li><strong>Maglev</strong> has low churn due to stable permutations</li>
                  <li><strong>CH-BL</strong> is the only one that enforces capacity limits</li>
                </>
              )}
              {selectedOperation === 'remove-node' && (
                <>
                  <li>All algorithms move all keys from the removed node</li>
                  <li><strong>Jump</strong> ensures even redistribution</li>
                  <li><strong>Maglev</strong> maintains uniform distribution after rebuild</li>
                  <li><strong>CH-BL</strong> respects capacity bounds during redistribution</li>
                </>
              )}
              {selectedOperation === 'new-key' && (
                <>
                  <li>Adding new keys never moves existing keys</li>
                  <li><strong>Jump</strong> has O(1) direct calculation</li>
                  <li><strong>Maglev</strong> has O(1) table lookup</li>
                  <li><strong>Ring</strong> and <strong>CH-BL</strong> use ring traversal</li>
                  <li><strong>CH-BL</strong> may need to walk if nodes are at capacity</li>
                </>
              )}
            </ul>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ComparisonPanel;

