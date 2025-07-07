import { openDB, IDBPDatabase } from 'idb';
import { OperationQueueItem } from '../types';

export class OperationQueue {
  private db: IDBPDatabase | null = null;
  private readonly DB_NAME = 'fernfs_queue';
  private readonly STORE_NAME = 'operations';
  private queueSize = 0;

  constructor() {
    this.initDB();
  }

  private async initDB(): Promise<void> {
    this.db = await openDB(this.DB_NAME, 1, {
      upgrade(db) {
        if (!db.objectStoreNames.contains('operations')) {
          db.createObjectStore('operations', { keyPath: 'id' });
        }
      },
    });
  }

  private async ensureDB(): Promise<IDBPDatabase> {
    if (!this.db) {
      await this.initDB();
    }
    if (!this.db) {
      throw new Error('Failed to initialize database');
    }
    return this.db;
  }

  /**
   * Add an operation to the queue
   * @param operation Operation to queue
   */
  async enqueue(operation: OperationQueueItem): Promise<void> {
    const db = await this.ensureDB();
    await db.put(this.STORE_NAME, operation);
    this.queueSize++;
  }

  /**
   * Get the next operation from the queue
   * @returns Next operation or null if queue is empty
   */
  async dequeue(): Promise<OperationQueueItem | null> {
    const db = await this.ensureDB();
    const operations = await db.getAll(this.STORE_NAME);

    if (operations.length === 0) {
      return null;
    }

    // Get oldest operation
    const operation = operations.sort((a, b) => a.timestamp - b.timestamp)[0];
    await db.delete(this.STORE_NAME, operation.id);
    this.queueSize = Math.max(0, this.queueSize - 1);

    return operation;
  }

  /**
   * Get all pending operations
   * @returns Array of pending operations
   */
  async getAll(): Promise<OperationQueueItem[]> {
    const db = await this.ensureDB();
    return db.getAll(this.STORE_NAME);
  }

  /**
   * Update an operation in the queue
   * @param operation Operation to update
   */
  async update(operation: OperationQueueItem): Promise<void> {
    const db = await this.ensureDB();
    await db.put(this.STORE_NAME, operation);
  }

  /**
   * Remove an operation from the queue
   * @param id Operation ID to remove
   */
  async remove(id: string): Promise<void> {
    const db = await this.ensureDB();
    await db.delete(this.STORE_NAME, id);
    this.queueSize = Math.max(0, this.queueSize - 1);
  }

  /**
   * Clear all operations from the queue
   */
  async clear(): Promise<void> {
    const db = await this.ensureDB();
    await db.clear(this.STORE_NAME);
    this.queueSize = 0;
  }

  /**
   * Get the number of pending operations
   * @returns Number of operations in the queue
   */
  size(): number {
    return this.queueSize;
  }

  async updateRetryCount(id: string): Promise<void> {
    const db = await this.ensureDB();
    const operation = await db.get(this.STORE_NAME, id);
    if (operation) {
      operation.retryCount = (operation.retryCount || 0) + 1;
      await db.put(this.STORE_NAME, operation);
    }
  }

  async close(): Promise<void> {
    if (this.db) {
      this.db.close();
      this.db = null;
    }
  }
} 