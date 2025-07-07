export interface DirEntry {
  name: string;
  type: 'file' | 'directory';
  size: number;
  mode: number;
  modifiedTime: number;
  createdTime: number;
}

export interface FileStat {
  size: number;
  mode: number;
  isFile: boolean;
  isDirectory: boolean;
  modifiedTime: number;
  createdTime: number;
}

export interface FileOptions {
  mode?: number;
  encoding?: string;
  create?: boolean;
  truncate?: boolean;
}

export interface StorageAdapter {
  mkdir(path: string, mode?: number): Promise<void>;
  rmdir(path: string): Promise<void>;
  readdir(path: string): Promise<DirEntry[]>;
  stat(path: string): Promise<FileStat>;
  readFile(path: string, encoding?: string): Promise<ArrayBuffer | string>;
  writeFile(path: string, data: ArrayBuffer | string, options?: FileOptions): Promise<void>;
  unlink(path: string): Promise<void>;
  rename(oldPath: string, newPath: string): Promise<void>;
  chmod(path: string, mode: number): Promise<void>;
  
  // Cache management (optional)
  cacheFile?(path: string, data: ArrayBuffer): Promise<void>;
  cacheDirectoryListing?(path: string, entries: DirEntry[]): Promise<void>;
  clearCache?(): Promise<void>;
}

export interface SyncOptions {
  retryAttempts?: number;
  retryDelay?: number;
}

export interface RemoteFSOptions {
  adapter: StorageAdapter;
  enableOffline?: boolean;
  serverUrl?: string;
  sync?: SyncOptions;
}

export interface SyncStatus {
  isSyncing: boolean;
  pendingOperations: number;
}

export interface OperationQueueItem {
  id: string;
  operation: string;
  path: string;
  data?: ArrayBuffer | string;
  options?: Record<string, unknown>;
  timestamp: number;
  retryCount?: number;
}

export interface WebSocketMessage {
  id: string;
  type: 'request' | 'response' | 'error';
  operation: string;
  path?: string;
  data?: any;
  options?: Record<string, unknown>;
  error?: string;
  timestamp?: number;
} 