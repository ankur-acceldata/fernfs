// Mock TextEncoder/TextDecoder
class MockTextEncoder {
  encode(input?: string): Uint8Array {
    return new Uint8Array(Buffer.from(input || ''));
  }
}

class MockTextDecoder {
  decode(input?: ArrayBuffer | ArrayBufferView | null): string {
    if (!input) return '';
    const buffer = input instanceof ArrayBuffer ? input : input.buffer;
    return Buffer.from(buffer).toString();
  }
}

global.TextEncoder = MockTextEncoder as any;
global.TextDecoder = MockTextDecoder as any;

// Mock WebSocket
class MockWebSocket {
  static readonly CONNECTING = 0;
  static readonly OPEN = 1;
  static readonly CLOSING = 2;
  static readonly CLOSED = 3;

  public onopen: ((event: Event) => void) | null = null;
  public onclose: ((event: Event) => void) | null = null;
  public onerror: ((event: Event) => void) | null = null;
  public onmessage: ((event: MessageEvent) => void) | null = null;
  public readyState: number = MockWebSocket.CONNECTING;

  constructor(public readonly url: string) {
    // Simulate connection
    setTimeout(() => {
      this.readyState = MockWebSocket.OPEN;
      if (this.onopen) {
        this.onopen(new Event('open'));
      }
    }, 0);
  }

  send(data: string): void {
    if (this.readyState !== MockWebSocket.OPEN) {
      if (this.onerror) {
        this.onerror(new Event('error'));
      }
      return;
    }

    // Echo the message back
    setTimeout(() => {
      if (this.onmessage) {
        this.onmessage(new MessageEvent('message', { data }));
      }
    }, 0);
  }

  close(): void {
    this.readyState = MockWebSocket.CLOSED;
    if (this.onclose) {
      this.onclose(new Event('close'));
    }
  }
}

global.WebSocket = MockWebSocket as any; 