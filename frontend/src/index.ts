// Core components
export { RemoteFS } from './RemoteFS';
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
  WebSocketMessage
} from './types/index';

// Sync components
export { WebSocketClient } from './network/WebSocketClient';
export { SyncManager } from './sync/SyncManager';
export { OperationQueue } from './sync/OperationQueue';

export type {
  WebSocketOptions
} from './network/WebSocketClient'; 