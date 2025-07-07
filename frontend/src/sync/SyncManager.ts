import { WebSocketClient } from '../network/WebSocketClient';
import { OperationQueue } from './OperationQueue';
import { StorageAdapter, WebSocketMessage, SyncStatus, SyncOptions } from '../types/index';

export class SyncManager {
  private isSyncing: boolean = false;
  private lastSyncTime: Date | null = null;
  private syncInterval: number | null = null;

  constructor(
    private adapter: StorageAdapter,
    private operationQueue: OperationQueue,
    private wsClient: WebSocketClient,
    private options: SyncOptions = {}
  ) {
    this.setupWebSocket();
    if (options.autoSync) {
      this.startAutoSync(options.syncInterval || 30000);
    }
  }

  private setupWebSocket(): void {
    this.wsClient.onMessage((message: WebSocketMessage) => {
      if (message.type === 'response') {
        // Handle successful response
        this.handleSyncResponse(message);
      } else if (message.type === 'error') {
        // Handle error response
        this.handleSyncError(message);
      }
    });
  }

  private async syncOperation(operation: OperationQueueItem): Promise<void> {
    const message: WebSocketMessage = {
      id: operation.id,
      type: 'request',
      operation: operation.operation,
      path: operation.path,
      data: operation.data,
      options: operation.options,
      timestamp: Date.now()
    };

    await this.wsClient.send(message);
  }

  private handleSyncResponse(message: WebSocketMessage): void {
    // Handle successful sync response
    this.operationQueue.remove(message.id);
    this.lastSyncTime = new Date();
  }

  private handleSyncError(message: WebSocketMessage): void {
    // Handle sync error
    console.error('Sync error:', message.error);
  }

  private startAutoSync(interval: number): void {
    this.syncInterval = window.setInterval(() => {
      this.sync().catch(console.error);
    }, interval);
  }

  async sync(): Promise<void> {
    if (this.isSyncing) {
      return;
    }

    this.isSyncing = true;
    try {
      const operations = await this.operationQueue.getAll();
      for (const operation of operations) {
        await this.syncOperation(operation);
      }
      this.lastSyncTime = new Date();
    } finally {
      this.isSyncing = false;
    }
  }

  getSyncStatus(): SyncStatus {
    return {
      isSyncing: this.isSyncing,
      lastSync: this.lastSyncTime,
      pendingOperations: this.operationQueue.size()
    };
  }

  async close(): Promise<void> {
    if (this.syncInterval !== null) {
      window.clearInterval(this.syncInterval);
      this.syncInterval = null;
    }
  }
}

interface OperationQueueItem {
  id: string;
  operation: string;
  path: string;
  data?: any;
  options?: Record<string, unknown>;
} 