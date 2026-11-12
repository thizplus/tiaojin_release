import type {
  WebSocketMessage,
  WebSocketConnectionStatus,
  WebSocketManagerConfig,
  SeatUpdateData,
  SeatUpdateCallback
} from '@/types';

export class WebSocketManager {
  private socket: WebSocket | null = null;
  private url: string;
  private reconnectInterval: number;
  private maxReconnectAttempts: number;
  private reconnectAttempts: number = 0;
  private reconnectTimer: NodeJS.Timeout | null = null;
  private debug: boolean;

  // Event listeners
  private onConnectionChange: ((status: WebSocketConnectionStatus) => void) | null = null;
  private onMessage: ((message: WebSocketMessage) => void) | null = null;
  private onSeatUpdate: ((data: SeatUpdateData) => void) | null = null;
  private onError: ((error: Event) => void) | null = null;

  constructor(config: WebSocketManagerConfig) {
    this.url = config.url;
    this.reconnectInterval = config.reconnectInterval || 3000;
    this.maxReconnectAttempts = config.maxReconnectAttempts || 10;
    this.debug = config.debug || false;
  }

  // Public methods
  connect(): void {
    if (this.socket?.readyState === WebSocket.OPEN) {
      this.log('WebSocket already connected');
      return;
    }

    if (this.socket?.readyState === WebSocket.CONNECTING) {
      this.log('WebSocket already connecting');
      return;
    }

    this.log(`Connecting to WebSocket: ${this.url}`);
    this.updateConnectionStatus('connecting');

    try {
      this.socket = new WebSocket(this.url);
      this.setupEventListeners();
    } catch (error) {
      this.log('Failed to create WebSocket connection:', error);
      this.updateConnectionStatus('error');
      this.scheduleReconnect();
    }
  }

  disconnect(): void {
    this.log('Manually disconnecting WebSocket');
    this.clearReconnectTimer();
    this.reconnectAttempts = 0;

    if (this.socket) {
      this.socket.onclose = null; // Prevent reconnection
      this.socket.close();
      this.socket = null;
    }

    this.updateConnectionStatus('disconnected');
  }

  reconnect(): void {
    this.log('Manual reconnection requested');
    this.disconnect();
    this.reconnectAttempts = 0;
    setTimeout(() => this.connect(), 100);
  }

  sendMessage(message: Record<string, unknown>): boolean {
    if (this.socket?.readyState === WebSocket.OPEN) {
      try {
        this.socket.send(JSON.stringify(message));
        this.log('Message sent:', message);
        return true;
      } catch (error) {
        this.log('Failed to send message:', error);
        return false;
      }
    } else {
      this.log('Cannot send message: WebSocket not connected');
      return false;
    }
  }

  // Event listener setters
  onConnectionStatusChange(callback: (status: WebSocketConnectionStatus) => void): void {
    this.onConnectionChange = callback;
  }

  onMessageReceived(callback: (message: WebSocketMessage) => void): void {
    this.onMessage = callback;
  }

  onSeatUpdateReceived(callback: SeatUpdateCallback): void {
    this.onSeatUpdate = callback;
  }

  onErrorOccurred(callback: (error: Event) => void): void {
    this.onError = callback;
  }

  // Getters
  get isConnected(): boolean {
    return this.socket?.readyState === WebSocket.OPEN;
  }

  get connectionStatus(): WebSocketConnectionStatus {
    if (!this.socket) return 'disconnected';

    switch (this.socket.readyState) {
      case WebSocket.CONNECTING:
        return 'connecting';
      case WebSocket.OPEN:
        return 'connected';
      case WebSocket.CLOSING:
      case WebSocket.CLOSED:
      default:
        return 'disconnected';
    }
  }

  // Private methods
  private setupEventListeners(): void {
    if (!this.socket) return;

    this.socket.onopen = () => {
      this.log('WebSocket connected successfully');
      this.reconnectAttempts = 0;
      this.clearReconnectTimer();
      this.updateConnectionStatus('connected');
    };

    this.socket.onmessage = (event) => {
      try {
        const message: WebSocketMessage = JSON.parse(event.data);
        this.log('Message received:', message);

        // Notify general message listener
        this.onMessage?.(message);

        // Handle specific message types
        if (message.type === 'seat_update') {
          this.onSeatUpdate?.(message.data);
        }
      } catch (error) {
        this.log('Failed to parse message:', error);
      }
    };

    this.socket.onclose = (event) => {
      this.log(`WebSocket disconnected. Code: ${event.code}, Reason: ${event.reason}`);
      this.updateConnectionStatus('disconnected');

      // Only attempt to reconnect if it wasn't a manual disconnection
      if (event.code !== 1000 && this.reconnectAttempts < this.maxReconnectAttempts) {
        this.scheduleReconnect();
      }
    };

    this.socket.onerror = (error) => {
      this.log('WebSocket error:', error);
      this.updateConnectionStatus('error');
      this.onError?.(error);
    };
  }

  private scheduleReconnect(): void {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      this.log('Max reconnection attempts reached');
      return;
    }

    this.reconnectAttempts++;
    this.updateConnectionStatus('reconnecting');

    const delay = this.reconnectInterval * Math.pow(1.5, this.reconnectAttempts - 1);
    this.log(`Scheduling reconnection attempt ${this.reconnectAttempts}/${this.maxReconnectAttempts} in ${delay}ms`);

    this.reconnectTimer = setTimeout(() => {
      this.connect();
    }, delay);
  }

  private clearReconnectTimer(): void {
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }
  }

  private updateConnectionStatus(status: WebSocketConnectionStatus): void {
    this.onConnectionChange?.(status);
  }

  private log(message: string, ...args: unknown[]): void {
    if (this.debug) {
      console.log(`[WebSocketManager] ${message}`, ...args);
    }
  }

  // Cleanup
  destroy(): void {
    this.log('Destroying WebSocket manager');
    this.clearReconnectTimer();

    if (this.socket) {
      // Remove all event listeners to prevent callbacks after destroy
      this.socket.onopen = null;
      this.socket.onmessage = null;
      this.socket.onclose = null;
      this.socket.onerror = null;

      if (this.socket.readyState === WebSocket.OPEN || this.socket.readyState === WebSocket.CONNECTING) {
        this.socket.close();
      }
      this.socket = null;
    }

    this.onConnectionChange = null;
    this.onMessage = null;
    this.onSeatUpdate = null;
    this.onError = null;
  }
}