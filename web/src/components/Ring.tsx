import React, { useEffect, useRef, useMemo } from 'react';
import * as d3 from 'd3';
import { VisualizerState } from '../types';

interface RingProps {
  state: VisualizerState;
  width: number;
  height: number;
  onNodeHover?: (nodeId: string | null) => void;
  onKeyHover?: (keyId: string | null) => void;
  highlightedKey?: string | null;
}

const Ring: React.FC<RingProps> = ({ state, width, height, onNodeHover, onKeyHover }) => {
  const svgRef = useRef<SVGSVGElement>(null);
  const centerX = width / 2;
  const centerY = height / 2;
  const radius = Math.min(width, height) * 0.32;
  const keyRadius = radius + 35;

  // Group keys by node for better visualization
  const keysByNode = useMemo(() => {
    const groups: Record<string, string[]> = {};
    state.keys.forEach((keyId) => {
      const node = state.assignments[keyId];
      if (!groups[node]) {
        groups[node] = [];
      }
      groups[node].push(keyId);
    });
    return groups;
  }, [state.keys, state.assignments]);

  // Determine if we should use aggregated visualization
  const useAggregated = state.keys.length > 50;
  const keysPerSegment = useAggregated ? Math.ceil(state.keys.length / 360) : 1;

  useEffect(() => {
    if (!svgRef.current) return;

    const svg = d3.select(svgRef.current);
    svg.selectAll('*').remove();

    // Create main group
    const g = svg.append('g')
      .attr('transform', `translate(${centerX}, ${centerY})`);

    // Draw ring circle with subtle styling
    g.append('circle')
      .attr('r', radius)
      .attr('fill', 'none')
      .attr('stroke', '#dadce0')
      .attr('stroke-width', 1.5)
      .attr('stroke-dasharray', '4,4')
      .attr('opacity', 0.6);

    // Draw position markers (0%, 25%, 50%, 75%) to show what "position around ring" means
    const positionMarkers = [0, 0.25, 0.5, 0.75];
    positionMarkers.forEach((pos) => {
      const angle = pos * 2 * Math.PI;
      const x = Math.cos(angle) * radius;
      const y = Math.sin(angle) * radius;
      const outerX = Math.cos(angle) * (radius + 15);
      const outerY = Math.sin(angle) * (radius + 15);

      // Marker line
      g.append('line')
        .attr('x1', x)
        .attr('y1', y)
        .attr('x2', outerX)
        .attr('y2', outerY)
        .attr('stroke', '#dadce0')
        .attr('stroke-width', 1)
        .attr('opacity', 0.4);

      // Position label
      g.append('text')
        .attr('x', outerX + Math.cos(angle) * 10)
        .attr('y', outerY + Math.sin(angle) * 10)
        .attr('text-anchor', 'middle')
        .attr('fill', '#80868b')
        .attr('font-size', '10px')
        .attr('font-weight', '500')
        .text(`${(pos * 100).toFixed(0)}%`);
    });

    // Google-style color palette
    const nodeColors = [
      '#4285f4', // Blue
      '#ea4335', // Red
      '#fbbc04', // Yellow
      '#34a853', // Green
      '#ff6d01', // Orange
      '#9334e6', // Purple
      '#e91e63', // Pink
      '#00acc1', // Cyan
    ];

    // Draw nodes with better styling
    state.nodes.forEach((nodeId, idx) => {
      const angle = state.nodeAngles[nodeId] || (2 * Math.PI * idx / state.nodes.length);
      const x = Math.cos(angle) * radius;
      const y = Math.sin(angle) * radius;

      const nodeGroup = g.append('g')
        .attr('class', 'node')
        .attr('data-node-id', nodeId)
        .style('cursor', 'pointer')
        .on('mouseenter', () => onNodeHover?.(nodeId))
        .on('mouseleave', () => onNodeHover?.(null));

      // Node circle with shadow effect
      const nodeCircle = nodeGroup.append('circle')
        .attr('cx', x)
        .attr('cy', y)
        .attr('r', 24)
        .attr('fill', nodeColors[idx % nodeColors.length])
        .attr('stroke', '#fff')
        .attr('stroke-width', 3)
        .attr('opacity', 1.0)
        .style('filter', 'drop-shadow(0 2px 4px rgba(0,0,0,0.2))');

      // Node label with better typography
      nodeGroup.append('text')
        .attr('x', x)
        .attr('y', y + 40)
        .attr('text-anchor', 'middle')
        .attr('fill', '#202124')
        .attr('font-size', '13px')
        .attr('font-weight', '500')
        .text(nodeId);

      // Key count badge
      const keyCount = keysByNode[nodeId]?.length || 0;
      if (keyCount > 0) {
        const badgeText = nodeGroup.append('text')
          .attr('x', x)
          .attr('y', y - 30)
          .attr('text-anchor', 'middle')
          .attr('fill', '#5f6368')
          .attr('font-size', '11px')
          .attr('font-weight', '400')
          .text(`${keyCount} keys`);

        // Add capacity indicator for CH-BL
        if (state.algorithm === 'chbl' && state.chblConfig) {
          const capacity = state.chblConfig.capacityPerNode;
          const usagePercent = (keyCount / capacity) * 100;
          const capacityColor = usagePercent >= 90 ? '#ea4335' : usagePercent >= 70 ? '#fbbc04' : '#34a853';
          
          badgeText
            .text(`${keyCount}/${capacity} keys`)
            .attr('fill', capacityColor);

          // Capacity bar
          const barWidth = 40;
          const barHeight = 4;
          const barX = x - barWidth / 2;
          const barY = y - 40;

          const barBg = nodeGroup.append('rect')
            .attr('x', barX)
            .attr('y', barY)
            .attr('width', barWidth)
            .attr('height', barHeight)
            .attr('fill', '#e8eaed')
            .attr('rx', 2);

          const barFill = nodeGroup.append('rect')
            .attr('x', barX)
            .attr('y', barY)
            .attr('width', 0)
            .attr('height', barHeight)
            .attr('fill', capacityColor)
            .attr('rx', 2);

          barFill
            .transition()
            .duration(500)
            .attr('width', (usagePercent / 100) * barWidth);
        }
      }

      // Subtle line from center to node
      g.append('line')
        .attr('x1', 0)
        .attr('y1', 0)
        .attr('x2', x)
        .attr('y2', y)
        .attr('stroke', '#e8eaed')
        .attr('stroke-width', 1)
        .attr('opacity', 0.5);
    });

    // Draw node segments (arcs showing which part of ring belongs to each node)
    const segmentsGroup = g.append('g').attr('class', 'segments');
    
    // Group keys by their assigned node and compute segments
    const nodeSegments: Array<{ nodeId: string; startAngle: number; endAngle: number; keys: string[] }> = [];
    
    state.nodes.forEach((nodeId, nodeIdx) => {
      const nodeKeys = state.keys.filter((keyId) => state.assignments[keyId] === nodeId);
      if (nodeKeys.length === 0) return;
      
      // Find the range of key positions for this node
      const keyPositions = nodeKeys
        .map((keyId) => state.positions[keyId] || 0)
        .sort((a, b) => a - b);
      
      if (keyPositions.length > 0) {
        // For visualization, show the arc from first to last key
        // In real consistent hashing, keys between nodes belong to the next node clockwise
        const startPos = keyPositions[0];
        const endPos = keyPositions[keyPositions.length - 1];
        
        // Handle wrap-around
        let startAngle = startPos * 2 * Math.PI;
        let endAngle = endPos * 2 * Math.PI;
        
        // If segment wraps around, we need to handle it differently
        if (endPos < startPos) {
          // Wrapped segment - draw two arcs
          nodeSegments.push({
            nodeId,
            startAngle: 0,
            endAngle: endPos * 2 * Math.PI,
            keys: nodeKeys.filter((k) => state.positions[k]! <= endPos),
          });
          nodeSegments.push({
            nodeId,
            startAngle: startPos * 2 * Math.PI,
            endAngle: 2 * Math.PI,
            keys: nodeKeys.filter((k) => state.positions[k]! >= startPos),
          });
        } else {
          nodeSegments.push({ nodeId, startAngle, endAngle, keys: nodeKeys });
        }
      }
    });
    
    // Draw segment arcs
    nodeSegments.forEach((segment) => {
      const nodeIdx = state.nodes.indexOf(segment.nodeId);
      const color = nodeIdx >= 0 ? nodeColors[nodeIdx % nodeColors.length] : '#888';
      
      const arc = d3.arc<{ startAngle: number; endAngle: number }>()
        .innerRadius(keyRadius - 5)
        .outerRadius(keyRadius + 5)
        .startAngle((d) => d.startAngle)
        .endAngle((d) => d.endAngle);
      
      segmentsGroup.append('path')
        .datum(segment)
        .attr('d', arc)
        .attr('fill', color)
        .attr('opacity', 0.1)
        .attr('stroke', color)
        .attr('stroke-width', 1)
        .attr('stroke-opacity', 0.3);
    });

    // Draw connection lines (behind keys)
    const connectionsGroup = g.append('g').attr('class', 'connections');
    
    if (!useAggregated) {
      // Draw individual key connections - only show on hover or make very subtle
      // The segments above show the assignment more clearly
    } else {
      // Draw aggregated segments
      const segmentAngle = (2 * Math.PI) / 360;
      for (let i = 0; i < 360; i++) {
        const angle = i * segmentAngle;
        const x = Math.cos(angle) * keyRadius;
        const y = Math.sin(angle) * keyRadius;
        
        // Find keys in this segment
        const segmentKeys = state.keys.filter((keyId) => {
          const keyAngle = (state.positions[keyId] || 0) * 2 * Math.PI;
          const normalizedAngle = ((keyAngle % (2 * Math.PI)) + (2 * Math.PI)) % (2 * Math.PI);
          const normalizedSegment = ((angle % (2 * Math.PI)) + (2 * Math.PI)) % (2 * Math.PI);
          const diff = Math.abs(normalizedAngle - normalizedSegment);
          return diff < segmentAngle / 2 || diff > (2 * Math.PI - segmentAngle / 2);
        });

        if (segmentKeys.length > 0) {
          const assignedNode = state.assignments[segmentKeys[0]];
          if (assignedNode) {
            const nodeAngle = state.nodeAngles[assignedNode] || 0;
            const nodeX = Math.cos(nodeAngle) * radius;
            const nodeY = Math.sin(nodeAngle) * radius;
            const nodeIdx = state.nodes.indexOf(assignedNode);
            const color = nodeIdx >= 0 ? nodeColors[nodeIdx % nodeColors.length] : '#888';

            connectionsGroup.append('line')
              .attr('x1', x)
              .attr('y1', y)
              .attr('x2', nodeX)
              .attr('y2', nodeY)
              .attr('stroke', color)
              .attr('stroke-width', Math.min(segmentKeys.length / 10, 2))
              .attr('opacity', Math.min(segmentKeys.length / 50, 0.3))
              .attr('stroke-dasharray', '1,2');
          }
        }
      }
    }

    // Draw keys
    if (!useAggregated) {
      // Individual keys
      state.keys.forEach((keyId) => {
        const position = state.positions[keyId] || 0;
        const angle = position * 2 * Math.PI;
        const x = Math.cos(angle) * keyRadius;
        const y = Math.sin(angle) * keyRadius;

        const assignedNode = state.assignments[keyId];
        const nodeIdx = state.nodes.indexOf(assignedNode);
        const color = nodeIdx >= 0 ? nodeColors[nodeIdx % nodeColors.length] : '#888';

        const keyGroup = g.append('g')
          .attr('class', 'key')
          .attr('data-key-id', keyId)
          .style('cursor', 'pointer')
          .on('mouseenter', () => onKeyHover?.(keyId))
          .on('mouseleave', () => onKeyHover?.(null));

        const keyCircle = keyGroup.append('circle')
          .attr('cx', 0)
          .attr('cy', 0)
          .attr('r', 5)
          .attr('fill', color)
          .attr('opacity', 1.0)
          .attr('stroke', '#fff')
          .attr('stroke-width', 2)
          .style('filter', 'drop-shadow(0 1px 2px rgba(0,0,0,0.2))');

        keyCircle
          .transition()
          .duration(800)
          .ease(d3.easeCubicOut)
          .attr('cx', x)
          .attr('cy', y);
      });
    } else {
      // Aggregated visualization - show density as arcs
      const segmentAngle = (2 * Math.PI) / 360;
      for (let i = 0; i < 360; i++) {
        const angle = i * segmentAngle;
        const x = Math.cos(angle) * keyRadius;
        const y = Math.sin(angle) * keyRadius;
        
        const segmentKeys = state.keys.filter((keyId) => {
          const keyAngle = (state.positions[keyId] || 0) * 2 * Math.PI;
          const normalizedAngle = ((keyAngle % (2 * Math.PI)) + (2 * Math.PI)) % (2 * Math.PI);
          const normalizedSegment = ((angle % (2 * Math.PI)) + (2 * Math.PI)) % (2 * Math.PI);
          const diff = Math.abs(normalizedAngle - normalizedSegment);
          return diff < segmentAngle / 2 || diff > (2 * Math.PI - segmentAngle / 2);
        });

        if (segmentKeys.length > 0) {
          const assignedNode = state.assignments[segmentKeys[0]];
          const nodeIdx = state.nodes.indexOf(assignedNode);
          const color = nodeIdx >= 0 ? nodeColors[nodeIdx % nodeColors.length] : '#888';
          const intensity = Math.min(segmentKeys.length / 20, 1);

          g.append('circle')
            .attr('cx', x)
            .attr('cy', y)
            .attr('r', Math.max(3, Math.min(segmentKeys.length / 5, 8)))
            .attr('fill', color)
            .attr('opacity', 0.6 + intensity * 0.4)
            .attr('stroke', '#fff')
            .attr('stroke-width', 1)
            .style('filter', 'drop-shadow(0 1px 2px rgba(0,0,0,0.15))');
        }
      }
    }

  }, [state, centerX, centerY, radius, keyRadius, keysByNode, useAggregated, onNodeHover, onKeyHover]);

  return (
    <svg
      ref={svgRef}
      width={width}
      height={height}
      style={{ display: 'block' }}
    />
  );
};

export default Ring;
