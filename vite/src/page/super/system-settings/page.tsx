
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Badge } from '@/components/ui/badge';
import {
  Settings,
  Clock,
  Database,
  Trash2,
  AlertTriangle,
  BarChart3
} from 'lucide-react';
import {
  ExpirationTiersCard,
  ConfigManagementCard,
  ProcessExpiredCard,
  BookingLimitsCard
} from '@/components/manage/config';
import { useConfigInfo, useExpirationStats } from '@/stores/configStore';

export default function SystemSettingsPage() {
  const { data: configInfo } = useConfigInfo();
  const { data: stats } = useExpirationStats();

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">ตั้งค่าระบบ</h1>
          <p className="text-muted-foreground">
            จัดการการตั้งค่าต่างๆ ของระบบและการหมดอายุการจอง
          </p>
        </div>
        {configInfo && (
          <Badge variant="outline" className="font-mono">
            {configInfo.version}
          </Badge>
        )}
      </div>

      {/* Quick Stats */}
      {stats && (
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">การจองทั้งหมด</CardTitle>
              <BarChart3 className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stats.total_bookings_with_expiration}</div>
              <p className="text-xs text-muted-foreground">
                รายการที่มีการตั้งค่าหมดอายุ
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">หมดอายุแล้ว</CardTitle>
              <Trash2 className="h-4 w-4 text-destructive" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-destructive">{stats.expired_bookings}</div>
              <p className="text-xs text-muted-foreground">
                รอการประมวลผล
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">ใกล้หมดอายุ</CardTitle>
              <AlertTriangle className="h-4 w-4 text-yellow-600" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-yellow-600">{stats.near_expiring_bookings}</div>
              <p className="text-xs text-muted-foreground">
                ต้องติดตาม
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">เวลาเฉลี่ย</CardTitle>
              <Clock className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stats.average_expiration_time}</div>
              <p className="text-xs text-muted-foreground">
                ระยะเวลาหมดอายุเฉลี่ย
              </p>
            </CardContent>
          </Card>
        </div>
      )}

      <Tabs defaultValue="expiration" className="space-y-6">
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="expiration" className="flex items-center gap-2">
            <Clock className="h-4 w-4" />
            การหมดอายุ
          </TabsTrigger>
          <TabsTrigger value="booking-limits" className="flex items-center gap-2">
            <Settings className="h-4 w-4" />
            จำกัดการจอง
          </TabsTrigger>
          <TabsTrigger value="configs" className="flex items-center gap-2">
            <Database className="h-4 w-4" />
            การตั้งค่า
          </TabsTrigger>
          <TabsTrigger value="maintenance" className="flex items-center gap-2">
            <Trash2 className="h-4 w-4" />
            บำรุงรักษา
          </TabsTrigger>
        </TabsList>

        <TabsContent value="expiration" className="space-y-6">
          <ExpirationTiersCard />
        </TabsContent>

        <TabsContent value="booking-limits" className="space-y-6">
          <BookingLimitsCard />
        </TabsContent>

        <TabsContent value="configs" className="space-y-6">
          <ConfigManagementCard />
        </TabsContent>

        <TabsContent value="maintenance" className="space-y-6">
          <ProcessExpiredCard />
        </TabsContent>
      </Tabs>
    </div>
  );
}