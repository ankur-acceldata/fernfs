import { WebSocketClient } from '../network/WebSocketClient';
import { OperationQueue } from './OperationQueue';
import { StorageAdapter, SyncOptions, SyncStatus } from '../types';

export class SyncManager {
  private isSyncing = false;
  private lastSyncTime: Date | undefined;
  private errors: Error[] = [];
  private retryAttempts: number;
  private retryDelay: number;

  constructor(
    private adapter: StorageAdapter,
    private queue: OperationQueue,
    private wsClient: WebSocketClient,
    options: SyncOptions = {}
  ) {
    this.retryAttempts = options.retryAttempts ?? 3;
    this.retryDelay = options.retryDelay ?? 1000;
  }

  async init(): Promise<void> {
    await this.wsClient.connect();
    this.setupEventListeners();
  }

  private setupEventListeners(): void {
    this.wsClient.onMessage((message) => {
      // Handle incoming sync messages
      console.log('Received sync message:', message);
    });

    this.wsClient.onError((error) => {
      console.error('WebSocket error:', error);
      this.errors.push(error);
    });
  }

  async sync(): Promise<void> {
    if (this.isSyncing) {
      return;
    }

    this.isSyncing = true;
    try {
      const operations = await this.queue.getAll();
      for (const operation of operations) {
        try {
          await this.wsClient.send(operation);
          await this.queue.remove(operation.id);
        } catch (error) {
          console.error('Failed to sync operation:', error);
          this.errors.push(error as Error);
          if (operation.retryCount && operation.retryCount >= this.retryAttempts) {
            await this.queue.remove(operation.id);
          } else {
            await this.queue.updateRetryCount(operation.id);
          }
        }
      }
      this.lastSyncTime = new Date();
    } finally {
      this.isSyncing = false;
    }
  }

  getStatus(): SyncStatus {
    return {
      isSyncing: this.isSyncing,
      lastSyncTime: this.lastSyncTime,
      pendingOperations: this.queue.size(),
      errors: this.errors.length > 0 ? [...this.errors] : undefined,
    };
  }
} 