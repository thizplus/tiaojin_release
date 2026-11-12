import { format } from 'date-fns';
import { th } from 'date-fns/locale';
import type { BookingResponse } from '@/types';

interface ExportBookingsOptions {
  tourTitle?: string;
  periodStartDate?: string;
  periodEndDate?: string;
}

export function exportBookingsToCSV(
  bookings: BookingResponse[],
  options: ExportBookingsOptions = {}
): void {
  if (!bookings || bookings.length === 0) {
    alert('ไม่มีข้อมูลสำหรับการส่งออก');
    return;
  }

  // สร้าง CSV header
  const headers = [
    'เลขที่ใบจอง',
    'วันที่',
    'โปรแกรม',
    'period date',
    'ชื่อบริษัท',
    'ชื่อลูกค้า',
    'เบอร์โทรศัพท์ลูกค้า',
    'ที่นั่ง',
    'ประเภท',
    'ห้องพัก',
    'รวมราคา',
    'ชำระแล้ว',
    'ค้างชำระ',
    'สถานะ',
    'เลขที่หนังสือเดินทาง',
    'คำนำหน้า(ไทย)',
    'ชื่อ(ไทย)',
    'นามสกุล(ไทย)',
    'คำนำหน้า(อังกฤษ)',
    'ชื่อ(อังกฤษ)',
    'นามสกุล(อังกฤษ)',
    'เพศ',
    'สัญชาติ',
    'วันเกิด',
    'เลขประจำตัวประชาชน',
    'สถานที่เกิด',
    'วันหมดอายุ'
  ];

  // สร้างข้อมูล CSV
  const rows = bookings.flatMap((booking) => {
    // ข้อมูลพื้นฐาน
    const createdDate = format(new Date(booking.created_at), 'd MMM yyyy / HH:mm น.', { locale: th });
    const tourTitle = booking.tour_details?.title || '';

    let periodDate = '';
    if (booking.period_details) {
      const startDate = format(new Date(booking.period_details.start_date), 'd MMM yyyy', { locale: th });
      const endDate = format(new Date(booking.period_details.end_date), 'd MMM yyyy', { locale: th });
      periodDate = `${startDate} - ${endDate}`;
    }

    const companyName = booking.company_profile?.name || '';
    const customerName = booking.customer_name || '';
    const customerPhone = booking.customer_phone || '';

    // คำนวณที่นั่งรวม
    const totalSeats = booking.sum_qty || 0;

    // สร้างข้อมูลประเภท (rate types)
    const rateTypesText = booking.rate_types
      ?.map(rt => `${rt.title} (${rt.quantity})`)
      .join(', ') || '';

    // สร้างข้อมูลห้องพัก
    const roomTypesText = booking.room_types
      ?.map(room => `${room.title} (${room.quantity})`)
      .join(', ') || '';

    // ข้อมูลการเงิน
    const totalPrice = booking.sum_rate || 0;
    const paidAmount = booking.sum_payment || 0;
    const remainingAmount = totalPrice - paidAmount;

    // สถานะ
    const statusMap: Record<string, string> = {
      'pending': 'รอดำเนินการ',
      'processing': 'กำลังดำเนินการ',
      'confirm': 'ยืนยันแล้ว',
      'completed': 'เสร็จสิ้น',
      'cancel': 'ยกเลิก',
      'expired': 'หมดอายุ'
    };
    const status = statusMap[booking.status] || booking.status;

    // ถ้ามี passport data ให้สร้างแถวต่อหนัง
    if (booking.passports && booking.passports.length > 0) {
      return booking.passports.map((passport, index) => {
        // แถวแรกใส่ข้อมูลทั่วไป, แถวถัดไปเว้นว่าง (หรือซ้ำข้อมูล)
        const isFirstRow = index === 0;

        return [
          isFirstRow ? (booking.invoice_id || '') : '',
          isFirstRow ? createdDate : '',
          isFirstRow ? tourTitle : '',
          isFirstRow ? periodDate : '',
          isFirstRow ? companyName : '',
          isFirstRow ? customerName : '',
          isFirstRow ? customerPhone : '',
          isFirstRow ? totalSeats : '',
          isFirstRow ? rateTypesText : '',
          isFirstRow ? roomTypesText : '',
          isFirstRow ? totalPrice : '',
          isFirstRow ? paidAmount : '',
          isFirstRow ? remainingAmount : '',
          isFirstRow ? status : '',
          passport.passport_number || '',
          passport.title_th || '',
          passport.name_th || '',
          passport.surname_th || '',
          passport.title_en || '',
          passport.name_en || '',
          passport.surname_en || '',
          passport.gender || '',
          passport.nationality || '',
          passport.date_of_birth ? format(new Date(passport.date_of_birth), 'd MMM yyyy', { locale: th }) : '',
          passport.id_number || '',
          passport.place_of_birth || '',
          passport.date_of_expiry ? format(new Date(passport.date_of_expiry), 'd MMM yyyy', { locale: th }) : ''
        ];
      });
    } else {
      // ไม่มี passport ก็แสดงแค่ข้อมูลการจอง
      return [[
        booking.invoice_id || '',
        createdDate,
        tourTitle,
        periodDate,
        companyName,
        customerName,
        customerPhone,
        totalSeats,
        rateTypesText,
        roomTypesText,
        totalPrice,
        paidAmount,
        remainingAmount,
        status,
        '', '', '', '', '', '', '', '', '', '', '', '', ''
      ]];
    }
  });

  // แปลงเป็น CSV string
  const csvContent = [
    '\ufeff', // BOM สำหรับ UTF-8
    headers.map(escapeCSVValue).join(','),
    ...rows.map(row => row.map(escapeCSVValue).join(','))
  ].join('\n');

  // สร้างชื่อไฟล์
  let filename = 'bookings';
  if (options.tourTitle) {
    filename = options.tourTitle;
  }
  if (options.periodStartDate && options.periodEndDate) {
    const startDate = format(new Date(options.periodStartDate), 'd-M-yyyy');
    const endDate = format(new Date(options.periodEndDate), 'd-M-yyyy');
    filename += ` ${startDate} ถึง ${endDate}`;
  }
  filename += '.csv';

  // ดาวน์โหลดไฟล์
  downloadCSV(csvContent, filename);
}

function escapeCSVValue(value: string | number | null | undefined): string {
  if (value === null || value === undefined) {
    return '""';
  }

  const stringValue = String(value);

  // ถ้ามี comma, quotes, หรือ newlines ต้องใส่ quotes
  if (stringValue.includes(',') || stringValue.includes('"') || stringValue.includes('\n')) {
    // Escape quotes โดยการใส่ double quotes
    return `"${stringValue.replace(/"/g, '""')}"`;
  }

  return `"${stringValue}"`;
}

function downloadCSV(csvContent: string, filename: string): void {
  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
  const link = document.createElement('a');

  if (link.download !== undefined) {
    const url = URL.createObjectURL(blob);
    link.setAttribute('href', url);
    link.setAttribute('download', filename);
    link.style.visibility = 'hidden';
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(url);
  }
}