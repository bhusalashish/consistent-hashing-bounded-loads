import React, { useState } from 'react';
import { Algorithm } from '../types';
import './EducationPanel.css';

interface EducationPanelProps {
  algorithm: Algorithm;
  onClose: () => void;
}

const EducationPanel: React.FC<EducationPanelProps> = ({ algorithm, onClose }) => {
  const [activeTab, setActiveTab] = useState<'how' | 'why' | 'compare'>('how');

  const algorithmInfo = {
    ring: {
      name: 'Ring Consistent Hash',
      how: [
        '1. Each node gets multiple virtual nodes (vnodes) placed on a hash ring',
        '2. Each key is hashed to a position on the ring (0 to 2^64)',
        '3. The key is assigned to the first node encountered when moving clockwise from the key\'s position',
        '4. This creates a circular mapping where keys "wrap around" the ring',
      ],
      why: [
        'When you add a node: New vnodes are added to the ring, and only keys between the new vnodes and previous nodes move',
        'When you remove a node: Keys assigned to that node\'s vnodes move to the next node clockwise',
        'Distribution: Keys are distributed based on hash positions, which can be uneven if hash values cluster',
        'Churn: ~1/N keys move when adding a node (where N is the number of nodes)',
      ],
      characteristics: {
        churn: 'Medium (~1/N when adding node)',
        distribution: 'Can be uneven',
        complexity: 'O(log N) lookup',
        useCase: 'Simple distributed systems, caching',
      },
    },
    jump: {
      name: 'Jump Consistent Hashing',
      how: [
        '1. Uses a mathematical function to map keys directly to bucket indices',
        '2. No ring structure - uses a deterministic "jump" algorithm',
        '3. Given a key hash, it computes which bucket (node) it belongs to',
        '4. The algorithm ensures minimal remapping when nodes are added/removed',
      ],
      why: [
        'When you add a node: Only ~1/(N+1) keys need to move (minimal churn!)',
        'When you remove a node: Keys from that node redistribute evenly',
        'Distribution: Very uniform - each node gets approximately the same number of keys',
        'Churn: Minimal - only about 1/(N+1) of keys move when adding a node',
        'Why it\'s special: The jump function ensures keys only move when necessary',
      ],
      characteristics: {
        churn: 'Minimal (~1/(N+1))',
        distribution: 'Very uniform',
        complexity: 'O(1) lookup',
        useCase: 'Google Bigtable, Cloud Pub/Sub',
      },
    },
    maglev: {
      name: 'Maglev Load Balancing',
      how: [
        '1. Creates a lookup table (permutation table) of size M (typically a prime like 65537)',
        '2. Each node gets a permutation (unique ordering) of table entries',
        '3. The table is filled by iterating through each node\'s permutation',
        '4. Keys are hashed to a table slot, which points to a node',
      ],
      why: [
        'When you add a node: The table is rebuilt, but most keys stay with the same node',
        'When you remove a node: Keys are reassigned, but distribution remains very uniform',
        'Distribution: Extremely uniform - the permutation ensures even spread',
        'Churn: Low - table rebuilds are efficient and most assignments stay stable',
        'Why it\'s special: The permutation table guarantees near-perfect load balancing',
      ],
      characteristics: {
        churn: 'Low',
        distribution: 'Near-perfect uniform',
        complexity: 'O(1) lookup',
        useCase: 'Google frontend load balancers',
      },
    },
    chbl: {
      name: 'Consistent Hashing with Bounded Loads (CH-BL)',
      how: [
        '1. Uses a virtual node ring (like Ring CH)',
        '2. Each node has a capacity limit: C = ceil(c × expected_keys / num_nodes)',
        '3. When assigning a key: Find the successor node on the ring',
        '4. If that node is at capacity, walk clockwise to find the next available node',
        '5. Uses "two-choice" fallback to avoid long walks',
      ],
      why: [
        'When you add a node: Keys redistribute, but no node exceeds its capacity',
        'When you remove a node: Keys from removed node redistribute, respecting capacity limits',
        'Distribution: Guaranteed to be bounded - no node gets more than c × average',
        'Churn: Similar to Ring CH, but with load guarantees',
        'Why it\'s special: Guarantees max load ≤ c × average load, even with skewed key distribution',
      ],
      characteristics: {
        churn: 'Medium (similar to Ring)',
        distribution: 'Bounded (guaranteed max load)',
        complexity: 'O(1) average, O(N) worst case',
        useCase: 'Caching, storage systems with capacity limits',
      },
    },
  };

  const info = algorithmInfo[algorithm];

  return (
    <div className="education-panel-overlay" onClick={onClose}>
      <div className="education-panel" onClick={(e) => e.stopPropagation()}>
        <div className="education-header">
          <div className="education-title">
            <h2>{info.name}</h2>
            <span className="education-subtitle">How it works under the hood</span>
          </div>
          <button className="education-close-btn" onClick={onClose} aria-label="Close">
            ✕
          </button>
        </div>

        <div className="education-tabs">
          <button
            className={`education-tab ${activeTab === 'how' ? 'active' : ''}`}
            onClick={() => setActiveTab('how')}
          >
            How It Works
          </button>
          <button
            className={`education-tab ${activeTab === 'why' ? 'active' : ''}`}
            onClick={() => setActiveTab('why')}
          >
            Why & When
          </button>
          <button
            className={`education-tab ${activeTab === 'compare' ? 'active' : ''}`}
            onClick={() => setActiveTab('compare')}
          >
            Comparison
          </button>
        </div>

        <div className="education-content">
          {activeTab === 'how' && (
            <div className="education-section">
              <h3>Step-by-Step Process</h3>
              <ol className="education-steps">
                {info.how.map((step, idx) => (
                  <li key={idx}>{step}</li>
                ))}
              </ol>
              <div className="education-visual">
                <div className="visual-box">
                  <strong>Visual Example:</strong>
                  <p>
                    {algorithm === 'ring' &&
                      'Imagine a clock face. Keys land at random positions. Each node has multiple "virtual" positions. The key belongs to the first node you encounter going clockwise.'}
                    {algorithm === 'jump' &&
                      'Think of it like a mathematical formula. Given a key, the jump function directly calculates which node it belongs to, without needing a ring structure.'}
                    {algorithm === 'maglev' &&
                      'Picture a lookup table with thousands of slots. Each slot is pre-assigned to a node using a special permutation. Keys hash to a slot, which tells you the node.'}
                    {algorithm === 'chbl' &&
                      'Similar to Ring CH, but with capacity limits. If a node is "full", the algorithm walks around the ring to find the next available node.'}
                  </p>
                </div>
              </div>
            </div>
          )}

          {activeTab === 'why' && (
            <div className="education-section">
              <h3>Rebalancing Behavior</h3>
              <ul className="education-list">
                {info.why.map((item, idx) => (
                  <li key={idx}>{item}</li>
                ))}
              </ul>
              <div className="education-characteristics">
                <h4>Key Characteristics</h4>
                <div className="characteristics-grid">
                  <div className="characteristic-item">
                    <span className="char-label">Churn:</span>
                    <span className="char-value">{info.characteristics.churn}</span>
                  </div>
                  <div className="characteristic-item">
                    <span className="char-label">Distribution:</span>
                    <span className="char-value">{info.characteristics.distribution}</span>
                  </div>
                  <div className="characteristic-item">
                    <span className="char-label">Lookup Time:</span>
                    <span className="char-value">{info.characteristics.complexity}</span>
                  </div>
                  <div className="characteristic-item">
                    <span className="char-label">Used In:</span>
                    <span className="char-value">{info.characteristics.useCase}</span>
                  </div>
                </div>
              </div>
            </div>
          )}

          {activeTab === 'compare' && (
            <div className="education-section">
              <h3>Algorithm Comparison</h3>
              <div className="comparison-table">
                <div className="comparison-row header">
                  <div className="comparison-cell">Algorithm</div>
                  <div className="comparison-cell">Churn (Add Node)</div>
                  <div className="comparison-cell">Distribution</div>
                  <div className="comparison-cell">Best For</div>
                </div>
                {Object.entries(algorithmInfo).map(([algo, info]) => (
                  <div key={algo} className={`comparison-row ${algo === algorithm ? 'highlight' : ''}`}>
                    <div className="comparison-cell">
                      <strong>{info.name}</strong>
                    </div>
                    <div className="comparison-cell">{info.characteristics.churn}</div>
                    <div className="comparison-cell">{info.characteristics.distribution}</div>
                    <div className="comparison-cell">{info.characteristics.useCase}</div>
                  </div>
                ))}
              </div>
              <div className="education-tip">
                <strong>💡 Tip:</strong> Try switching between algorithms and watch how the key distribution changes!
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default EducationPanel;

