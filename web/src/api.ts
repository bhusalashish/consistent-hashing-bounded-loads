import { VisualizerState } from './types';

const API_BASE = 'http://localhost:8080';

export async function fetchState(): Promise<VisualizerState> {
  const response = await fetch(`${API_BASE}/state`);
  if (!response.ok) {
    throw new Error(`Failed to fetch state: ${response.statusText}`);
  }
  const data = await response.json();
  return data.state;
}

export async function addNode(): Promise<VisualizerState> {
  const response = await fetch(`${API_BASE}/add-node`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({}),
  });
  if (!response.ok) {
    throw new Error(`Failed to add node: ${response.statusText}`);
  }
  const data = await response.json();
  return data.state;
}

export async function removeNode(nodeId: string): Promise<VisualizerState> {
  const response = await fetch(`${API_BASE}/remove-node`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ nodeId }),
  });
  if (!response.ok) {
    throw new Error(`Failed to remove node: ${response.statusText}`);
  }
  const data = await response.json();
  return data.state;
}

export async function regenerateKeys(count: number): Promise<VisualizerState> {
  const response = await fetch(`${API_BASE}/regenerate-keys`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ count }),
  });
  if (!response.ok) {
    throw new Error(`Failed to regenerate keys: ${response.statusText}`);
  }
  const data = await response.json();
  return data.state;
}

export async function setAlgorithm(algorithm: string): Promise<VisualizerState> {
  const response = await fetch(`${API_BASE}/set-algorithm`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ algorithm }),
  });
  if (!response.ok) {
    throw new Error(`Failed to set algorithm: ${response.statusText}`);
  }
  const data = await response.json();
  return data.state;
}

export async function setKeyCount(count: number): Promise<VisualizerState> {
  const response = await fetch(`${API_BASE}/set-key-count?count=${count}`, {
    method: 'POST',
  });
  if (!response.ok) {
    throw new Error(`Failed to set key count: ${response.statusText}`);
  }
  const data = await response.json();
  return data.state;
}

export async function setCHBLConfig(loadFactor: number, expectedKeys: number): Promise<VisualizerState> {
  const response = await fetch(`${API_BASE}/set-chbl-config`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ loadFactor, expectedKeys }),
  });
  if (!response.ok) {
    throw new Error(`Failed to set CH-BL config: ${response.statusText}`);
  }
  const data = await response.json();
  return data.state;
}

export interface AlgorithmComparison {
  algorithm: string;
  state: VisualizerState;
  stats: Statistics;
}

export async function compareOperation(operation: string, nodeId?: string): Promise<AlgorithmComparison[]> {
  const response = await fetch(`${API_BASE}/compare-operation`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ operation, nodeId }),
  });
  if (!response.ok) {
    throw new Error(`Failed to compare operation: ${response.statusText}`);
  }
  const data = await response.json();
  return data.comparison;
}

