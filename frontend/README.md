# FernFS Browser Client

A browser-based remote file management library that provides POSIX-like file operations with offline support and real-time synchronization.

## Features

- POSIX-like file operations (mkdir, readdir, readFile, writeFile, etc.)
- Offline support with IndexedDB storage
- Real-time synchronization via WebSocket
- Multiple storage backend support through adapters
- TypeScript support with full type definitions
- Comprehensive test coverage

## Installation

```bash
npm install @fernfs/client
```

## Quick Start

```typescript
import { RemoteFS, LocalStorageAdapter } from '@fernfs/client';

// Create a storage adapter
const adapter = new LocalStorageAdapter();

// Initialize RemoteFS
const fs = new RemoteFS({
  adapter,
  enableOffline: true,
  enableWebWorker: false,
});

// Use POSIX-like file operations
async function example() {
  // Create a directory
  await fs.mkdir('/documents');

  // Write a file
  await fs.writeFile('/documents/hello.txt', 'Hello, World!');

  // Read a file
  const content = await fs.readFile('/documents/hello.txt', 'utf8');
  console.log(content); // "Hello, World!"

  // List directory contents
  const entries = await fs.readdir('/documents');
  console.log(entries); // [{ name: 'hello.txt', type: 'file', ... }]

  // Get file stats
  const stats = await fs.stat('/documents/hello.txt');
  console.log(stats.size); // 13
}
```

## API Reference

### RemoteFS

The main class that provides file system operations.

```typescript
class RemoteFS {
  constructor(options: RemoteFSOptions);

  // Directory operations
  async mkdir(path: string, mode?: number): Promise<void>;
  async rmdir(path: string): Promise<void>;
  async readdir(path: string): Promise<DirEntry[]>;

  // File operations
  async readFile(path: string, encoding?: string): Promise<ArrayBuffer | string>;
  async writeFile(path: string, data: ArrayBuffer | string, options?: FileOptions): Promise<void>;
  async unlink(path: string): Promise<void>;
  async rename(oldPath: string, newPath: string): Promise<void>;
  async chmod(path: string, mode: number): Promise<void>;
  async stat(path: string): Promise<FileStat>;

  // Synchronization
  async sync(): Promise<void>;
  async getSyncStatus(): Promise<SyncStatus>;
}
```

### Storage Adapters

Storage adapters provide the actual storage implementation. The library includes:

- `LocalStorageAdapter`: Uses IndexedDB for local storage
- Custom adapters can be implemented by following the `StorageAdapter` interface

### Types

```typescript
interface DirEntry {
  name: string;
  type: 'file' | 'directory';
  size?: number;
  modifiedAt?: Date;
  createdAt?: Date;
}

interface FileStat {
  size: number;
  mode: number;
  type: 'file' | 'directory';
  modifiedAt: Date;
  createdAt: Date;
  isFile(): boolean;
  isDirectory(): boolean;
}

interface FileOptions {
  mode?: number;
  encoding?: string;
  create?: boolean;
  truncate?: boolean;
}

interface SyncStatus {
  isSyncing: boolean;
  lastSyncTime?: Date;
  pendingOperations: number;
  errors?: Error[];
}
```

## Development

### Setup

```bash
# Install dependencies
npm install

# Run tests
npm test

# Build library
npm run build

# Run linter
npm run lint

# Format code
npm run format
```

### Testing

The library uses Jest for testing. Run tests with:

```bash
# Run all tests
npm test

# Run tests in watch mode
npm run test:watch

# Run tests with coverage
npm run test:coverage
```

### Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT License - see the [LICENSE](LICENSE) file for details 