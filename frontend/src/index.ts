// Core classes
export { RemoteFS } from './RemoteFS';

// Storage adapters
export { LocalStorageAdapter } from './adapters/LocalStorageAdapter';

// Types
export type {
  DirEntry,
  FileStat,
  FileOptions,
  StorageAdapter,
  SyncOptions,
  SyncStatus,
  RemoteFSOptions,
  OperationQueueItem,
} from './types';

// Sync components
export { WebSocketClient } from './network/WebSocketClient';
export { SyncManager } from './sync/SyncManager';
export { OperationQueue } from './sync/OperationQueue';

export type {
  WebSocketMessage,
  WebSocketOptions
} from './network/WebSocketClient'; 