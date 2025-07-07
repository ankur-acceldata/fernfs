import { DirEntry, FileStat, FileMode } from '../types';

export interface StorageAdapter {
  // Core operations
  mkdir(path: string, mode: FileMode): Promise<void>;
  rmdir(path: string): Promise<void>;
  readdir(path: string): Promise<DirEntry[]>;
  stat(path: string): Promise<FileStat>;
  readFile(path: string): Promise<ArrayBuffer>;
  writeFile(path: string, data: ArrayBuffer, mode?: FileMode): Promise<void>;
  unlink(path: string): Promise<void>;
  rename(oldPath: string, newPath: string): Promise<void>;
  chmod(path: string, mode: FileMode): Promise<void>;

  // Cache management
  cacheFile(path: string, data: ArrayBuffer): Promise<void>;
  cacheDirectoryListing(path: string, entries: DirEntry[]): Promise<void>;
  clearCache(): Promise<void>;
} 