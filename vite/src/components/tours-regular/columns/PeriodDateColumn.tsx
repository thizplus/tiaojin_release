import { Badge } from '@/components/ui/badge';

type Props = {
  startDate: string;
  endDate: string;
  isNew?: boolean;
}

const GetDateGroup = ({ startDate, endDate, isNew }: Props) => {
  const start = new Date(startDate);
  const end = new Date(endDate);

  const formatDay = (date: Date) => date.getDate();

  const months = [
    "ม.ค.", "ก.พ.", "มี.ค.", "เม.ย.", "พ.ค.", "มิ.ย.",
    "ก.ค.", "ส.ค.", "ก.ย.", "ต.ค.", "พ.ย.", "ธ.ค."
  ];

  const formatMonth = (date: Date) => months[date.getMonth()];

  const formatMonthYear = (date: Date) => {
    return `${months[date.getMonth()]} ${date.getFullYear() + 543}`;
  };

  return (
    <div className="relative flex flex-col items-center text-center min-w-[90px] px-4 py-1 text-foreground rounded-md">
      {isNew && (
        <div className="absolute -top-2 -right-6 z-10">
          <Badge variant="destructive" className="text-xs px-1 py-0 h-5">
            NEW
          </Badge>
        </div>
      )}

      {(start.getMonth() !== end.getMonth() || start.getFullYear() !== end.getFullYear()) ? (
        <>
          <div className="flex items-center space-x-1 text-primary font-bold text-lg">
            <span className='text-foreground text-2xl leading-none'>{formatDay(start)}</span>
            <span className='text-foreground text-xs leading-none'>{formatMonth(start)}</span>
            <span className="text-foreground text-2xl leading-none">-</span>
            <span className='text-foreground text-2xl leading-none'>{formatDay(end)}</span>
            <span className='text-foreground text-xs leading-none'>{formatMonth(end)}</span>
          </div>
          <div className="text-xs text-foreground leading-tight">
            {formatMonthYear(end)}
          </div>
        </>
      ) : (
        <>
          <div className="flex items-center space-x-1 text-primary font-bold text-lg">
            <span className='text-foreground text-2xl leading-none'>{formatDay(start)}</span>
            <span className="text-foreground text-2xl leading-none">-</span>
            <span className='text-foreground text-2xl leading-none'>{formatDay(end)}</span>
          </div>
          <div className="text-xs text-foreground leading-tight">
            {formatMonthYear(start)}
          </div>
        </>
      )}

    </div>
  );
};

interface PeriodDateColumnProps {
  startDate: string;
  endDate: string;
  createdAt?: string;
  status: string;
  // Props สำหรับ mobile
  tourTitle?: string;
  tourThumbnail?: string;
  // Expiration Override (ไม่ใช้แล้ว - แสดงใน Sheet แทน)
  expirationOverride?: boolean;
  expirationDisplay?: string;
}

export function PeriodDateColumn({
  startDate,
  endDate,
  createdAt,
  status,
  tourTitle,
}: PeriodDateColumnProps) {
  const isNew = status === 'publish' && createdAt &&
    new Date(createdAt) > new Date(Date.now() - 7 * 24 * 60 * 60 * 1000);

  return (
    <>
      {/* Desktop: แสดงแค่วันที่ */}
      <div className="hidden md:flex items-center">
        <GetDateGroup
          startDate={startDate}
          endDate={endDate}
          isNew={isNew || false}
        />
      </div>

      {/* Mobile: แสดงภาพปก + ชื่อทัวร์ + วันที่ */}
      <div className="block md:hidden">
        <div className="flex items-start space-x-3 p-2 bg-muted/30 rounded-lg">

          {/* ข้อมูลทัวร์ */}
          <div className="flex-1 min-w-0">
            <div className="flex items-start justify-between">
              <div className="flex-1 min-w-0">
                <h4 className="text-xs font-medium text-foreground line-clamp-3 leading-tight">
                  {tourTitle || 'ชื่อทัวร์'}
                </h4>
              </div>
              {/* วันที่แบบย่อ */}
              <div className="flex-shrink-0 ml-2">
                <GetDateGroup
                  startDate={startDate}
                  endDate={endDate}
                  isNew={isNew || false}
                />
              </div>
            </div>
          </div>
        </div>
      </div>
    </>
  );
}