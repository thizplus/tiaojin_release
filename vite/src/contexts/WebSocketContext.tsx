import React, { useEffect, useState, useRef, useCallback } from 'react';
import { WebSocketManager } from '@/lib/websocket-manager';
import { useToast } from '@/hooks/use-toast';
import { WebSocketContext } from '@/hooks/useWebSocketHooks';
import type {
  WebSocketContextType,
  WebSocketConnectionStatus,
  WebSocketMessage,
  SeatUpdateData,
  SeatUpdateType,
  NotificationConfig
} from '@/types';

interface WebSocketProviderProps {
  children: React.ReactNode;
  url?: string;
  autoConnect?: boolean;
  showNotifications?: boolean;
}


export function WebSocketProvider({
  children,
  url = 'ws://localhost:8080/ws',
  autoConnect = true,
  showNotifications = true
}: WebSocketProviderProps) {
  const { toast } = useToast();
  const [connectionStatus, setConnectionStatus] = useState<WebSocketConnectionStatus>('disconnected');
  const [lastMessage, setLastMessage] = useState<WebSocketMessage | null>(null);
  const wsManagerRef = useRef<WebSocketManager | null>(null);

  // Global cache for seat data by period_id
  const seatDataCacheRef = useRef<Map<string, SeatUpdateData>>(new Map());


  // Seat update callbacks registry
  const seatUpdateCallbacksRef = useRef<Map<string, (data: SeatUpdateData) => void>>(new Map());

  // Initialize WebSocket Manager
  useEffect(() => {
    const manager = new WebSocketManager({
      url,
      reconnectInterval: 3000,
      maxReconnectAttempts: 10,
      debug: import.meta.env.DEV
    });

    wsManagerRef.current = manager;

    // Set up event listeners
    manager.onConnectionStatusChange((status) => {
      setConnectionStatus(status);

      if (showNotifications) {
        switch (status) {
          case 'connected':
            toast({
              title: "เชื่อมต่อสำเร็จ",
              description: "ระบบ Real-time พร้อมใช้งาน",
            });
            break;
          case 'disconnected':
            toast({
              title: "การเชื่อมต่อขาดหาย",
              description: "กำลังพยายามเชื่อมต่อใหม่...",
              variant: "destructive",
            });
            break;
          case 'error':
            toast({
              title: "เกิดข้อผิดพลาด",
              description: "ไม่สามารถเชื่อมต่อ WebSocket ได้",
              variant: "destructive",
            });
            break;
        }
      }
    });

    manager.onMessageReceived((message) => {
      setLastMessage(message);
    });

    manager.onSeatUpdateReceived((data) => {
      // Cache seat data by period_id
      seatDataCacheRef.current.set(data.period_id, data);

      // Notify all registered callbacks
      seatUpdateCallbacksRef.current.forEach((callback) => {
        callback(data);
      });

      // Show notification based on update type
      if (showNotifications) {
        const notification = getSeatUpdateNotification(data.update_type);
        if (notification) {
          toast({
            title: notification.title,
            description: notification.description,
            variant: notification.type === 'error' ? 'destructive' : 'default',
          });
        }
      }
    });

    manager.onErrorOccurred((error) => {
      console.error('WebSocket Error:', error);
    });

    // Cleanup on unmount
    return () => {
      manager.destroy();
      wsManagerRef.current = null;
    };
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [url]);

  // Separate useEffect for autoConnect to prevent dependency array size changes
  useEffect(() => {
    if (autoConnect && wsManagerRef.current) {
      wsManagerRef.current.connect();
    }
  }, [autoConnect]);

  // Helper function to get notification config for seat updates
  const getSeatUpdateNotification = (updateType: SeatUpdateType): NotificationConfig | null => {
    const notifications: Record<SeatUpdateType, NotificationConfig> = {
      seat_held: {
        type: 'info',
        title: 'ที่นั่งถูกจอง',
        description: 'มีการจองที่นั่งชั่วคราว',
      },
      seat_released: {
        type: 'success',
        title: 'ที่นั่งถูกปล่อย',
        description: 'ที่นั่งพร้อมให้จองแล้ว',
      },
      seat_hold_expired: {
        type: 'info',
        title: 'การจองหมดอายุ',
        description: 'การจองชั่วคราวหมดอายุแล้ว',
      },
      booking_created: {
        type: 'success',
        title: 'มีการจองใหม่',
        description: 'มีลูกค้าจองทัวร์เรียบร้อยแล้ว',
      },
      booking_updated: {
        type: 'info',
        title: 'อัพเดทการจอง',
        description: 'ข้อมูลการจองถูกแก้ไข',
      },
      booking_cancelled: {
        type: 'warning',
        title: 'ยกเลิกการจอง',
        description: 'การจองถูกยกเลิกแล้ว',
      },
      booking_expired: {
        type: 'warning',
        title: 'การจองหมดอายุ',
        description: 'การจองหมดอายุและถูกยกเลิกแล้ว',
      },
    };

    return notifications[updateType] || null;
  };

  // Context methods
  const sendMessage = useCallback((message: Record<string, unknown>): boolean => {
    return wsManagerRef.current?.sendMessage(message) || false;
  }, []);

  const reconnect = useCallback(() => {
    wsManagerRef.current?.reconnect();
  }, []);

  const disconnect = useCallback(() => {
    wsManagerRef.current?.disconnect();
  }, []);

  const connect = useCallback(() => {
    wsManagerRef.current?.connect();
  }, []);

  // Register/unregister seat update callbacks
  const subscribeSeatUpdates = useCallback((id: string, callback: (data: SeatUpdateData) => void) => {
    seatUpdateCallbacksRef.current.set(id, callback);

    return () => {
      seatUpdateCallbacksRef.current.delete(id);
    };
  }, []);

  // Get cached seat data for a specific period
  const getCachedSeatData = useCallback((periodId: string): SeatUpdateData | null => {
    return seatDataCacheRef.current.get(periodId) || null;
  }, []);

  const value: WebSocketContextType & {
    connect: () => void;
    subscribeSeatUpdates: (id: string, callback: (data: SeatUpdateData) => void) => () => void;
    getCachedSeatData: (periodId: string) => SeatUpdateData | null;
  } = {
    isConnected: connectionStatus === 'connected',
    connectionStatus,
    lastMessage,
    sendMessage,
    reconnect,
    disconnect,
    connect,
    subscribeSeatUpdates,
    getCachedSeatData,
  };

  return (
    <WebSocketContext.Provider value={value}>
      {children}
    </WebSocketContext.Provider>
  );
}

