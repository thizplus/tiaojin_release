import { format } from "date-fns";
import { th } from "date-fns/locale";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Badge } from "@/components/ui/badge";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  Calendar,
  Users,
  Plane,
  MapPin,
  DollarSign,
  Percent,
  Clock,
  FileText,
  Bed,
  Plus,
  Flame,
  Info
} from "lucide-react";
import type { TourPeriod } from "@/types";
import { getPeriodDialogStatusIcon, getPeriodDialogStatusLabel } from "@/utils/status-helpers";

interface PeriodDetailDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  period: TourPeriod | null;
}

export function PeriodDetailDialog({
  open,
  onOpenChange,
  period,
}: PeriodDetailDialogProps) {
  if (!period) return null;

  const startDate = new Date(period.start_date);
  const endDate = new Date(period.end_date);
  const diffTime = Math.abs(endDate.getTime() - startDate.getTime());
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24)) + 1;

  const formatDeposit = () => {
    if (period.deposit_type === 'percent') {
      return `${period.deposit}%`;
    }
    return `฿${period.deposit.toLocaleString()}`;
  };

  const getDepositLabel = () => {
    switch (period.deposit_type) {
      case 'per_unit': return 'ต่อคน';
      case 'percent': return 'เปอร์เซ็นต์';
      default: return period.deposit_type;
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case "publish":
        return "bg-green-100 text-green-800 border-green-200";
      case "draft":
        return "bg-gray-100 text-gray-800 border-gray-200";
      default:
        return "bg-gray-100 text-gray-800 border-gray-200";
    }
  };

  const getRateTypeLabel = (rateType: string) => {
    switch (rateType) {
      case 'per_unit': return 'ต่อหน่วย';
      case 'per_booking': return 'ต่อการจอง';
      default: return rateType;
    }
  };

  // Helper functions สำหรับราคาและคอมมิชชั่น
  const getCurrentPrice = (rateItem: (typeof period.rate_items)[0]) => {
    return period.rate_type === 'sale' ? rateItem.sale : rateItem.regular;
  };

  const getCurrentCommission = (commissionItem: (typeof period.commission_items)[0]) => {
    return period.commission_type === 'sale' ? commissionItem.sale : commissionItem.regular;
  };

  const getPriceTypeLabel = () => {
    return period.rate_type === 'sale' ? 'ราคาพิเศษ' : 'ราคาปกติ';
  };

  const getCommissionTypeLabel = () => {
    return period.commission_type === 'sale' ? 'คอมพิเศษ' : 'คอมปกติ';
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px] max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="flex items-center space-x-3">
            <Calendar className="h-6 w-6" />
            <span>รายละเอียดช่วงเวลาทัวร์</span>
            <Badge variant="outline" className={getStatusColor(period.status)}>
              {getPeriodDialogStatusIcon(period.status)}
              <span className="ml-1">{getPeriodDialogStatusLabel(period.status)}</span>
            </Badge>
          </DialogTitle>
        </DialogHeader>

        <div className="space-y-6">
          {/* Basic Information */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <Calendar className="h-5 w-5" />
                <span>ข้อมูลพื้นฐาน</span>
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                <div className="space-y-2">
                  <div className="text-sm text-muted-foreground">ระยะเวลา</div>
                  <div className="font-medium">
                    {format(startDate, "dd MMM yyyy", { locale: th })} - {format(endDate, "dd MMM yyyy", { locale: th })}
                  </div>
                  <div className="text-sm text-muted-foreground flex items-center">
                    <Clock className="h-4 w-4 mr-1" />
                    {diffDays} วัน
                  </div>
                </div>

                <div className="space-y-2">
                  <div className="text-sm text-muted-foreground">สายการบิน</div>
                  <div className="flex items-center space-x-2">
                    <Plane className="h-4 w-4 text-muted-foreground" />
                    <span className="font-medium">{period.airline}</span>
                  </div>
                  <div className="flex items-center space-x-2">
                    <MapPin className="h-4 w-4 text-muted-foreground" />
                    <span className="text-sm text-muted-foreground">{period.airport}</span>
                  </div>
                </div>

                <div className="space-y-2">
                  <div className="text-sm text-muted-foreground">มัดจำ</div>
                  <div className="font-medium text-lg">{formatDeposit()}</div>
                  <div className="text-sm text-muted-foreground">{getDepositLabel()}</div>
                </div>
              </div>

              {period.note && (
                <>
                  <Separator className="my-4" />
                  <div className="space-y-2">
                    <div className="text-sm text-muted-foreground flex items-center">
                      <FileText className="h-4 w-4 mr-1" />
                      หมายเหตุ
                    </div>
                    <div className="text-sm bg-muted p-3 rounded-md">{period.note}</div>
                  </div>
                </>
              )}
            </CardContent>
          </Card>

          {/* Expiration Override - ทัวร์ไฟไหม้ */}
          {period.expiration_override && (
            <Card className="border-orange-500/30 bg-orange-500/5">
              <CardHeader>
                <CardTitle className="flex items-center space-x-2 text-orange-600">
                  <Flame className="h-5 w-5" />
                  <span>การหมดอายุ Invoice (ทัวร์ไฟไหม้)</span>
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  <div className="flex items-center space-x-3">
                    <div className="p-3 bg-orange-500/10 rounded-lg">
                      <Clock className="h-8 w-8 text-orange-500" />
                    </div>
                    <div>
                      <div className="text-2xl font-bold text-orange-600">
                        {period.expiration_display || 'กำหนดเอง'}
                      </div>
                      <div className="text-sm text-orange-600/80">
                        ระยะเวลาหมดอายุ Invoice
                      </div>
                    </div>
                  </div>

                  {period.expiration_message && (
                    <div className="p-3 bg-orange-500/10 rounded-md">
                      <div className="text-sm text-orange-600">
                        <span className="font-medium">ข้อความ:</span> {period.expiration_message}
                      </div>
                    </div>
                  )}

                  <div className="text-xs text-orange-600/70 flex items-center space-x-1">
                    <Info className="h-3 w-3" />
                    <span>Invoice จะหมดอายุเร็วกว่าปกติตามที่กำหนด</span>
                  </div>
                </div>
              </CardContent>
            </Card>
          )}

          {/* Seats Information */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <Users className="h-5 w-5" />
                <span>ข้อมูลที่นั่ง</span>
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                <div className="text-center p-3 bg-muted rounded-lg">
                  <div className="text-2xl font-bold text-muted-foreground">{period.seat}</div>
                  <div className="text-sm text-muted-foreground">ทั้งหมด</div>
                </div>
                <div className="text-center p-3 bg-muted rounded-lg">
                  <div className="text-2xl font-bold text-chart-2">{period.booking}</div>
                  <div className="text-sm text-chart-2">จองแล้ว</div>
                </div>
                <div className="text-center p-3 bg-muted rounded-lg">
                  <div className="text-2xl font-bold text-muted-foreground">{period.held_seats}</div>
                  <div className="text-sm text-muted-foreground">จำนอง</div>
                </div>
                <div className="text-center p-3 bg-muted rounded-lg">
                  <div className="text-2xl font-bold text-muted-foreground">{period.available_seats}</div>
                  <div className="text-sm text-muted-foreground">ว่าง</div>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Flight Information */}
          {(period.sector_1?.length > 0 || period.sector_2?.length > 0) && (
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center space-x-2">
                  <Plane className="h-5 w-5" />
                  <span>ข้อมูลเที่ยวบิน</span>
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div>
                    <h4 className="font-medium text-muted-foreground mb-3">เที่ยวบินไป</h4>
                    <div className="space-y-2">
                      {period.sector_1?.length > 0 ? (
                        period.sector_1.map((flight, index) => (
                          <div key={index} className="p-3 border rounded-md">
                            <div className="font-medium">{flight.flight || '-'}</div>
                            <div className="text-sm text-muted-foreground">{flight.route || '-'}</div>
                            <div className="text-sm text-muted-foreground">
                              {flight.date && !isNaN(new Date(flight.date).getTime())
                                ? format(new Date(flight.date), "dd MMM yyyy", { locale: th })
                                : '-'}
                            </div>
                          </div>
                        ))
                      ) : (
                        <div className="text-sm text-muted-foreground">ไม่มีข้อมูล</div>
                      )}
                    </div>
                  </div>

                  <div>
                    <h4 className="font-medium text-muted-foreground mb-3">เที่ยวบินกลับ</h4>
                    <div className="space-y-2">
                      {period.sector_2?.length > 0 ? (
                        period.sector_2.map((flight, index) => (
                          <div key={index} className="p-3 border rounded-md">
                            <div className="font-medium">{flight.flight || '-'}</div>
                            <div className="text-sm text-muted-foreground">{flight.route || '-'}</div>
                            <div className="text-sm text-muted-foreground">
                              {flight.date && !isNaN(new Date(flight.date).getTime())
                                ? format(new Date(flight.date), "dd MMM yyyy", { locale: th })
                                : '-'}
                            </div>
                          </div>
                        ))
                      ) : (
                        <div className="text-sm text-muted-foreground">ไม่มีข้อมูล</div>
                      )}
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          )}

          {/* Pricing and Commission Information with Tabs */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <DollarSign className="h-5 w-5" />
                <span>ข้อมูลราคาและคอมมิชชั่น</span>
              </CardTitle>
            </CardHeader>
            <CardContent>
              <Tabs defaultValue="pricing" className="w-full">
                <TabsList className="grid w-full grid-cols-2">
                  <TabsTrigger value="pricing" className="flex items-center space-x-2">
                    <DollarSign className="h-4 w-4" />
                    <span>ราคา</span>
                    {period.rate_type === 'sale' ? (
                      <Badge variant="outline" className="text-xs text-primary border-primary flex items-center space-x-1">
                        <Flame className="h-3 w-3" />
                        <span>{getPriceTypeLabel()}</span>
                      </Badge>
                    ) : (
                      <Badge variant="outline" className="text-xs">
                        {getPriceTypeLabel()}
                      </Badge>
                    )}
                  </TabsTrigger>
                  <TabsTrigger value="commission" className="flex items-center space-x-2">
                    <Percent className="h-4 w-4" />
                    <span>คอมมิชชั่น</span>
                    {period.commission_type === 'sale' ? (
                      <Badge variant="outline" className="text-xs text-primary border-primary flex items-center space-x-1">
                        <Flame className="h-3 w-3" />
                        <span>{getCommissionTypeLabel()}</span>
                      </Badge>
                    ) : (
                      <Badge variant="outline" className="text-xs">
                        {getCommissionTypeLabel()}
                      </Badge>
                    )}
                  </TabsTrigger>
                </TabsList>

                <TabsContent value="pricing" className="mt-6">
                  <div className="space-y-3">
                    {period.rate_items.map((rate) => (
                      <div key={rate.id} className="flex justify-between items-center p-3 border rounded-md">
                        <div className="flex items-center space-x-3">
                          <div>
                            <div className="font-medium">{rate.title}</div>
                            {rate.exclude && (
                              <Badge variant="outline" className="text-xs mt-1 bg-muted text-muted-foreground">
                                ไม่ได้ส่วนลด
                              </Badge>
                            )}
                          </div>
                        </div>
                        <div className="text-right">
                          <div className="text-lg font-bold text-foreground">
                            ฿{getCurrentPrice(rate).toLocaleString()}
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>

                  <div className="mt-4 p-3 bg-muted/30 rounded-md">
                    <div className="flex items-center space-x-2 text-sm text-muted-foreground">
                      <Info className="h-4 w-4" />
                      <span>แสดงราคา{getPriceTypeLabel()}ทั้งหมด</span>
                    </div>
                  </div>
                </TabsContent>

                <TabsContent value="commission" className="mt-6">
                  <div className="space-y-3">
                    {period.commission_items.map((commission, index) => (
                      <div key={index} className="flex justify-between items-center p-3 border rounded-md">
                        <div className="font-medium">ค่าคอมมิชชั่น {index + 1}</div>
                        <div className="text-right">
                          <div className="text-lg font-bold text-foreground">
                            ฿{getCurrentCommission(commission).toLocaleString()}
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>

                  <div className="mt-4 p-3 bg-muted/30 rounded-md">
                    <div className="flex items-center space-x-2 text-sm text-muted-foreground">
                      <Info className="h-4 w-4" />
                      <span>แสดง{getCommissionTypeLabel()}ทั้งหมด</span>
                    </div>
                  </div>
                </TabsContent>
              </Tabs>
            </CardContent>
          </Card>

          {/* Room Types */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <Bed className="h-5 w-5" />
                <span>ประเภทห้องพัก</span>
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-3">
                {period.room_types.map((room) => (
                  <div key={room.id} className="flex justify-between items-center p-3 border rounded-md">
                    <div className="font-medium">{room.title}</div>
                    <div className="text-right">
                      {room.price === 0 ? (
                        <span className="text-foreground font-medium">ฟรี</span>
                      ) : room.price > 0 ? (
                        <span className="text-foreground font-medium">
                          +฿{room.price.toLocaleString()}
                        </span>
                      ) : (
                        <span className="text-foreground font-medium">
                          ฿{room.price.toLocaleString()}
                        </span>
                      )}
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>

          {/* Additional Services */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <Plus className="h-5 w-5" />
                <span>บริการเสริม</span>
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-3">
                {period.additional_services.map((service) => (
                  <div key={service.id} className="flex justify-between items-center p-3 border rounded-md">
                    <div>
                      <div className="font-medium">{service.title}</div>
                      <div className="text-sm text-muted-foreground">
                        {getRateTypeLabel(service.rate_type)} • {service.unit_name}
                      </div>
                    </div>
                    <div className="text-right">
                      <div className="text-lg font-bold text-foreground">
                        ฿{service.price.toLocaleString()}
                      </div>
                      <div className="text-xs text-muted-foreground">
                        ต่อ{service.unit_name}
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>

          {/* Additional Information */}
          <Card>
            <CardHeader>
              <CardTitle>ข้อมูลเพิ่มเติม</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <div className="text-sm text-muted-foreground">วันที่สร้าง</div>
                  <div className="font-medium">
                    {format(new Date(period.created_at), "dd MMM yyyy HH:mm", { locale: th })}
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </DialogContent>
    </Dialog>
  );
}