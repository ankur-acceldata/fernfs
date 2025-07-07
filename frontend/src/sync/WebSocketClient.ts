export interface WebSocketMessage {
  id: string;
  operation: string;
  path: string;
  data?: string;
  options?: Record<string, unknown>;
  timestamp: number;
}

export interface WebSocketOptions {
  reconnectAttempts?: number;
  reconnectDelay?: number;
  heartbeatInterval?: number;
}

export class WebSocketClient {
  private ws: WebSocket | null = null;
  private reconnectAttempts: number;
  private reconnectDelay: number;
  private heartbeatInterval: number;
  private messageHandlers: Map<string, (response: WebSocketMessage) => void>;
  private connectionPromise: Promise<void> | null = null;

  constructor(
    private serverUrl: string,
    options: WebSocketOptions = {}
  ) {
    this.reconnectAttempts = options.reconnectAttempts ?? 5;
    this.reconnectDelay = options.reconnectDelay ?? 1000;
    this.heartbeatInterval = options.heartbeatInterval ?? 30000;
    this.messageHandlers = new Map();

    // Bind methods
    this.handleMessage = this.handleMessage.bind(this);
    this.handleClose = this.handleClose.bind(this);
    this.handleError = this.handleError.bind(this);
  }

  private async connect(): Promise<void> {
    if (this.ws?.readyState === WebSocket.OPEN) {
      return;
    }

    if (this.connectionPromise) {
      return this.connectionPromise;
    }

    this.connectionPromise = new Promise((resolve, reject) => {
      try {
        this.ws = new WebSocket(this.serverUrl);

        this.ws.onopen = () => {
          this.startHeartbeat();
          resolve();
        };

        this.ws.onmessage = this.handleMessage;
        this.ws.onclose = this.handleClose;
        this.ws.onerror = this.handleError;
      } catch (error) {
        reject(error);
      }
    });

    return this.connectionPromise;
  }

  private handleMessage(event: MessageEvent): void {
    try {
      const message: WebSocketMessage = JSON.parse(event.data);
      const handler = this.messageHandlers.get(message.id);
      if (handler) {
        handler(message);
        this.messageHandlers.delete(message.id);
      }
    } catch (error) {
      console.error('Error handling WebSocket message:', error);
    }
  }

  private handleClose(event: CloseEvent): void {
    console.log('WebSocket connection closed:', event.code, event.reason);
    this.ws = null;
    this.connectionPromise = null;
    this.reconnect();
  }

  private handleError(event: Event): void {
    console.error('WebSocket error:', event);
  }

  private async reconnect(attempt = 0): Promise<void> {
    if (attempt >= this.reconnectAttempts) {
      console.error('Max reconnection attempts reached');
      return;
    }

    const delay = this.reconnectDelay * Math.pow(2, attempt);
    await new Promise(resolve => setTimeout(resolve, delay));

    try {
      await this.connect();
    } catch (error) {
      console.error('Reconnection failed:', error);
      this.reconnect(attempt + 1);
    }
  }

  private startHeartbeat(): void {
    setInterval(() => {
      if (this.ws?.readyState === WebSocket.OPEN) {
        this.ws.send(JSON.stringify({ type: 'ping' }));
      }
    }, this.heartbeatInterval);
  }

  async send(message: WebSocketMessage): Promise<WebSocketMessage> {
    await this.connect();

    return new Promise((resolve, reject) => {
      if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
        reject(new Error('WebSocket is not connected'));
        return;
      }

      this.messageHandlers.set(message.id, resolve);
      this.ws.send(JSON.stringify(message));

      // Set timeout for response
      setTimeout(() => {
        this.messageHandlers.delete(message.id);
        reject(new Error('WebSocket request timeout'));
      }, 30000);
    });
  }

  async close(): Promise<void> {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
      this.connectionPromise = null;
    }
  }
} 