import { PDFDocument, rgb } from 'pdf-lib';
import fontkit from '@pdf-lib/fontkit';
import type { BookingResponse } from '@/types';

const imageFile = '/images/logo/logo.png';

function getDate(isoDate: string) {
  if (!isoDate) {
    return '-';
  }

  const date = new Date(isoDate);

  const months = [
    "ม.ค.", "ก.พ.", "มี.ค.", "เม.ย.", "พ.ค.", "มิ.ย.",
    "ก.ค.", "ส.ค.", "ก.ย.", "ต.ค.", "พ.ย.", "ธ.ค."
  ];

  const day = date.getDate();
  const month = months[date.getMonth()];
  const year = date.getFullYear();

  return `${day} ${month} ${year}`;
}

function numberFormat(number: number) {
  return new Intl.NumberFormat('en-US').format(number);
}

const thaiNumerals = [
  'ศูนย์', 'หนึ่ง', 'สอง', 'สาม', 'สี่',
  'ห้า', 'หก', 'เจ็ด', 'แปด', 'เก้า'
];

const thaiUnits = [
  '', 'สิบ', 'ร้อย', 'พัน', 'หมื่น', 'แสน', 'ล้าน'
];

function numberToThaiText(number: number) {
  if (number === 0) return 'ศูนย์บาทถ้วน';

  const numStr = number.toString();
  let result = '';
  const len = numStr.length;

  for (let i = 0; i < len; i++) {
    const digit = parseInt(numStr[i]);
    const position = len - i - 1;

    if (digit !== 0) {
      if (position === 1 && digit === 1 && len > 1) {
        result += 'สิบ';
      } else if (position === 1 && digit === 2) {
        result += 'ยี่สิบ';
      } else {
        result += thaiNumerals[digit];
        if (position > 0) {
          result += thaiUnits[position];
        }
      }
    }
  }

  return result + 'บาทถ้วน';
}


