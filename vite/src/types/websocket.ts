// WebSocket Types
export interface SeatUpdateData {
  period_id: string;
  total_seats: number;
  booked_seats: number;
  held_seats: number;
  available_seats: number;
  update_type: SeatUpdateType;
  timestamp: string;
}

export type SeatUpdateType =
  | 'seat_held'
  | 'seat_released'
  | 'seat_hold_expired'
  | 'booking_created'
  | 'booking_updated'
  | 'booking_cancelled'
  | 'booking_expired';

export interface WebSocketMessage {
  type: 'seat_update';
  data: SeatUpdateData;
}

export interface WebSocketContextType {
  isConnected: boolean;
  connectionStatus: WebSocketConnectionStatus;
  lastMessage: WebSocketMessage | null;
  sendMessage: (message: any) => void;
  reconnect: () => void;
  disconnect: () => void;
}

export type WebSocketConnectionStatus =
  | 'connecting'
  | 'connected'
  | 'disconnected'
  | 'reconnecting'
  | 'error';

export interface SeatUpdateCallback {
  (data: SeatUpdateData): void;
}

export interface NotificationConfig {
  type: 'info' | 'success' | 'warning' | 'error';
  title: string;
  description?: string;
  duration?: number;
}

export interface SeatDisplayData {
  available: number;
  booked: number;
  held: number;
  total: number;
}

export interface WebSocketManagerConfig {
  url: string;
  reconnectInterval?: number;
  maxReconnectAttempts?: number;
  debug?: boolean;
}