import { RemoteFS } from '../RemoteFS';
import { LocalStorageAdapter } from '../adapters/LocalStorageAdapter';
import { StorageAdapter } from '../types/index';
import 'fake-indexeddb/auto';

describe('RemoteFS', () => {
  let fs: RemoteFS;
  let adapter: StorageAdapter;

  beforeEach(() => {
    adapter = new LocalStorageAdapter();
    fs = new RemoteFS({ adapter });
  });

  afterEach(async () => {
    await fs.close();
  });

  describe('Directory Operations', () => {
    it('should create and list directories', async () => {
      await fs.mkdir('/test');
      const entries = await fs.readdir('/');
      expect(entries).toHaveLength(1);
      expect(entries[0].name).toBe('test');
      expect(entries[0].type).toBe('directory');
    });

    it('should remove empty directories', async () => {
      await fs.mkdir('/test');
      await fs.rmdir('/test');
      const entries = await fs.readdir('/');
      expect(entries).toHaveLength(0);
    });

    it('should fail to remove non-empty directories', async () => {
      await fs.mkdir('/test');
      await fs.writeFile('/test/file.txt', 'content');
      await expect(fs.rmdir('/test')).rejects.toThrow();
    });
  });

  describe('File Operations', () => {
    it('should write and read files', async () => {
      await fs.writeFile('/test.txt', 'content');
      const content = await fs.readFile('/test.txt', 'utf-8');
      expect(content).toBe('content');
    });

    it('should handle binary data', async () => {
      const data = new Uint8Array([1, 2, 3, 4]);
      await fs.writeFile('/test.bin', data.buffer);
      const content = await fs.readFile('/test.bin');
      expect(content).toBeInstanceOf(ArrayBuffer);
      expect(new Uint8Array(content as ArrayBuffer)).toEqual(data);
    });

    it('should delete files', async () => {
      await fs.writeFile('/test.txt', 'content');
      await fs.unlink('/test.txt');
      await expect(fs.readFile('/test.txt')).rejects.toThrow();
    });

    it('should rename files', async () => {
      await fs.writeFile('/test.txt', 'content');
      await fs.rename('/test.txt', '/renamed.txt');
      const content = await fs.readFile('/renamed.txt', 'utf-8');
      expect(content).toBe('content');
      await expect(fs.readFile('/test.txt')).rejects.toThrow();
    });
  });

  describe('File Stats', () => {
    it('should get file stats', async () => {
      await fs.writeFile('/test.txt', 'content');
      const stats = await fs.stat('/test.txt');
      expect(stats.isFile).toBe(true);
      expect(stats.isDirectory).toBe(false);
      expect(stats.size).toBe(7);
    });

    it('should get directory stats', async () => {
      await fs.mkdir('/test');
      const stats = await fs.stat('/test');
      expect(stats.isFile).toBe(false);
      expect(stats.isDirectory).toBe(true);
    });

    it('should update file permissions', async () => {
      await fs.writeFile('/test.txt', 'content');
      await fs.chmod('/test.txt', 0o644);
      const stats = await fs.stat('/test.txt');
      expect(stats.mode).toBe(0o644);
    });
  });

  describe('Error Handling', () => {
    it('should handle non-existent files', async () => {
      await expect(fs.readFile('/nonexistent.txt')).rejects.toThrow();
    });

    it('should handle invalid operations', async () => {
      await fs.mkdir('/test');
      await expect(fs.mkdir('/test')).rejects.toThrow();
    });
  });

  describe('Offline Support', () => {
    let onlineFs: RemoteFS;
    let offlineFs: RemoteFS;

    beforeEach(() => {
      const adapter = new LocalStorageAdapter();
      onlineFs = new RemoteFS({
        adapter,
        enableOffline: true,
        serverUrl: 'ws://localhost:8080'
      });

      offlineFs = new RemoteFS({
        adapter,
        enableOffline: true,
        serverUrl: 'ws://localhost:8080'
      });
    });

    afterEach(async () => {
      await onlineFs.close();
      await offlineFs.close();
    });

    it('should sync changes when online', async () => {
      await onlineFs.writeFile('/test.txt', 'content');
      await onlineFs.sync();
      const status = onlineFs.getSyncStatus();
      expect(status.isSyncing).toBe(false);
      expect(status.pendingOperations).toBe(0);
    });
  });
}); 