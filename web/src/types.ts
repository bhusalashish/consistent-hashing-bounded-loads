export interface VisualizerState {
  nodes: string[];
  keys: string[];
  positions: Record<string, number>; // key/node → position (0..1)
  assignments: Record<string, string>; // key → node
  nodeAngles: Record<string, number>; // node → angle in radians
  algorithm: string;
  stats?: Statistics;
  chblConfig?: CHBLConfig;
}

export interface CHBLConfig {
  loadFactor: number;
  expectedKeys: number;
  capacityPerNode: number;
}

export interface AlgorithmComparison {
  algorithm: string;
  state: VisualizerState;
  stats: Statistics;
}

export interface Statistics {
  operation: string;
  totalKeys: number;
  keysMoved: number;
  keysMovedPercent: number;
  movementByNode: Record<string, number>;
  distribution: Record<string, number>;
  previousDist: Record<string, number>;
  keyMovements: KeyMovement[];
  capacityInfo?: CapacityInfo;
}

export interface CapacityInfo {
  nodesAtCapacity: string[];
  unassignedKeys: number;
  capacityPerNode: Record<string, number>;
  currentLoad: Record<string, number>;
  loadPercentage: Record<string, number>;
}

export interface KeyMovement {
  keyId: string;
  fromNode: string;
  toNode: string;
}

export type Algorithm = 'ring' | 'jump' | 'maglev' | 'chbl';