export async function generateBookingPDF(booking: BookingResponse): Promise<void> {
  try {
    const pdfDoc = await PDFDocument.create();

    // Register fontkit instance
    pdfDoc.registerFontkit(fontkit);

    // Load the Thai font
    const fontNormal = await fetch('/fonts/CSPraJad.otf').then(res => res.arrayBuffer());
    const fontBold = await fetch('/fonts/CSPraJad-bold.otf').then(res => res.arrayBuffer());

    // Embed the font into the document
    const normal = await pdfDoc.embedFont(fontNormal);
    const bold = await pdfDoc.embedFont(fontBold);

    // Add a new page to the document
    const page = pdfDoc.addPage();
    const { width, height } = page.getSize();
    const lg = 16;
    const base = 11;

    // Skip image loading for now to avoid errors
    let image = null;
    try {
      const imageBytes = await fetch(imageFile).then(res => res.arrayBuffer());
      image = await pdfDoc.embedPng(imageBytes);
    } catch (error) {
      console.warn('Could not load logo image:', error);
    }

    // Draw the image on the page if loaded successfully
    if (image) {
      page.drawImage(image, {
        x: 15,
        y: height - 80,
        width: 70,
        height: 70,
      });
    }

    // Company header
    page.drawText('BOOKING CHINA (บุ๊คกิ้ง ไช่น่า)', {
      x: 90,
      y: height - 35,
      size: lg,
      font: bold,
      color: rgb(0, 0, 0)
    });

    page.drawText('บริษัท บุ๊คกิ้ง ไช่น่า จำกัด', {
      x: 90,
      y: height - 55,
      size: base,
      font: normal,
      color: rgb(0, 0, 0)
    });

    page.drawText('เลขประจำตัวผู้เสียภาษี: 0107569000012', {
      x: 90,
      y: height - 70,
      size: base,
      font: normal,
      color: rgb(0, 0, 0)
    });

    // Invoice header
    page.drawText('INVOICE', {
      x: width - 150,
      y: height - 35,
      size: 24,
      font: bold,
      color: rgb(0, 0, 0)
    });

    page.drawText(booking.invoice_id, {
      x: width - 150,
      y: height - 60,
      size: base,
      font: normal,
      color: rgb(0, 0, 0)
    });

    // Date
    page.drawText(`วันที่: ${getDate(booking.created_at)}`, {
      x: width - 150,
      y: height - 80,
      size: base,
      font: normal,
      color: rgb(0, 0, 0)
    });

    // Customer info
    page.drawText('ลูกค้า:', {
      x: 20,
      y: height - 120,
      size: base,
      font: bold,
      color: rgb(0, 0, 0)
    });

    page.drawText(booking.customer_name || '', {
      x: 60,
      y: height - 120,
      size: base,
      font: normal,
      color: rgb(0, 0, 0)
    });

    page.drawText('เบอร์โทร:', {
      x: 20,
      y: height - 140,
      size: base,
      font: normal,
      color: rgb(0, 0, 0)
    });

    page.drawText(booking.customer_phone || '', {
      x: 70,
      y: height - 140,
      size: base,
      font: normal,
      color: rgb(0, 0, 0)
    });

    // Tour info
    page.drawText('ทัวร์:', {
      x: 20,
      y: height - 170,
      size: base,
      font: bold,
      color: rgb(0, 0, 0)
    });

    page.drawText(booking.tour_details?.title || '', {
      x: 55,
      y: height - 170,
      size: base,
      font: normal,
      color: rgb(0, 0, 0)
    });

    page.drawText('รหัสทัวร์:', {
      x: 20,
      y: height - 190,
      size: base,
      font: normal,
      color: rgb(0, 0, 0)
    });

    page.drawText(booking.tour_details?.code || '', {
      x: 75,
      y: height - 190,
      size: base,
      font: normal,
      color: rgb(0, 0, 0)
    });

    // Period info
    if (booking.period_details) {
      page.drawText('วันเดินทาง:', {
        x: 20,
        y: height - 210,
        size: base,
        font: normal,
        color: rgb(0, 0, 0)
      });

      page.drawText(`${getDate(booking.period_details.start_date)} - ${getDate(booking.period_details.end_date)}`, {
        x: 80,
        y: height - 210,
        size: base,
        font: normal,
        color: rgb(0, 0, 0)
      });
    }

    // Table header
    page.drawText('รายการ', {
      x: 30,
      y: height - 260,
      size: base,
      font: bold
    });
    page.drawText('จำนวน', {
      x: 200,
      y: height - 260,
      size: base,
      font: bold
    });
    page.drawText('ราคา', {
      x: 270,
      y: height - 260,
      size: base,
      font: bold
    });
    page.drawText('รวม', {
      x: 350,
      y: height - 260,
      size: base,
      font: bold
    });

    // Table line
    page.drawLine({
      start: { x: 20, y: height - 275 },
      end: { x: width - 20, y: height - 275 },
      thickness: 1,
      color: rgb(0, 0, 0),
      opacity: 0.8
    });

    // Create table items from booking data
    const tableItems: Array<{
      id: number;
      title: string;
      qty: number;
      price: number;
      total: number;
    }> = [];
    let itemIndex = 1;

    // Add rate types
    if (booking.rate_types) {
      booking.rate_types.forEach(rate => {
        if (rate.price > 0) {
          tableItems.push({
            id: itemIndex++,
            title: rate.title,
            qty: rate.quantity,
            price: rate.price,
            total: rate.total
          });
        }
      });
    }

    // Add room types
    if (booking.room_types) {
      booking.room_types.forEach(room => {
        if (room.quantity > 0 && room.price > 0) {
          tableItems.push({
            id: itemIndex++,
            title: room.title,
            qty: room.quantity,
            price: room.price,
            total: room.total
          });
        }
      });
    }

    // Add additional services
    if (booking.additional_services) {
      booking.additional_services.forEach(service => {
        if (service.price > 0) {
          tableItems.push({
            id: itemIndex++,
            title: service.title,
            qty: service.quantity,
            price: service.price,
            total: service.total
          });
        }
      });
    }

    // Draw table items
    tableItems.forEach((item, index) => {
      const yPos = height - 290 - (index * 15);

      page.drawText(item.title, {
        x: 30,
        y: yPos,
        size: 10,
        font: normal
      });

      page.drawText(item.qty.toString(), {
        x: 200,
        y: yPos,
        size: 10,
        font: normal
      });

      page.drawText(numberFormat(item.price), {
        x: 270,
        y: yPos,
        size: 10,
        font: normal
      });

      page.drawText(numberFormat(item.total), {
        x: 350,
        y: yPos,
        size: 10,
        font: normal
      });
    });

    // Summary
    const summaryY = height - 290 - (tableItems.length * 15) - 30;

    page.drawText('รวม:', {
      x: 300,
      y: summaryY,
      size: base,
      font: bold
    });

    page.drawText(numberFormat(booking.sum_net), {
      x: 350,
      y: summaryY,
      size: base,
      font: bold
    });

    // Room details
    page.drawText('ห้องพักที่เลือก :', {
      x: 20,
      y: height - 535,
      size: 10,
      font: normal,
      opacity: 0.5
    });

    const roomDetails = booking.room_types
      ? booking.room_types
          .filter(room => room.quantity > 0)
          .map(room => `${room.title} ${room.quantity} ห้อง`)
          .join(', ')
      : '-';

    page.drawText(roomDetails, {
      x: 100,
      y: height - 535,
      size: 10,
      font: normal,
      opacity: 0.5
    });

    // Amount in words
    page.drawText('จำนวนเงินตัวอักษร', {
      x: 30,
      y: height - height + 280,
      size: base,
      font: normal
    });

    page.drawText(numberToThaiText(booking.sum_net), {
      x: 30,
      y: height - height + 260,
      size: base,
      font: normal
    });

    // Save and download
    const pdfBytes = await pdfDoc.save();
    const blob = new Blob([new Uint8Array(pdfBytes)], { type: 'application/pdf' });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = `invoice-${booking.invoice_id}.pdf`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(url);

  } catch (error) {
    console.error('Error generating PDF:', error);
    throw error;
  }
}