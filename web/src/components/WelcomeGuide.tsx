import React from 'react';
import './WelcomeGuide.css';

interface WelcomeGuideProps {
  onClose: () => void;
}

const WelcomeGuide: React.FC<WelcomeGuideProps> = ({ onClose }) => {
  return (
    <div className="welcome-overlay" onClick={onClose}>
      <div className="welcome-panel" onClick={(e) => e.stopPropagation()}>
        <div className="welcome-header">
          <h2>Welcome to Consistent Hashing Visualizer! 🎓</h2>
          <button className="welcome-close-btn" onClick={onClose} aria-label="Close">
            ✕
          </button>
        </div>

        <div className="welcome-content">
          <div className="welcome-section">
            <h3>🎯 What is this?</h3>
            <p>
              This is an interactive playground to understand how consistent hashing algorithms work.
              Perfect for students, engineers, or anyone curious about distributed systems!
            </p>
          </div>

          <div className="welcome-section">
            <h3>🚀 Quick Start</h3>
            <ol className="welcome-steps">
              <li>
                <strong>Explore algorithms:</strong> Click the "📚 Learn" button next to the algorithm dropdown to understand how each one works
              </li>
              <li>
                <strong>Hover over keys/nodes:</strong> See why keys are assigned to specific nodes
              </li>
              <li>
                <strong>Add/Remove nodes:</strong> Watch how keys rebalance and see statistics
              </li>
              <li>
                <strong>Switch algorithms:</strong> Compare how different algorithms distribute keys
              </li>
            </ol>
          </div>

          <div className="welcome-section">
            <h3>💡 Key Features</h3>
            <ul className="welcome-features">
              <li>
                <strong>Visual explanations:</strong> Hover over any key or node to see why assignments happen
              </li>
              <li>
                <strong>Statistics:</strong> After each operation, see detailed stats about key movements
              </li>
              <li>
                <strong>Algorithm education:</strong> Learn the step-by-step process for each algorithm
              </li>
              <li>
                <strong>Real-time visualization:</strong> Watch keys animate as they move during rebalancing
              </li>
            </ul>
          </div>

          <div className="welcome-section">
            <h3>🎓 Learning Path</h3>
            <div className="learning-path">
              <div className="path-step">
                <div className="path-number">1</div>
                <div className="path-content">
                  <strong>Start with Ring CH</strong>
                  <p>Understand the basic ring structure and how keys map to nodes</p>
                </div>
              </div>
              <div className="path-step">
                <div className="path-number">2</div>
                <div className="path-content">
                  <strong>Try Jump Hashing</strong>
                  <p>See how minimal churn is achieved with mathematical precision</p>
                </div>
              </div>
              <div className="path-step">
                <div className="path-number">3</div>
                <div className="path-content">
                  <strong>Explore Maglev</strong>
                  <p>Discover how permutation tables ensure uniform distribution</p>
                </div>
              </div>
              <div className="path-step">
                <div className="path-number">4</div>
                <div className="path-content">
                  <strong>Understand CH-BL</strong>
                  <p>Learn how bounded loads guarantee capacity limits</p>
                </div>
              </div>
            </div>
          </div>

          <div className="welcome-tip">
            <strong>💡 Pro Tip:</strong> Add a node and watch the statistics panel to see exactly how many keys moved and why!
          </div>
        </div>

        <div className="welcome-footer">
          <button className="welcome-start-btn" onClick={onClose}>
            Let's Start Exploring! 🚀
          </button>
        </div>
      </div>
    </div>
  );
};

export default WelcomeGuide;

