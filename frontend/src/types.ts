export * from './types/index';

export type FileMode = number;

export interface DirEntry {
  name: string;
  type: 'file' | 'directory';
  mode: FileMode;
}

export interface FileStat {
  size: number;
  mode: FileMode;
  isDirectory: boolean;
  isFile: boolean;
  modifiedTime: number;
  createdTime: number;
}

export interface FileOperationBase {
  type: string;
  path: string;
  timestamp: number;
}

export interface MkdirOperation extends FileOperationBase {
  type: 'mkdir';
  mode: FileMode;
}

export interface RmdirOperation extends FileOperationBase {
  type: 'rmdir';
}

export interface WriteFileOperation extends FileOperationBase {
  type: 'writeFile';
  data: ArrayBuffer;
  mode?: FileMode;
}

export interface UnlinkOperation extends FileOperationBase {
  type: 'unlink';
}

export interface RenameOperation extends FileOperationBase {
  type: 'rename';
  oldPath: string;
  newPath: string;
}

export interface ChmodOperation extends FileOperationBase {
  type: 'chmod';
  mode: FileMode;
}

export type FileOperation = 
  | MkdirOperation 
  | RmdirOperation 
  | WriteFileOperation 
  | UnlinkOperation 
  | RenameOperation 
  | ChmodOperation;

export interface WebSocketMessage {
  id: string;
  type: 'request' | 'response' | 'error';
  operation: string;
  data?: any;
  error?: string;
} 