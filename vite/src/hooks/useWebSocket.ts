// Re-export WebSocket hooks from hooks for convenience
export {
  useWebSocket,
  useSeatUpdates,
  useWebSocketStatus
} from '@/hooks/useWebSocketHooks';

// Additional specialized hooks
import { useEffect, useState, useCallback } from 'react';
import { useWebSocket } from '@/hooks/useWebSocketHooks';
import type { WebSocketMessage } from '@/types';

/**
 * Hook for WebSocket connection health monitoring
 */
export function useConnectionHealth() {
  const { isConnected, connectionStatus } = useWebSocket();
  const [connectionHistory, setConnectionHistory] = useState<{
    timestamp: string;
    status: string;
  }[]>([]);

  useEffect(() => {
    setConnectionHistory(prev => [
      ...prev.slice(-9), // Keep last 10 entries
      {
        timestamp: new Date().toISOString(),
        status: connectionStatus
      }
    ]);
  }, [connectionStatus]);

  const isHealthy = isConnected;
  const lastDisconnection = connectionHistory
    .slice()
    .reverse()
    .find(entry => entry.status === 'disconnected');

  return {
    isHealthy,
    connectionStatus,
    connectionHistory,
    lastDisconnection,
    isStable: connectionHistory.slice(-5).every(entry => entry.status === 'connected'),
  };
}

/**
 * Hook for WebSocket message debugging
 */
export function useWebSocketDebug() {
  const { lastMessage, isConnected, connectionStatus } = useWebSocket();
  const [messageHistory, setMessageHistory] = useState<(WebSocketMessage & { receivedAt: string })[]>([]);

  useEffect(() => {
    if (lastMessage) {
      setMessageHistory(prev => [
        ...prev.slice(-49), // Keep last 50 messages
        {
          ...lastMessage,
          receivedAt: new Date().toISOString(),
        }
      ]);
    }
  }, [lastMessage]);

  const clearHistory = useCallback(() => {
    setMessageHistory([]);
  }, []);

  const getMessageStats = useCallback(() => {
    const stats = messageHistory.reduce((acc, msg) => {
      const type = msg.data?.update_type || 'unknown';
      acc[type] = (acc[type] || 0) + 1;
      return acc;
    }, {} as Record<string, number>);

    return {
      total: messageHistory.length,
      byType: stats,
      latest: messageHistory[messageHistory.length - 1],
    };
  }, [messageHistory]);

  return {
    isConnected,
    connectionStatus,
    lastMessage,
    messageHistory,
    clearHistory,
    getMessageStats,
  };
}