import { useState } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Badge } from '@/components/ui/badge';
import { Switch } from '@/components/ui/switch';
import { Separator } from '@/components/ui/separator';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
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
  Trash2,
  Play,
  AlertTriangle,
  CheckCircle,
  XCircle,
 
  Loader2,
  Eye
} from 'lucide-react';
import { useProcessExpiredBookings } from '@/stores/configStore';
import { useToast } from '@/hooks/use-toast';
import { format } from 'date-fns';
import { th } from 'date-fns/locale';
import type { ProcessExpiredResult } from '@/types/config';

export function ProcessExpiredCard() {
  const { toast } = useToast();
  const processMutation = useProcessExpiredBookings();

  const [isDryRun, setIsDryRun] = useState(true);
  const [limit, setLimit] = useState(100);
  const [result, setResult] = useState<ProcessExpiredResult | null>(null);
  const [showConfirmDialog, setShowConfirmDialog] = useState(false);

  const handleProcess = async (dryRun: boolean = isDryRun) => {
    try {
      const processResult = await processMutation.mutateAsync({
        dry_run: dryRun,
        limit: limit || undefined,
      });

      setResult(processResult);

      if (dryRun) {
        toast({
          title: "ทดสอบการประมวลผลสำเร็จ",
          description: `พบการจองที่หมดอายุ ${processResult.expired_found} รายการ`,
        });

        if (processResult.expired_found > 0) {
          setShowConfirmDialog(true);
        }
      } else {
        toast({
          title: "ประมวลผลสำเร็จ",
          description: `ประมวลผลแล้ว ${processResult.processed_count} รายการ`,
        });
      }
    } catch (error) {
      toast({
        title: "เกิดข้อผิดพลาด",
        description: error instanceof Error ? error.message : "ไม่สามารถประมวลผลได้",
        variant: "destructive",
      });
    }
  };

  const handleConfirmProcess = () => {
    setShowConfirmDialog(false);
    handleProcess(false);
  };

  const getStatusColor = (success: boolean) => {
    return success ? 'default' : 'destructive';
  };

  return (
    <>
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-4">
          <div className="flex items-center space-x-2">
            <Trash2 className="h-5 w-5" />
            <div>
              <CardTitle className="text-lg">ประมวลผลการจองที่หมดอายุ</CardTitle>
              <CardDescription>
                ตรวจสอบและจัดการการจองที่หมดอายุแล้ว
              </CardDescription>
            </div>
          </div>
        </CardHeader>

        <CardContent className="space-y-6">
          {/* Controls */}
          <div className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="limit">จำนวนที่ตรวจสอบ</Label>
                <Input
                  id="limit"
                  type="number"
                  value={limit}
                  onChange={(e) => setLimit(parseInt(e.target.value) || 100)}
                  min="1"
                  max="1000"
                  placeholder="100"
                />
              </div>
              <div className="space-y-2">
                <Label>โหมดการทำงาน</Label>
                <div className="flex items-center space-x-2">
                  <Switch
                    id="dry_run"
                    checked={isDryRun}
                    onCheckedChange={setIsDryRun}
                  />
                  <Label htmlFor="dry_run">
                    {isDryRun ? 'ทดสอบ (ไม่ลบจริง)' : 'ลบจริง'}
                  </Label>
                </div>
              </div>
            </div>

            <div className="flex gap-2">
              <Button
                onClick={() => handleProcess(true)}
                disabled={processMutation.isPending}
                variant="outline"
              >
                {processMutation.isPending ? (
                  <Loader2 className="h-4 w-4 animate-spin mr-2" />
                ) : (
                  <Eye className="h-4 w-4 mr-2" />
                )}
                ทดสอบ
              </Button>
              <Button
                onClick={() => isDryRun ? handleProcess(true) : setShowConfirmDialog(true)}
                disabled={processMutation.isPending}
                variant={isDryRun ? "default" : "destructive"}
              >
                {processMutation.isPending ? (
                  <Loader2 className="h-4 w-4 animate-spin mr-2" />
                ) : (
                  <Play className="h-4 w-4 mr-2" />
                )}
                {isDryRun ? 'ทดสอบการประมวลผล' : 'ประมวลผลจริง'}
              </Button>
            </div>
          </div>

          {/* Results */}
          {result && (
            <>
              <Separator />
              <div className="space-y-4">
                <h4 className="font-medium flex items-center">
                  <CheckCircle className="h-4 w-4 mr-2" />
                  ผลการประมวลผล
                </h4>

                {/* Summary */}
                <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                  <div className="bg-muted p-3 rounded-lg text-center">
                    <p className="text-2xl font-bold text-primary">{result.total_checked}</p>
                    <p className="text-sm text-muted-foreground">ตรวจสอบแล้ว</p>
                  </div>
                  <div className="bg-muted p-3 rounded-lg text-center">
                    <p className="text-2xl font-bold text-destructive">{result.expired_found}</p>
                    <p className="text-sm text-muted-foreground">หมดอายุ</p>
                  </div>
                  <div className="bg-muted p-3 rounded-lg text-center">
                    <p className="text-2xl font-bold text-green-600">{result.processed_count}</p>
                    <p className="text-sm text-muted-foreground">ประมวลผลแล้ว</p>
                  </div>
                  <div className="bg-muted p-3 rounded-lg text-center">
                    <p className="text-2xl font-bold text-red-600">{result.failed_count}</p>
                    <p className="text-sm text-muted-foreground">ล้มเหลว</p>
                  </div>
                </div>

                {/* Processing Info */}
                <div className="bg-muted p-3 rounded-lg">
                  <div className="grid grid-cols-2 gap-4 text-sm">
                    <div>
                      <strong>ประมวลผลเมื่อ:</strong> {format(new Date(result.processed_at), 'dd MMM yyyy HH:mm:ss', { locale: th })}
                    </div>
                    <div>
                      <strong>โหมด:</strong> {result.dry_run ? 'ทดสอบ' : 'ประมวลผลจริง'}
                    </div>
                  </div>
                </div>

                {/* Processed Bookings */}
                {result.processed_bookings && result.processed_bookings.length > 0 && (
                  <div className="space-y-2">
                    <h5 className="font-medium">รายการที่ประมวลผล:</h5>
                    <div className="border rounded-lg">
                      <Table>
                        <TableHeader>
                          <TableRow>
                            <TableHead>เลขที่การจอง</TableHead>
                            <TableHead>ลูกค้า</TableHead>
                            <TableHead>หมดอายุเมื่อ</TableHead>
                            <TableHead>การดำเนินการ</TableHead>
                            <TableHead>สถานะ</TableHead>
                          </TableRow>
                        </TableHeader>
                        <TableBody>
                          {result.processed_bookings.map((booking) => (
                            <TableRow key={booking.booking_id}>
                              <TableCell className="font-mono text-sm">
                                {booking.invoice_id}
                              </TableCell>
                              <TableCell>{booking.customer_name}</TableCell>
                              <TableCell className="text-sm">
                                {format(new Date(booking.expires_at), 'dd MMM yyyy HH:mm', { locale: th })}
                              </TableCell>
                              <TableCell>{booking.action}</TableCell>
                              <TableCell>
                                <Badge variant={getStatusColor(booking.success)}>
                                  {booking.success ? 'สำเร็จ' : 'ล้มเหลว'}
                                </Badge>
                              </TableCell>
                            </TableRow>
                          ))}
                        </TableBody>
                      </Table>
                    </div>
                  </div>
                )}

                {/* Errors */}
                {result.errors && result.errors.length > 0 && (
                  <div className="space-y-2">
                    <h5 className="font-medium text-destructive flex items-center">
                      <XCircle className="h-4 w-4 mr-2" />
                      ข้อผิดพลาด:
                    </h5>
                    <div className="bg-destructive/10 border border-destructive rounded-lg p-3">
                      {result.errors.map((error, index) => (
                        <p key={index} className="text-sm text-destructive">• {error}</p>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            </>
          )}
        </CardContent>
      </Card>

      {/* Confirmation Dialog */}
      <AlertDialog open={showConfirmDialog} onOpenChange={setShowConfirmDialog}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle className="flex items-center">
              <AlertTriangle className="h-5 w-5 mr-2 text-destructive" />
              ยืนยันการประมวลผล
            </AlertDialogTitle>
            <AlertDialogDescription>
              พบการจองที่หมดอายุ <strong>{result?.expired_found}</strong> รายการ
              <br />
              <br />
              การดำเนินการนี้จะลบการจองที่หมดอายุออกจากระบบ <strong>และไม่สามารถยกเลิกได้</strong>
              <br />
              <br />
              คุณต้องการดำเนินการต่อหรือไม่?
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>ยกเลิก</AlertDialogCancel>
            <AlertDialogAction
              onClick={handleConfirmProcess}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              ยืนยัน ประมวลผลจริง
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}