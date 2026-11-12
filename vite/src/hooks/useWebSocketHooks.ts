import { useContext, useState, useRef, useEffect } from 'react';
import { createContext } from 'react';
import type {
  WebSocketContextType,
  SeatUpdateData
} from '@/types';

// Context definition (for type safety)
export const WebSocketContext = createContext<WebSocketContextType | null>(null);

export function useWebSocket() {
  const context = useContext(WebSocketContext);
  if (!context) {
    throw new Error('useWebSocket must be used within a WebSocketProvider');
  }
  return context;
}

// Specialized hook for seat updates
export function useSeatUpdates(periodId?: string) {
  const { subscribeSeatUpdates, getCachedSeatData } = useWebSocket() as WebSocketContextType & {
    subscribeSeatUpdates: (id: string, callback: (data: SeatUpdateData) => void) => () => void;
    getCachedSeatData: (periodId: string) => SeatUpdateData | null;
  };
  const [seatData, setSeatData] = useState<SeatUpdateData | null>(null);
  const callbackId = useRef(`seat-updates-${Math.random().toString(36).substring(2, 9)}`);

  useEffect(() => {
    // Initialize with cached data if available
    if (periodId) {
      const cachedData = getCachedSeatData(periodId);
      if (cachedData) {
        setSeatData(cachedData);
      }
    }

    const unsubscribe = subscribeSeatUpdates(callbackId.current, (data: SeatUpdateData) => {
      // Only update if it's for the specific period or if no period filter is set
      if (!periodId || data.period_id === periodId) {
        setSeatData(data);
      }
    });

    return unsubscribe;
  }, [subscribeSeatUpdates, getCachedSeatData, periodId]);

  return {
    seatData,
    lastUpdate: seatData?.timestamp,
    updateType: seatData?.update_type,
  };
}

// Hook for connection status monitoring
export function useWebSocketStatus() {
  const { isConnected, connectionStatus, reconnect, connect, disconnect } = useWebSocket() as WebSocketContextType & {
    connect: () => void;
  };

  return {
    isConnected,
    connectionStatus,
    reconnect,
    connect,
    disconnect,
  };
}