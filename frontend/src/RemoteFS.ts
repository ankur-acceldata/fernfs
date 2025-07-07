import { WebSocketClient } from './network/WebSocketClient';
import { SyncManager } from './sync/SyncManager';
import { OperationQueue } from './sync/OperationQueue';
import {
  StorageAdapter,
  DirEntry,
  FileStat,
  FileOptions,
  RemoteFSOptions,
  SyncStatus
} from './types';

export class RemoteFS {
  private adapter: StorageAdapter;
  private syncManager: SyncManager | null = null;
  private operationQueue: OperationQueue | null = null;
  private webSocketClient: WebSocketClient | null = null;

  constructor(options: RemoteFSOptions) {
    this.adapter = options.adapter;

    if (options.enableOffline) {
      this.operationQueue = new OperationQueue();
      this.webSocketClient = new WebSocketClient(options.serverUrl || 'ws://localhost:8080');
      this.syncManager = new SyncManager(
        this.adapter,
        this.operationQueue,
        this.webSocketClient,
        options.sync
      );
    }
  }

  async init(): Promise<void> {
    if (this.syncManager) {
      await this.syncManager.init();
    }
  }

  async close(): Promise<void> {
    if (this.webSocketClient) {
      await this.webSocketClient.close();
    }
    if (this.operationQueue) {
      await this.operationQueue.close();
    }
  }

  async mkdir(path: string, mode: number = 0o755): Promise<void> {
    await this.adapter.mkdir(path, mode);
  }

  async rmdir(path: string): Promise<void> {
    await this.adapter.rmdir(path);
  }

  async readdir(path: string): Promise<DirEntry[]> {
    return this.adapter.readdir(path);
  }

  async stat(path: string): Promise<FileStat> {
    return this.adapter.stat(path);
  }

  async readFile(path: string, encoding?: string): Promise<ArrayBuffer | string> {
    return this.adapter.readFile(path, encoding);
  }

  async writeFile(path: string, data: ArrayBuffer | string, options: FileOptions = {}): Promise<void> {
    await this.adapter.writeFile(path, data, options);
  }

  async unlink(path: string): Promise<void> {
    await this.adapter.unlink(path);
  }

  async rename(oldPath: string, newPath: string): Promise<void> {
    await this.adapter.rename(oldPath, newPath);
  }

  async chmod(path: string, mode: number): Promise<void> {
    await this.adapter.chmod(path, mode);
  }

  async sync(): Promise<void> {
    if (this.syncManager) {
      await this.syncManager.sync();
    }
  }

  async getSyncStatus(): Promise<SyncStatus> {
    if (!this.syncManager) {
      return {
        isSyncing: false,
        pendingOperations: 0
      };
    }
    return this.syncManager.getStatus();
  }
} 