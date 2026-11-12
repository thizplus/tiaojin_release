import { useState } from 'react';
import { format } from 'date-fns';
import { th } from 'date-fns/locale';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog';
import {
  Users,
  Clock,
  Calendar,
  FileText,
  Unlock,
  AlertTriangle,
  CheckCircle,
  XCircle,
} from 'lucide-react';
import { useReleaseSeatHold } from '@/hooks/useSeatHold';
import { SEAT_HOLD_STATUS, SEAT_HOLD_STATUS_COLORS } from '@/constants/seat-hold';
import type { SeatHold } from '@/types/models/seat-hold';

interface SeatHoldCardProps {
  seatHold: SeatHold;
}

export function SeatHoldCard({ seatHold }: SeatHoldCardProps) {
  const [showReleaseDialog, setShowReleaseDialog] = useState(false);
  const releaseMutation = useReleaseSeatHold();

  const handleRelease = async () => {
    try {
      await releaseMutation.mutateAsync({
        holdId: seatHold.id,
        data: { reason: 'ปล่อยโดยผู้ดูแลระบบ' },
      });
      setShowReleaseDialog(false);
    } catch (error) {
      console.error('Failed to release seat hold:', error);
    }
  };

  const getStatusIcon = (status: SeatHold['status']) => {
    switch (status) {
      case 'active':
        return <CheckCircle className="h-4 w-4 text-green-600" />;
      case 'released':
        return <Unlock className="h-4 w-4 text-gray-600" />;
      case 'expired':
        return <XCircle className="h-4 w-4 text-red-600" />;
      default:
        return <AlertTriangle className="h-4 w-4 text-yellow-600" />;
    }
  };

  const formatCreatedAt = (dateString: string) => {
    return format(new Date(dateString), 'dd MMM yy HH:mm', { locale: th });
  };

  const formatExpiresAt = (dateString?: string) => {
    if (!dateString) return 'ไม่กำหนด';
    return format(new Date(dateString), 'dd MMM yy HH:mm', { locale: th });
  };

  return (
    <>
      <Card className="hover:shadow-sm  transition-shadow leading-none">
        <CardContent className="px-4">
          <div className="flex items-center justify-between mb-3">
            <div className="flex items-center space-x-2">
              <Users className="h-4 w-4 text-muted-foreground" />
              <span className="font-semibold">{seatHold.held_seats} ที่นั่ง</span>
              {seatHold.time_remaining && !seatHold.is_expired && (
                <Badge variant="outline" className=" bg-muted text-muted-foreground">
                  {seatHold.time_remaining}
                </Badge>
              )}
              {seatHold.is_expired && (
                <Badge variant="outline" className=" bg-muted text-muted-foreground">
                  หมดอายุ
                </Badge>
              )}
            </div>
            <div className="flex items-center space-x-2">
              <Badge
                variant="outline"
                className={`flex items-center space-x-1 ${SEAT_HOLD_STATUS_COLORS[seatHold.status]}`}
              >
                {getStatusIcon(seatHold.status)}
                <span>{SEAT_HOLD_STATUS[seatHold.status]}</span>
              </Badge>
              {seatHold.can_release && seatHold.status === 'active' && (
                <Badge
                  variant="outline"
                  className="flex items-center space-x-1  text-foreground bg-primary  cursor-pointer"
                  onClick={() => setShowReleaseDialog(true)}
                >
                  <Unlock className="h-3 w-3" />
                  <span>ปล่อย</span>
                </Badge>
              )}
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4 text-sm">
            <div className="flex items-center space-x-2">
              <Clock className="h-4 w-4 text-muted-foreground" />
              <span className="text-muted-foreground">สร้าง: {formatCreatedAt(seatHold.created_at)}</span>
            </div>
            <div className="flex items-center space-x-2">
              <Calendar className="h-4 w-4 text-muted-foreground" />
              <span className="text-muted-foreground">หมดอายุ: {formatExpiresAt(seatHold.expires_at)}</span>
            </div>
          </div>

          {seatHold.reason && (
            <div className=" text-sm">
              <div className="flex items-center space-x-2">
                <FileText className="h-4 w-4 text-muted-foreground" />
                <span className="text-muted-foreground">เหตุผล: {seatHold.reason}</span>
              </div>
            </div>
          )}

        </CardContent>
      </Card>

      {/* Release Confirmation Dialog */}
      <AlertDialog open={showReleaseDialog} onOpenChange={setShowReleaseDialog}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle className="flex items-center space-x-2">
              <AlertTriangle className="h-5 w-5 text-destructive" />
              <span>ยืนยันการปล่อยที่นั่งจำนอง</span>
            </AlertDialogTitle>
            <AlertDialogDescription>
              คุณต้องการปล่อยที่นั่งจำนอง <strong>{seatHold.held_seats} ที่นั่ง</strong> ใช่หรือไม่?
              <br />
              <br />
              การดำเนินการนี้ไม่สามารถย้อนกลับได้ และที่นั่งเหล่านี้จะกลับมาว่างให้จองได้ทันที
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>ยกเลิก</AlertDialogCancel>
            <AlertDialogAction
              onClick={handleRelease}
              disabled={releaseMutation.isPending}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              {releaseMutation.isPending ? (
                <Clock className="mr-2 h-4 w-4 animate-spin" />
              ) : (
                <Unlock className="mr-2 h-4 w-4" />
              )}
              ปล่อยที่นั่ง
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}