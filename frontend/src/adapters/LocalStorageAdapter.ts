import { openDB, IDBPDatabase } from 'idb';
import { StorageAdapter, DirEntry, FileStat, FileOptions } from '../types/index';

interface FileEntry {
  path: string;
  content: ArrayBuffer | string;
  type: 'file' | 'directory';
  mode: number;
  size: number;
  modifiedAt: Date;
  createdAt: Date;
}

class FileStatImpl implements FileStat {
  public readonly isFile: boolean;
  public readonly isDirectory: boolean;
  public readonly size: number;
  public readonly mode: number;
  public readonly modifiedTime: number;
  public readonly createdTime: number;

  constructor(
    size: number,
    mode: number,
    type: 'file' | 'directory',
    modifiedAt: Date,
    createdAt: Date
  ) {
    this.size = size;
    this.mode = mode;
    this.isFile = type === 'file';
    this.isDirectory = type === 'directory';
    this.modifiedTime = modifiedAt.getTime();
    this.createdTime = createdAt.getTime();
  }
}

export class LocalStorageAdapter implements StorageAdapter {
  private db: IDBPDatabase | null = null;
  private readonly DB_NAME = 'fernfs_storage';
  private readonly STORE_NAME = 'files';

  constructor() {
    this.initDB();
  }

  private async initDB(): Promise<void> {
    this.db = await openDB(this.DB_NAME, 1, {
      upgrade(db) {
        if (!db.objectStoreNames.contains('files')) {
          db.createObjectStore('files', { keyPath: 'path' });
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

  private async getEntry(path: string): Promise<FileEntry | undefined> {
    const db = await this.ensureDB();
    return db.get(this.STORE_NAME, path);
  }

  private async putEntry(entry: FileEntry): Promise<void> {
    const db = await this.ensureDB();
    await db.put(this.STORE_NAME, entry);
  }

  private async deleteEntry(path: string): Promise<void> {
    const db = await this.ensureDB();
    await db.delete(this.STORE_NAME, path);
  }

  private async getAllEntries(): Promise<FileEntry[]> {
    const db = await this.ensureDB();
    return db.getAll(this.STORE_NAME);
  }

  private normalizePath(path: string): string {
    return path.startsWith('/') ? path : `/${path}`;
  }

  private getParentPath(path: string): string {
    const normalized = this.normalizePath(path);
    const lastSlash = normalized.lastIndexOf('/');
    return lastSlash === 0 ? '/' : normalized.slice(0, lastSlash);
  }

  private getBaseName(path: string): string {
    const normalized = this.normalizePath(path);
    const lastSlash = normalized.lastIndexOf('/');
    return normalized.slice(lastSlash + 1);
  }

  async mkdir(path: string, mode: number = 0o755): Promise<void> {
    const normalized = this.normalizePath(path);
    const existing = await this.getEntry(normalized);
    if (existing) {
      throw new Error(`Path already exists: ${normalized}`);
    }

    const entry: FileEntry = {
      path: normalized,
      content: '',
      type: 'directory',
      mode,
      size: 0,
      modifiedAt: new Date(),
      createdAt: new Date()
    };

    await this.putEntry(entry);
  }

  async rmdir(path: string): Promise<void> {
    const normalized = this.normalizePath(path);
    const entry = await this.getEntry(normalized);
    if (!entry) {
      throw new Error(`Directory not found: ${normalized}`);
    }
    if (entry.type !== 'directory') {
      throw new Error(`Not a directory: ${normalized}`);
    }

    const children = await this.readdir(normalized);
    if (children.length > 0) {
      throw new Error(`Directory not empty: ${normalized}`);
    }

    await this.deleteEntry(normalized);
  }

  async readdir(path: string): Promise<DirEntry[]> {
    const normalized = this.normalizePath(path);
    const entries = await this.getAllEntries();
    const parent = normalized === '/' ? '' : normalized;
    
    return entries
      .filter(entry => {
        const entryParent = this.getParentPath(entry.path);
        return entryParent === parent;
      })
      .map(entry => ({
        name: this.getBaseName(entry.path),
        type: entry.type,
        size: entry.size,
        mode: entry.mode,
        modifiedTime: entry.modifiedAt.getTime(),
        createdTime: entry.createdAt.getTime()
      }));
  }

  async stat(path: string): Promise<FileStat> {
    const normalized = this.normalizePath(path);
    const entry = await this.getEntry(normalized);
    if (!entry) {
      throw new Error(`Path not found: ${normalized}`);
    }

    return new FileStatImpl(
      entry.size,
      entry.mode,
      entry.type,
      entry.modifiedAt,
      entry.createdAt
    );
  }

  async readFile(path: string, encoding?: string): Promise<ArrayBuffer | string> {
    const normalized = this.normalizePath(path);
    const entry = await this.getEntry(normalized);
    if (!entry) {
      throw new Error(`File not found: ${normalized}`);
    }
    if (entry.type !== 'file') {
      throw new Error(`Not a file: ${normalized}`);
    }

    return entry.content;
  }

  async writeFile(path: string, data: ArrayBuffer | string, options: FileOptions = {}): Promise<void> {
    const normalized = this.normalizePath(path);
    const existing = await this.getEntry(normalized);
    if (existing && existing.type === 'directory') {
      throw new Error(`Cannot write to directory: ${normalized}`);
    }

    const entry: FileEntry = {
      path: normalized,
      content: data,
      type: 'file',
      mode: options.mode ?? 0o644,
      size: data instanceof ArrayBuffer ? data.byteLength : data.length,
      modifiedAt: new Date(),
      createdAt: existing?.createdAt ?? new Date()
    };

    await this.putEntry(entry);
  }

  async unlink(path: string): Promise<void> {
    const normalized = this.normalizePath(path);
    const entry = await this.getEntry(normalized);
    if (!entry) {
      throw new Error(`File not found: ${normalized}`);
    }
    if (entry.type !== 'file') {
      throw new Error(`Not a file: ${normalized}`);
    }

    await this.deleteEntry(normalized);
  }

  async rename(oldPath: string, newPath: string): Promise<void> {
    const normalizedOld = this.normalizePath(oldPath);
    const normalizedNew = this.normalizePath(newPath);
    const entry = await this.getEntry(normalizedOld);
    if (!entry) {
      throw new Error(`Path not found: ${normalizedOld}`);
    }

    const newEntry = { ...entry, path: normalizedNew };
    await this.putEntry(newEntry);
    await this.deleteEntry(normalizedOld);
  }

  async chmod(path: string, mode: number): Promise<void> {
    const normalized = this.normalizePath(path);
    const entry = await this.getEntry(normalized);
    if (!entry) {
      throw new Error(`Path not found: ${normalized}`);
    }

    const updatedEntry = { ...entry, mode };
    await this.putEntry(updatedEntry);
  }
} 