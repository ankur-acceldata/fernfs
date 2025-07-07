import { WebSocket, WebSocketServer } from 'ws';
import { WebSocketMessage, DirEntry, FileStat } from '../types';

export class MockWebSocketServer {
  private clients: Set<WebSocket> = new Set();
  private files: Map<string, ArrayBuffer> = new Map();
  private metadata: Map<string, FileStat> = new Map();
  private directories: Map<string, DirEntry[]> = new Map();
  private server: WebSocketServer;

  constructor() {
    // Create WebSocket server
    this.server = new WebSocketServer({ port: 8080 });

    this.server.on('connection', (ws) => {
      this.clients.add(ws);

      ws.on('message', (data) => {
        this.handleMessage(ws, data.toString());
      });

      ws.on('close', () => {
        this.clients.delete(ws);
      });
    });
  }

  private handleMessage(ws: WebSocket, data: string): void {
    try {
      const message: WebSocketMessage = JSON.parse(data);
      const response = this.processRequest(message);
      ws.send(JSON.stringify(response));
    } catch (err) {
      const error = err as Error;
      ws.send(JSON.stringify({
        id: 'error',
        type: 'error',
        operation: 'unknown',
        error: error.message
      }));
    }
  }

  private processRequest(message: WebSocketMessage): WebSocketMessage {
    const { id, operation, data } = message;
    let response: any;

    switch (operation) {
      case 'mkdir':
        this.metadata.set(data.path, {
          size: 0,
          mode: data.mode || 0o755,
          isFile: false,
          isDirectory: true,
          modifiedTime: Date.now(),
          createdTime: Date.now()
        });
        this.directories.set(data.path, []);
        response = null;
        break;

      case 'rmdir':
        this.metadata.delete(data.path);
        this.directories.delete(data.path);
        response = null;
        break;

      case 'readdir':
        response = this.directories.get(data.path) || [];
        break;

      case 'stat':
        const stat = this.metadata.get(data.path);
        if (!stat) {
          throw new Error(`File not found: ${data.path}`);
        }
        response = stat;
        break;

      case 'readFile':
        const file = this.files.get(data.path);
        if (!file) {
          throw new Error(`File not found: ${data.path}`);
        }
        response = file;
        break;

      case 'writeFile':
        this.files.set(data.path, data.data);
        this.metadata.set(data.path, {
          size: data.data.byteLength,
          mode: data.mode || 0o644,
          isFile: true,
          isDirectory: false,
          modifiedTime: Date.now(),
          createdTime: Date.now()
        });
        response = null;
        break;

      case 'unlink':
        this.files.delete(data.path);
        this.metadata.delete(data.path);
        response = null;
        break;

      case 'rename':
        const fileData = this.files.get(data.oldPath);
        const fileMeta = this.metadata.get(data.oldPath);
        if (fileData && fileMeta) {
          this.files.set(data.newPath, fileData);
          this.metadata.set(data.newPath, { ...fileMeta });
          this.files.delete(data.oldPath);
          this.metadata.delete(data.oldPath);
        }
        response = null;
        break;

      case 'chmod':
        const meta = this.metadata.get(data.path);
        if (meta) {
          this.metadata.set(data.path, { ...meta, mode: data.mode });
        }
        response = null;
        break;

      default:
        throw new Error(`Unknown operation: ${operation}`);
    }

    return {
      id,
      type: 'response',
      operation,
      data: response
    };
  }

  close(): void {
    this.server.close();
    for (const client of this.clients) {
      client.close();
    }
    this.clients.clear();
    this.files.clear();
    this.metadata.clear();
    this.directories.clear();
  }
} 