import  { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';
import {
  Activity,
  Trash2,
  Wifi,
  WifiOff,
  RotateCcw,
  //MessageSquare,
  BarChart3,
  Clock,
  Bug
} from 'lucide-react';
import { useWebSocketDebug, useWebSocketStatus } from '@/hooks/useWebSocket';
import { format } from 'date-fns';

export function WebSocketDebugPanel() {
  const {
    isConnected,
    connectionStatus,
    lastMessage,
    messageHistory,
    clearHistory,
    getMessageStats
  } = useWebSocketDebug();

  const { reconnect, disconnect } = useWebSocketStatus();

  const [isOpen, setIsOpen] = useState(false);

  const stats = getMessageStats();

  const getStatusColor = () => {
    switch (connectionStatus) {
      case 'connected': return 'text-green-600';
      case 'connecting': return 'text-blue-600';
      case 'reconnecting': return 'text-yellow-600';
      case 'error': return 'text-red-600';
      default: return 'text-gray-600';
    }
  };

  const getStatusIcon = () => {
    switch (connectionStatus) {
      case 'connected': return <Wifi className="h-4 w-4" />;
      case 'connecting': return <Activity className="h-4 w-4 animate-spin" />;
      case 'reconnecting': return <RotateCcw className="h-4 w-4 animate-spin" />;
      default: return <WifiOff className="h-4 w-4" />;
    }
  };

  if (!import.meta.env.DEV) {
    return null; // Only show in development
  }

  return (
    <Dialog open={isOpen} onOpenChange={setIsOpen}>
      <DialogTrigger asChild>
        <Button
          variant="outline"
          size="sm"
          className="fixed bottom-4 right-4 z-50 shadow-lg"
        >
          <Bug className="h-4 w-4 mr-2" />
          WebSocket Debug
        </Button>
      </DialogTrigger>

      <DialogContent className="max-w-4xl max-h-[80vh]">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Bug className="h-5 w-5" />
            WebSocket Debug Panel
          </DialogTitle>
          <DialogDescription>
            Monitor WebSocket connection and message history
          </DialogDescription>
        </DialogHeader>

        <Tabs defaultValue="status" className="w-full">
          <TabsList className="grid w-full grid-cols-4">
            <TabsTrigger value="status">Status</TabsTrigger>
            <TabsTrigger value="messages">Messages</TabsTrigger>
            <TabsTrigger value="stats">Statistics</TabsTrigger>
            <TabsTrigger value="controls">Controls</TabsTrigger>
          </TabsList>

          <TabsContent value="status" className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  {getStatusIcon()}
                  Connection Status
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <div className="text-sm text-muted-foreground">Status</div>
                    <Badge variant={isConnected ? 'default' : 'destructive'} className={getStatusColor()}>
                      {connectionStatus}
                    </Badge>
                  </div>
                  <div className="space-y-2">
                    <div className="text-sm text-muted-foreground">Connected</div>
                    <Badge variant={isConnected ? 'default' : 'outline'}>
                      {isConnected ? 'Yes' : 'No'}
                    </Badge>
                  </div>
                </div>

                {lastMessage && (
                  <div className="space-y-2">
                    <div className="text-sm text-muted-foreground">Last Message</div>
                    <div className="p-3 bg-muted rounded text-xs">
                      <div className="font-mono">
                        {JSON.stringify(lastMessage, null, 2)}
                      </div>
                    </div>
                  </div>
                )}
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="messages" className="space-y-4">
            <div className="flex items-center justify-between">
              <h3 className="text-lg font-medium">Message History</h3>
              <Button
                variant="outline"
                size="sm"
                onClick={clearHistory}
                disabled={messageHistory.length === 0}
              >
                <Trash2 className="h-4 w-4 mr-2" />
                Clear
              </Button>
            </div>

            <ScrollArea className="h-[400px] border rounded">
              <div className="p-4 space-y-2">
                {messageHistory.length === 0 ? (
                  <div className="text-center text-muted-foreground py-8">
                    No messages yet
                  </div>
                ) : (
                  messageHistory.slice().reverse().map((message, index) => (
                    <Card key={index} className="p-3">
                      <div className="flex items-center justify-between mb-2">
                        <Badge variant="outline">
                          {message.type}
                        </Badge>
                        <div className="text-xs text-muted-foreground flex items-center gap-1">
                          <Clock className="h-3 w-3" />
                          {format(new Date(message.receivedAt), 'HH:mm:ss.SSS')}
                        </div>
                      </div>
                      <div className="text-xs font-mono bg-muted p-2 rounded overflow-x-auto">
                        {JSON.stringify(message.data, null, 2)}
                      </div>
                    </Card>
                  ))
                )}
              </div>
            </ScrollArea>
          </TabsContent>

          <TabsContent value="stats" className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <BarChart3 className="h-5 w-5" />
                  Message Statistics
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <div className="text-sm text-muted-foreground">Total Messages</div>
                    <div className="text-2xl font-bold">{stats.total}</div>
                  </div>
                  <div className="space-y-2">
                    <div className="text-sm text-muted-foreground">Latest Update</div>
                    <div className="text-sm">
                      {stats.latest ? format(new Date(stats.latest.receivedAt), 'HH:mm:ss') : 'None'}
                    </div>
                  </div>
                </div>

                <Separator />

                <div className="space-y-2">
                  <div className="text-sm text-muted-foreground">By Update Type</div>
                  <div className="grid gap-2">
                    {Object.entries(stats.byType).map(([type, count]) => (
                      <div key={type} className="flex items-center justify-between text-sm">
                        <span className="capitalize">{type.replace('_', ' ')}</span>
                        <Badge variant="outline">{count as number}</Badge>
                      </div>
                    ))}
                  </div>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="controls" className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle>Connection Controls</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  <Button
                    variant="outline"
                    onClick={reconnect}
                    disabled={connectionStatus === 'connecting'}
                  >
                    <RotateCcw className="h-4 w-4 mr-2" />
                    Reconnect
                  </Button>
                  <Button
                    variant="outline"
                    onClick={disconnect}
                    disabled={!isConnected}
                  >
                    <WifiOff className="h-4 w-4 mr-2" />
                    Disconnect
                  </Button>
                </div>

                <Separator />

                <div className="space-y-2">
                  <div className="text-sm text-muted-foreground">Actions</div>
                  <Button
                    variant="outline"
                    onClick={clearHistory}
                    disabled={messageHistory.length === 0}
                    className="w-full"
                  >
                    <Trash2 className="h-4 w-4 mr-2" />
                    Clear Message History
                  </Button>
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </DialogContent>
    </Dialog>
  );
}