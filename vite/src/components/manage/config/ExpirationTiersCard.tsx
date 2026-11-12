import { useState } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import {
  Settings,
  Clock,
  Calendar,
  Edit3,
  BarChart3,
  Loader2,
  CheckCircle,
  XCircle
} from 'lucide-react';
import { useExpirationTiers, useExpirationStats } from '@/stores/configStore';
import { ExpirationTiersDialog } from './ExpirationTiersDialog';
import { format } from 'date-fns';
import { th } from 'date-fns/locale';

export function ExpirationTiersCard() {
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const { data: tiersData, isLoading: isLoadingTiers, error: tiersError } = useExpirationTiers();
  const { data: stats } = useExpirationStats();

  const formatExpireTime = (tier: {expire_minutes?: number | null; expire_hours?: number | null; expire_days?: number | null}) => {
    console.log('formatExpireTime tier:', tier); // Debug log

    if (tier.expire_minutes) {
      return `${tier.expire_minutes} นาที`;
    }
    if (tier.expire_hours) {
      return `${tier.expire_hours} ชั่วโมง`;
    }
    if (tier.expire_days) {
      return `${tier.expire_days} วัน`;
    }
    return 'ไม่ระบุ';
  };

  const getTierColor = (sortOrder: number) => {
    switch (sortOrder) {
      case 0: return 'destructive'; // Most urgent (นาที)
      case 1: return 'destructive'; // Very urgent (ชั่วโมง)
      case 2: return 'default';     // Normal (วัน)
      case 3: return 'default';     // Normal (วัน)
      case 4: return 'secondary';   // Default (วัน)
      default: return 'secondary';
    }
  };

  if (isLoadingTiers) {
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

  if (tiersError) {
    return (
      <Card>
        <CardContent className="p-6">
          <div className="flex items-center text-destructive">
            <XCircle className="h-5 w-5 mr-2" />
            <span>เกิดข้อผิดพลาด: {tiersError.message}</span>
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
            <Settings className="h-5 w-5" />
            <div>
              <CardTitle className="text-lg">การตั้งค่าวันหมดอายุการจอง</CardTitle>
              <CardDescription>
                กำหนดระยะเวลาหมดอายุการจองตามประเภททัวร์
              </CardDescription>
            </div>
          </div>
          <Button onClick={() => setIsDialogOpen(true)}>
            <Edit3 className="h-4 w-4 mr-2" />
            แก้ไข
          </Button>
        </CardHeader>

        <CardContent className="space-y-4">
          {/* Tier Configuration */}
          <div className="space-y-3">
            <h4 className="font-medium flex items-center">
              <Clock className="h-4 w-4 mr-2" />
              ระดับการหมดอายุ
            </h4>

            {tiersData?.tiers?.map((tier) => (
              <div key={tier.name} className="flex items-center justify-between p-3 bg-muted rounded-lg">
                <div className="flex items-center space-x-3">
                  <Badge variant={getTierColor(tier.sort_order)}>
                    {tier.name}
                  </Badge>
                  <div>
                    <p className="font-medium">{tier.description}</p>
                    <p className="text-sm text-muted-foreground">
                      {tier.threshold_days
                        ? `ทัวร์ที่เริ่มภายใน ${tier.threshold_days} วัน`
                        : 'ทัวร์อื่นๆ ทั้งหมด'
                      }
                    </p>
                  </div>
                </div>
                <div className="text-right">
                  <p className="font-medium">{formatExpireTime(tier)}</p>
                  <p className="text-sm text-muted-foreground">หมดอายุใน</p>
                </div>
              </div>
            ))}
          </div>

          <Separator />

          {/* Statistics */}
          {stats && (
            <div className="space-y-3">
              <h4 className="font-medium flex items-center">
                <BarChart3 className="h-4 w-4 mr-2" />
                สถิติการจอง
              </h4>

              <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                <div className="bg-muted p-3 rounded-lg text-center">
                  <p className="text-2xl font-bold text-primary">{stats.total_bookings_with_expiration}</p>
                  <p className="text-sm text-muted-foreground">การจองทั้งหมด</p>
                </div>
                <div className="bg-muted p-3 rounded-lg text-center">
                  <p className="text-2xl font-bold text-destructive">{stats.expired_bookings}</p>
                  <p className="text-sm text-muted-foreground">หมดอายุแล้ว</p>
                </div>
                <div className="bg-muted p-3 rounded-lg text-center">
                  <p className="text-2xl font-bold text-yellow-600">{stats.near_expiring_bookings}</p>
                  <p className="text-sm text-muted-foreground">ใกล้หมดอายุ</p>
                </div>
                <div className="bg-muted p-3 rounded-lg text-center">
                  <p className="text-2xl font-bold text-green-600">
                    {stats.total_bookings_with_expiration - stats.expired_bookings - stats.near_expiring_bookings}
                  </p>
                  <p className="text-sm text-muted-foreground">ปกติ</p>
                </div>
              </div>

              {/* Bookings by Tier */}
              {Object.keys(stats.bookings_by_tier).length > 0 && (
                <div className="space-y-2">
                  <h5 className="text-sm font-medium">การจองแยกตามระดับ:</h5>
                  <div className="grid grid-cols-3 gap-2">
                    {Object.entries(stats.bookings_by_tier).map(([tierName, count]) => (
                      <div key={tierName} className="flex justify-between text-sm">
                        <span className="capitalize">{tierName}:</span>
                        <span className="font-medium">{count}</span>
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </div>
          )}

          {/* Metadata */}
          {tiersData?.metadata && (
            <>
              <Separator />
              <div className="flex items-center justify-between text-sm text-muted-foreground">
                <div className="flex items-center">
                  <CheckCircle className="h-4 w-4 mr-1" />
                  Version {tiersData.metadata.version}
                </div>
                <div className="flex items-center">
                  <Calendar className="h-4 w-4 mr-1" />
                  อัพเดทล่าสุด: {format(new Date(tiersData.metadata.last_updated), 'dd MMM yyyy HH:mm', { locale: th })}
                </div>
              </div>
            </>
          )}
        </CardContent>
      </Card>

      <ExpirationTiersDialog
        open={isDialogOpen}
        onOpenChange={setIsDialogOpen}
        initialData={tiersData}
      />
    </>
  );
}