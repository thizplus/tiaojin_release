
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
 
  WifiOff,
  RotateCcw,
  AlertCircle,
  CheckCircle2,
  Loader2
} from 'lucide-react';
import { useWebSocketStatus } from '@/hooks/useWebSocket';
import { cn } from '@/lib/utils';
import type { WebSocketConnectionStatus } from '@/types';

interface WebSocketStatusProps {
  className?: string;
  showReconnectButton?: boolean;
  compact?: boolean;
}

export function WebSocketStatus({
  className,
  showReconnectButton = true,
  compact = false
}: WebSocketStatusProps) {
  const { isConnected, connectionStatus, reconnect } = useWebSocketStatus();

  const getStatusConfig = (status: WebSocketConnectionStatus) => {
    switch (status) {
      case 'connected':
        return {
          icon: CheckCircle2,
          text: 'เชื่อมต่อแล้ว',
          variant: 'default' as const,
          color: 'text-green-600',
          bgColor: 'bg-green-50',
          borderColor: 'border-green-200',
        };
      case 'connecting':
        return {
          icon: Loader2,
          text: 'กำลังเชื่อมต่อ...',
          variant: 'secondary' as const,
          color: 'text-blue-600',
          bgColor: 'bg-blue-50',
          borderColor: 'border-blue-200',
          animate: true,
        };
      case 'reconnecting':
        return {
          icon: RotateCcw,
          text: 'กำลังเชื่อมต่อใหม่...',
          variant: 'secondary' as const,
          color: 'text-yellow-600',
          bgColor: 'bg-yellow-50',
          borderColor: 'border-yellow-200',
          animate: true,
        };
      case 'error':
        return {
          icon: AlertCircle,
          text: 'เกิดข้อผิดพลาด',
          variant: 'destructive' as const,
          color: 'text-red-600',
          bgColor: 'bg-red-50',
          borderColor: 'border-red-200',
        };
      case 'disconnected':
      default:
        return {
          icon: WifiOff,
          text: 'ไม่เชื่อมต่อ',
          variant: 'outline' as const,
          color: 'text-gray-600',
          bgColor: 'bg-gray-50',
          borderColor: 'border-gray-200',
        };
    }
  };

  const config = getStatusConfig(connectionStatus);
  const Icon = config.icon;

  if (compact) {
    return (
      <div className={cn('flex items-center gap-2', className)}>
        <div className={cn(
          'flex items-center gap-1 px-2 py-1 rounded text-xs border',
          config.bgColor,
          config.borderColor
        )}>
          <Icon className={cn(
            'h-3 w-3',
            config.color,
            config.animate && 'animate-spin'
          )} />
          <span className={config.color}>{config.text}</span>
        </div>
        {showReconnectButton && !isConnected && (
          <Button
            size="sm"
            variant="outline"
            onClick={reconnect}
            className="h-6 px-2 text-xs"
          >
            <RotateCcw className="h-3 w-3" />
          </Button>
        )}
      </div>
    );
  }

  return (
    <div className={cn('flex items-center gap-2', className)}>
      <Badge variant={config.variant} className="flex items-center gap-1">
        <Icon className={cn(
          'h-3 w-3',
          config.animate && 'animate-spin'
        )} />
        {config.text}
      </Badge>
      {showReconnectButton && !isConnected && (
        <Button
          size="sm"
          variant="outline"
          onClick={reconnect}
          className="h-7"
        >
          <RotateCcw className="h-4 w-4 mr-1" />
          เชื่อมต่อใหม่
        </Button>
      )}
    </div>
  );
}

// Connection indicator for minimal display
export function WebSocketIndicator({ className }: { className?: string }) {
  const {  connectionStatus } = useWebSocketStatus();

  const getIndicatorColor = () => {
    switch (connectionStatus) {
      case 'connected':
        return 'bg-green-500';
      case 'connecting':
      case 'reconnecting':
        return 'bg-yellow-500 animate-pulse';
      case 'error':
        return 'bg-red-500';
      case 'disconnected':
      default:
        return 'bg-gray-400';
    }
  };

  return (
    <div
      className={cn(
        'h-2 w-2 rounded-full',
        getIndicatorColor(),
        className
      )}
      title={`WebSocket: ${connectionStatus}`}
    />
  );
}