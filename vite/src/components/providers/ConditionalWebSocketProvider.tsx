import { WebSocketProvider } from '@/contexts/WebSocketContext';
//import { WebSocketDebugPanel } from '@/components/debug/WebSocketDebugPanel';
//import { WebSocketControl } from '@/components/debug/WebSocketControl';
import { useAuth } from '@/hooks/useAuth';

interface ConditionalWebSocketProviderProps {
  children: React.ReactNode;
}

export function ConditionalWebSocketProvider({ children }: ConditionalWebSocketProviderProps) {
  const { isAuthenticated } = useAuth();

  // Only provide WebSocket when user is authenticated
  if (isAuthenticated) {
    return (
      <WebSocketProvider
        url={import.meta.env.VITE_WS_URL || 'ws://localhost:8080/ws'}
        autoConnect={true}
        showNotifications={import.meta.env.VITE_WS_NOTIFICATIONS !== 'false'}
      >
        {children}
   {/**
    * 
    *   <WebSocketControl />
        <WebSocketDebugPanel />
    * 
    */}
      </WebSocketProvider>
    );
  }

  // Return children without WebSocket when not authenticated
  return <>{children}</>;
}