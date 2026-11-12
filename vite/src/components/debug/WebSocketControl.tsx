
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { WebSocketStatus } from '@/components/ui/websocket-status';
import { useWebSocketStatus } from '@/hooks/useWebSocket';
import {
  Wifi,
  WifiOff,
  RotateCcw,
  Settings
} from 'lucide-react';

export function WebSocketControl() {
  const { isConnected, connectionStatus, connect, disconnect, reconnect } = useWebSocketStatus();

  if (!import.meta.env.DEV) {
    return null; // Only show in development
  }

  return (
    <Card className="fixed top-4 right-4 z-50 w-80 shadow-lg">
      <CardHeader className="pb-3">
        <CardTitle className="flex items-center gap-2 text-sm">
          <Settings className="h-4 w-4" />
          WebSocket Control (Dev)
        </CardTitle>
      </CardHeader>

      <CardContent className="space-y-3">
        <WebSocketStatus compact={true} showReconnectButton={false} />

        <div className="flex gap-2">
          {!isConnected ? (
            <Button
              size="sm"
              onClick={connect}
              disabled={connectionStatus === 'connecting'}
              className="flex-1"
            >
              <Wifi className="h-3 w-3 mr-1" />
              Connect
            </Button>
          ) : (
            <Button
              size="sm"
              variant="outline"
              onClick={disconnect}
              className="flex-1"
            >
              <WifiOff className="h-3 w-3 mr-1" />
              Disconnect
            </Button>
          )}

          <Button
            size="sm"
            variant="outline"
            onClick={reconnect}
            disabled={connectionStatus === 'connecting'}
            className="flex-1"
          >
            <RotateCcw className="h-3 w-3 mr-1" />
            Reconnect
          </Button>
        </div>

        <div className="text-xs text-muted-foreground">
          Status: {connectionStatus}
        </div>
      </CardContent>
    </Card>
  );
}