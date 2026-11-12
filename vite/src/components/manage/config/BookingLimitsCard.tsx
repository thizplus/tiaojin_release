import { useState } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import {
  Shield,
  Users,
  Settings,
  Edit3,
  Loader2,
  XCircle,
  CheckCircle,
  AlertCircle,
  User
} from 'lucide-react';
import { useBookingLimitsConfig } from '@/stores/configStore';
import { useUser } from '@/hooks/useUser';
import { BookingLimitsDialog } from './BookingLimitsDialog';

export function BookingLimitsCard() {
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const { data: limitsConfig, isLoading, error } = useBookingLimitsConfig();

  // Get agent users data to display names instead of UUIDs
  const { usersList } = useUser(
    { role: 'agent', limit: 100 },
    undefined,
    { enableList: true, enableById: false, enableStats: false }
  );

  console.log('🔍 BookingLimitsCard Debug:', {
    limitsConfig,
    isLoading,
    error,
    hasLimitsConfig: !!limitsConfig,
    usersList: usersList.length
  });

  const getAgentDisplayInfo = (agentId: string) => {
    const agent = usersList.find(user => user.id === agentId);
    if (agent) {
      return {
        email: agent.email,
        role: agent.role,
        found: true
      };
    }
    return {
      email: agentId,
      role: 'unknown',
      found: false
    };
  };

  if (isLoading) {
    return (
      <Card>
        <CardContent className="p-6">
          <div className="flex items-center justify-center">
            <Loader2 className="h-6 w-6 animate-spin" />
            <span className="ml-2">กำลังโหลดข้อมูล...</span>
          </div>
        </CardContent>
      </Card>
    );
  }

  if (error) {
    return (
      <Card>
        <CardContent className="p-6">
          <div className="flex items-center text-destructive">
            <XCircle className="h-5 w-5 mr-2" />
            <span>เกิดข้อผิดพลาด: {error.message}</span>
          </div>
        </CardContent>
      </Card>
    );
  }

  if (!limitsConfig) {
    return (
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <div className="flex items-center space-x-2">
            <Shield className="h-5 w-5" />
            <div>
              <CardTitle className="text-lg">การจำกัดการจอง</CardTitle>
              <CardDescription>
                กำหนดข้อจำกัดการจองและผู้ใช้ที่ได้รับการยกเว้น
              </CardDescription>
            </div>
          </div>
          <Button onClick={() => setIsDialogOpen(true)}>
            <Settings className="h-4 w-4 mr-2" />
            ตั้งค่า
          </Button>
        </CardHeader>

        <CardContent>
          <div className="text-center py-8 text-muted-foreground">
            <AlertCircle className="h-8 w-8 mx-auto mb-2" />
            <p>ยังไม่มีการตั้งค่าการจำกัดการจอง</p>
            <p className="text-sm">กดปุ่มตั้งค่าเพื่อเริ่มต้น</p>
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <>
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <div className="flex items-center space-x-2">
            <Shield className="h-5 w-5" />
            <div>
              <CardTitle className="text-lg">การจำกัดการจอง</CardTitle>
              <CardDescription>
                กำหนดข้อจำกัดการจองและผู้ใช้ที่ได้รับการยกเว้น
              </CardDescription>
            </div>
          </div>
          <Button onClick={() => setIsDialogOpen(true)}>
            <Edit3 className="h-4 w-4 mr-2" />
            แก้ไข
          </Button>
        </CardHeader>

        <CardContent className="space-y-4">
          {/* Status */}
          <div className="flex items-center justify-between p-3 bg-muted rounded-lg">
            <div className="flex items-center space-x-2">
              <CheckCircle className="h-4 w-4" />
              <span className="font-medium">สถานะระบบ</span>
            </div>
            <Badge variant={limitsConfig.global_settings.enabled ? 'default' : 'secondary'}>
              {limitsConfig.global_settings.enabled ? 'เปิดใช้งาน' : 'ปิดใช้งาน'}
            </Badge>
          </div>

          {/* Global Settings */}
          <div className="space-y-3">
            <h4 className="font-medium flex items-center">
              <Settings className="h-4 w-4 mr-2" />
              การตั้งค่าทั่วไป
            </h4>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="bg-muted p-3 rounded-lg">
                <p className="text-sm text-muted-foreground">จำนวนจำกัดต่อรอบ</p>
                <p className="text-2xl font-bold text-primary">
                  {limitsConfig.global_settings.default_limit_per_period}
                </p>
                <p className="text-xs text-muted-foreground">การจองสูงสุด</p>
              </div>

              <div className="bg-muted p-3 rounded-lg">
                <p className="text-sm text-muted-foreground">สถานะที่นับรวม</p>
                <p className="text-lg font-bold text-primary">
                  {limitsConfig.global_settings.count_all_statuses.length}
                </p>
                <p className="text-xs text-muted-foreground">ประเภทสถานะ</p>
              </div>
            </div>

            {/* Status List */}
            <div className="space-y-2">
              <p className="text-sm font-medium">สถานะการจองที่นับรวม:</p>
              <div className="flex flex-wrap gap-2">
                {limitsConfig.global_settings.count_all_statuses.map((status) => (
                  <Badge key={status} variant="outline" className="text-xs">
                    {status}
                  </Badge>
                ))}
              </div>
            </div>
          </div>

          <Separator />

          {/* Agent Bypasses */}
          <div className="space-y-3">
            <h4 className="font-medium flex items-center">
              <Users className="h-4 w-4 mr-2" />
              ผู้ใช้ที่ได้รับการยกเว้น
              <Badge variant="secondary" className="ml-2 text-xs">
                {limitsConfig.agent_global_bypasses.length}
              </Badge>
            </h4>

            {limitsConfig.agent_global_bypasses.length === 0 ? (
              <div className="text-center py-4 text-muted-foreground">
                <p className="text-sm">ไม่มีผู้ใช้ที่ได้รับการยกเว้น</p>
              </div>
            ) : (
              <div className="space-y-2">
                <div className="grid grid-cols-1 gap-2 max-h-32 overflow-y-auto">
                  {limitsConfig.agent_global_bypasses.map((agentId, index) => {
                    const agentInfo = getAgentDisplayInfo(agentId);
                    return (
                      <div
                        key={agentId}
                        className="flex items-center justify-between p-2 bg-muted rounded border"
                      >
                        <div className="flex items-center space-x-2 min-w-0 flex-1">
                          <User className="h-3 w-3 text-muted-foreground shrink-0" />
                          <div className="min-w-0 flex-1">
                            <p className="font-mono text-xs truncate">
                              {agentInfo.email}
                            </p>
                            {agentInfo.found && (
                              <p className="text-xs text-muted-foreground">
                                ID: {agentId}
                              </p>
                            )}
                          </div>
                        </div>
                        <div className="flex items-center gap-1 shrink-0">
                          <Badge variant={agentInfo.found ? 'default' : 'destructive'} className="text-xs">
                            {agentInfo.role}
                          </Badge>
                          <Badge variant="outline" className="text-xs">
                            #{index + 1}
                          </Badge>
                        </div>
                      </div>
                    );
                  })}
                </div>
              </div>
            )}
          </div>
        </CardContent>
      </Card>

      <BookingLimitsDialog
        open={isDialogOpen}
        onOpenChange={setIsDialogOpen}
        initialData={limitsConfig}
      />
    </>
  );
}