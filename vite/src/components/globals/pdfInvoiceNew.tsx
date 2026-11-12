import { Button } from '@/components/ui/button';
import { PDFDocument, rgb } from 'pdf-lib';
import fontkit from '@pdf-lib/fontkit'
import { FileText } from 'lucide-react';
import type { BookingResponse } from '@/types';

type PropsPdf = {
  booking: BookingResponse;
  variant?: "default" | "destructive" | "outline" | "secondary" | "ghost" | "link";
  size?: "default" | "sm" | "lg" | "icon";
  style?: "row" | "column";
}

const imageFile = '/images/logo/logo.png';

const PdfInvoiceNew = ({ booking, variant = "outline", size = "sm", style = "row" }: PropsPdf) => {

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
    if (number === 0) {
      return `${thaiNumerals[0]}บาทถ้วน`;
    }

    const parts = number.toFixed(2).split('.');
    const integerPart = parts[0];
    const decimalPart = parts[1];

    const integerText = convertIntegerPartToThaiText(integerPart);
    const decimalText = convertDecimalPartToThaiText(decimalPart);

    if (parseInt(decimalPart, 10) === 0) {
      return `${integerText}บาทถ้วน`;
    }

    return `${integerText}บาท${decimalText}สตางค์`;
  }

  function convertIntegerPartToThaiText(number: string) {
    let text = '';
    const digits = number.split('').reverse();

    digits.forEach((digit: string, index: number) => {
      const digitValue = parseInt(digit, 10);

      if (digitValue === 0) {
        return;
      }

      if (index === 1 && digitValue === 1) {
        text = `สิบ${text}`;
      } else if (index === 1 && digitValue === 2) {
        text = `ยี่สิบ${text}`;
      } else if (index === 0 && digitValue === 1 && digits.length > 1) {
        text = `เอ็ด${text}`;
      } else {
        text = `${thaiNumerals[digitValue]}${thaiUnits[index]}${text}`;
      }
    });

    return text;
  }

  function convertDecimalPartToThaiText(number: string) {
    return number.split('').map((digit: string) => thaiNumerals[parseInt(digit, 10)]).join('');
  }

  function getDateDiff(startDate: string, createDate: string) {
    const start = new Date(startDate);
    const create = new Date(createDate);
    const diffTime = Math.abs(start.getTime() - create.getTime());
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
    return diffDays;
  }

  async function createPdf() {
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
        y: height - 25,
        size: lg,
        font: bold
      });

      page.drawText('บริษัท เที่ยวจีน จำกัด', {
        x: 90,
        y: height - 40,
        size: 13,
        font: bold
      });

      page.drawText('เลขที่ 120/28 ถนน จันทร์ทองเอี่ยม ตำบล บางรักพัฒนา อำเภอบางบัวทอง นนทบุรี 11110', {
        x: 90,
        y: height - 52,
        size: base,
        font: normal
      });

      page.drawText('เลขที่ผู้เสียภาษีอากร : 0125566041751', {
        x: 90,
        y: height - 64,
        size: base,
        font: normal
      });

      page.drawText('โทรศัพท์ : 090-969-2436', {
        x: 90,
        y: height - 76,
        size: base,
        font: normal
      });

      // Red line
      page.drawRectangle({
        x: width - 20,
        y: height - 80,
        width: 4,
        height: 70,
        color: rgb(1, 0, 0),
        borderWidth: 0,
      });

      // Invoice title
      page.drawText('ใบแจ้งหนี้', {
        x: width - 85,
        y: height - 35,
        size: lg,
        font: bold
      });

      page.drawText('INVOICE', {
        x: width - 67,
        y: height - 47,
        size: base,
        font: normal
      });

      page.drawText('สำหรับลูกค้า', {
        x: width - 76,
        y: height - 60,
        size: base,
        font: normal
      });

      // Customer info labels (left side)
      const leftLabels = [
        { text: 'รหัสลูกค้า', y: height - 110 },
        { text: 'บริษัท', y: height - 125 },
        { text: 'เลขผู้เสียภาษี', y: height - 140 },
        { text: 'ที่อยู่', y: height - 155 },
        { text: 'โทรศัพท์', y: height - 180 },
        { text: 'อีเมล', y: height - 195 }
      ];

      leftLabels.forEach(label => {
        page.drawText(label.text, {
          x: 30,
          y: label.y,
          size: base,
          font: normal
        });
      });

      // Invoice info labels (right side)
      const rightLabels = [
        { text: 'เลขที่เอกสาร', y: height - 110 },
        { text: 'วันที่', y: height - 125 },
        { text: 'ผู้จอง (Sale)', y: height - 140 },
        { text: 'พนักงานขาย', y: height - 155 },
        { text: 'ข้อมูลผู้เดินทาง', y: height - 180 },
        { text: 'โทรศัพท์', y: height - 195 }
      ];

      rightLabels.forEach(label => {
        page.drawText(label.text, {
          x: 380,
          y: label.y,
          size: base,
          font: normal
        });
      });

      // Customer data (left side)
      page.drawText(booking.company_profile?.id || '-', {
        x: 110,
        y: height - 110,
        size: base,
        font: normal
      });

      page.drawText(booking.company_profile?.name || '-', {
        x: 110,
        y: height - 125,
        size: base,
        font: normal
      });

      page.drawText(booking.company_profile?.tax_id || '-', { // No tax ID in new system
        x: 110,
        y: height - 140,
        size: base,
        font: normal
      });

      page.drawText(booking.company_profile?.address + ' ' + booking.company_profile?.province + ' ' + booking.company_profile?.zipcode || '-', { // No address in new system
        x: 110,
        y: height - 155,
        size: base,
        font: normal
      });

      page.drawText(booking.company_profile?.phone || '-', {
        x: 110,
        y: height - 180,
        size: base,
        font: normal
      });

      page.drawText(booking.user_details?.email || '-', {
        x: 110,
        y: height - 195,
        size: base,
        font: normal
      });

      // Invoice data (right side)
      page.drawText(booking.invoice_id, {
        x: 450,
        y: height - 110,
        size: base,
        font: normal
      });

      page.drawText(getDate(booking.created_at), {
        x: 450,
        y: height - 125,
        size: base,
        font: normal
      });

      page.drawText(booking.user_details?.email || '-', {
        x: 450,
        y: height - 140,
        size: base,
        font: normal
      });

      page.drawText('[' +booking.sales_details?.code + '] ' + booking.sales_details?.first_name + ' ' + booking.sales_details?.last_name || '-', { // No sales person in new system
        x: 450,
        y: height - 155,
        size: base,
        font: normal
      });

      page.drawText(booking.customer_name, {
        x: 450,
        y: height - 180,
        size: base,
        font: normal
      });

      page.drawText(booking.customer_phone, {
        x: 450,
        y: height - 195,
        size: base,
        font: normal
      });

      // Tour info box
      page.drawRectangle({
        x: 20,
        y: height - 210,
        width: 559,
        height: 120,
        borderColor: rgb(0, 0, 0),
        borderWidth: 1,
        borderOpacity: 0.1,
        color: undefined,
      });

      // Tour details
      const tourTitle = booking.tour_details
        ? `[${booking.tour_details.code}] ${booking.tour_details.title}`
        : 'ไม่มีข้อมูลทัวร์';

      page.drawText(tourTitle, {
        x: 20,
        y: height - 230,
        size: lg,
        font: bold,
        lineHeight: 15,
        maxWidth: 500,
        wordBreaks: [" "]
      });

      // Date info
      if (booking.period_details) {
        page.drawText('เช็คอิน', {
          x: 20,
          y: height - 265,
          size: base,
          font: normal
        });

        page.drawText(getDate(booking.period_details.start_date), {
          x: 60,
          y: height - 265,
          size: base,
          font: normal
        });

        page.drawText('เช็คเอาท์', {
          x: 200,
          y: height - 265,
          size: base,
          font: normal
        });

        page.drawText(getDate(booking.period_details.end_date), {
          x: 255,
          y: height - 265,
          size: base,
          font: normal
        });
      }

      // Table header
      page.drawRectangle({
        x: 20,
        y: height - 315,
        width: 559,
        height: 40,
        borderColor: rgb(0, 0, 0),
        borderWidth: 1,
        borderOpacity: 0.1,
        color: undefined,
      });

      // Table headers
      const tableHeaders = [
        { text: 'ลำดับ', subText: 'Item', x: 30, align: 'left' as const },
        { text: 'รายการ', subText: 'Description', x: 90, align: 'left' as const },
        { text: 'จำนวน', subText: 'Quantity', x: 390, align: 'right' as const },
        { text: 'ราคาต่อหน่วย', subText: 'Unit Price', x: 480, align: 'right' as const },
        { text: 'จำนวนเงิน', subText: 'Amount', x: 560, align: 'right' as const }
      ];

      tableHeaders.forEach(header => {
        if (header.align === 'right') {
          // ชิดขวา
          const tw = normal.widthOfTextAtSize(header.text, base);
          page.drawText(header.text, { x: header.x - tw, y: height - 295, size: base, font: normal });
          const sw = normal.widthOfTextAtSize(header.subText, 8);
          page.drawText(header.subText, { x: header.x - sw, y: height - 305, size: 8, font: normal, opacity: 0.8 });
        } else {
          page.drawText(header.text, { x: header.x, y: height - 295, size: base, font: normal });
          page.drawText(header.subText, { x: header.x, y: height - 305, size: 8, font: normal, opacity: 0.8 });
        }
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

      // Helper: วาดข้อความชิดขวาที่ตำแหน่ง x
      const drawTextRight = (text: string, rightX: number, y: number, size: number, f: typeof normal) => {
        const w = f.widthOfTextAtSize(text, size);
        page.drawText(text, { x: rightX - w, y, size, font: f, opacity: 0.9 });
      };

      // Column right edges สำหรับตาราง
      const colQtyRight = 390;      // จำนวน
      const colPriceRight = 480;    // ราคาต่อหน่วย
      const colAmountRight = 560;   // จำนวนเงิน

      // Draw table items
      tableItems.forEach((item, index) => {
        const yPos = height - 330 - (index * 15);

        page.drawText(item.id.toString(), {
          x: 30, y: yPos, size: 10, font: normal, opacity: 0.9
        });

        page.drawText(item.title, {
          x: 90, y: yPos, size: 10, font: normal, opacity: 0.9
        });

        drawTextRight(item.qty.toString(), colQtyRight, yPos, 10, normal);
        drawTextRight(numberFormat(item.price), colPriceRight, yPos, 10, normal);
        drawTextRight(numberFormat(item.total), colAmountRight, yPos, 10, normal);
      });

      // Additional details
      page.drawText('รายละเอียดเพิ่มเติม :', {
        x: 20,
        y: height - 520,
        size: 10,
        font: normal,
        opacity: 0.5
      });

      page.drawText(booking.note || '-', {
        x: 100,
        y: height - 520,
        size: 10,
        font: normal,
        opacity: 0.5
      });

      // room details
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

      page.drawText('Amount in words', {
        x: 30,
        y: height - height + 270,
        size: 8,
        font: normal,
        opacity: 0.8
      });

      // ใช้ grand_total ถ้ามี ไม่งั้นใช้ sum_net
      const totalForWords = booking.grand_total || booking.sum_net;
      const bahtText = `-${numberToThaiText(totalForWords)}-`;

      page.drawText(bahtText, {
        x: 150,
        y: height - height + 280,
        size: base,
        font: normal
      });

      // Payment schedule
      if (booking.date_deposit && booking.date_final) {
        const dateDiff = getDateDiff(booking.period_details?.start_date || booking.created_at, booking.created_at);
        // ใช้ grand_total ถ้ามี (รวม VAT - WHT) ไม่งั้นใช้ sum_net
        const totalPayable = booking.grand_total || booking.sum_net;

        if (dateDiff < 21) {
          page.drawText('ชำระเงินทั้งหมดภายในวันที่', {
            x: 30,
            y: height - height + 250,
            size: base,
            font: normal
          });

          page.drawText(getDate(booking.date_final), {
            x: 150,
            y: height - height + 250,
            size: base,
            font: normal
          });

          page.drawText('จำนวน', {
            x: 250,
            y: height - height + 250,
            size: base,
            font: normal
          });

          page.drawText(numberFormat(totalPayable), {
            x: 290,
            y: height - height + 250,
            size: base,
            font: normal
          });
        } else {
          page.drawText('ชำระเงินมัดจำภายในวันที่', {
            x: 30,
            y: height - height + 250,
            size: base,
            font: normal
          });

          page.drawText('ชำระส่วนที่เหลือภายในวันที่', {
            x: 30,
            y: height - height + 220,
            size: base,
            font: normal
          });

          page.drawText(getDate(booking.date_deposit), {
            x: 150,
            y: height - height + 250,
            size: base,
            font: normal
          });

          page.drawText(getDate(booking.date_final), {
            x: 150,
            y: height - height + 220,
            size: base,
            font: normal
          });

          page.drawText('จำนวน', {
            x: 250,
            y: height - height + 250,
            size: base,
            font: normal
          });

          page.drawText('จำนวน', {
            x: 250,
            y: height - height + 220,
            size: base,
            font: normal
          });

          const depositAmount = booking.sum_deposit > totalPayable ? totalPayable : booking.sum_deposit;

          page.drawText(numberFormat(depositAmount), {
            x: 290,
            y: height - height + 250,
            size: base,
            font: normal
          });

          page.drawText(numberFormat(totalPayable - depositAmount), {
            x: 290,
            y: height - height + 220,
            size: base,
            font: normal
          });
        }
      }

      // Summary - ปรับ Y positions ให้รองรับ VAT/WHT
      let summaryY = height - height + 280;
      const lineSpacing = 15;
      const summaryRightEdge = 560; // ขอบขวาสำหรับจัดตัวเลขชิดขวา

      // Helper: วาดตัวเลขชิดขวา
      const drawAmountRight = (text: string, y: number, useBold = false) => {
        const f = useBold ? bold : normal;
        const textWidth = f.widthOfTextAtSize(text, base);
        page.drawText(text, {
          x: summaryRightEdge - textWidth,
          y,
          size: base,
          font: f
        });
      };

      // Helper: วาด label ชิดขวาที่ขอบซ้ายของตัวเลข
      const summaryLabelRight = 490; // ขอบขวาของ label (ก่อนตัวเลข)
      const drawLabelRight = (text: string, y: number, useBold = false) => {
        const f = useBold ? bold : normal;
        const textWidth = f.widthOfTextAtSize(text, base);
        page.drawText(text, {
          x: summaryLabelRight - textWidth,
          y,
          size: base,
          font: f
        });
      };

      // มูลค่ารวม
      drawLabelRight('มูลค่ารวม', summaryY);
      drawAmountRight(numberFormat((booking.sum_room || 0) + (booking.sum_rate || 0) + (booking.sum_additional_services || 0)), summaryY);
      summaryY -= lineSpacing;

      // ส่วนลด
      drawLabelRight('ส่วนลด', summaryY);
      drawAmountRight(numberFormat(booking.discount || 0), summaryY);
      summaryY -= lineSpacing;

      // ยอดก่อนภาษี (Net Before Tax)
      drawLabelRight('ยอดก่อนภาษี', summaryY);
      drawAmountRight(numberFormat(booking.net_before_tax || booking.sum_net), summaryY);
      summaryY -= lineSpacing;

      // VAT 7% (ถ้าเลือก - รวมในค่าบริการแล้ว)
      if (booking.include_vat) {
        drawLabelRight(`VAT ${booking.vat_rate || 7}% (รวมในราคา)`, summaryY);
        drawAmountRight(numberFormat(booking.vat_amount || 0), summaryY);
        summaryY -= lineSpacing;
      }

      // หัก ณ ที่จ่าย 3% (ถ้าเลือก)
      if (booking.include_wht) {
        drawLabelRight(`หัก ณ ที่จ่าย ${booking.wht_rate || 3}%`, summaryY);
        drawAmountRight(`-${numberFormat(booking.wht_amount || 0)}`, summaryY);
        summaryY -= lineSpacing;
      }

      // ยอดที่ต้องชำระ (Grand Total) - ใช้ bold
      drawLabelRight('ยอดที่ต้องชำระ', summaryY, true);
      const payableAmount = booking.grand_total || booking.sum_net;
      drawAmountRight(numberFormat(payableAmount), summaryY, true);

      // Payment info box
      page.drawRectangle({
        x: 20,
        y: height - height + 200,
        width: 559,
        height: 100,
        borderColor: rgb(0, 0, 0),
        borderWidth: 1,
        borderOpacity: 0.1,
        color: rgb(0, 0, 0),
        opacity: 0.1
      });

      page.drawText('ข้อมูลการชำระเงิน', {
        x: 20,
        y: height - height + 180,
        size: 14,
        font: bold
      });

      page.drawText('กสิกรไทย', {
        x: 20,
        y: height - height + 165,
        size: base,
        font: normal
      });

      page.drawText('บจก.เที่ยวจีน', {
        x: 80,
        y: height - height + 165,
        size: base,
        font: normal
      });

      page.drawText('174-8-16031-7', {
        x: 150,
        y: height - height + 165,
        size: base,
        font: normal
      });

      page.drawText('สาขา เซนทรัล เวสต์เกต', {
        x: 250,
        y: height - height + 165,
        size: base,
        font: normal
      });

      page.drawText('หมายเหตุ', {
        x: 20,
        y: height - height + 140,
        size: 14,
        font: bold
      });

      page.drawText('กรุณาส่งหลักฐานการชำระค่าบริการหรือใบ pay-in มาที่ E-mail : bookingchina23@gmail.com', {
        x: 20,
        y: height - height + 125,
        size: base,
        font: normal
      });

      // WHT note — แสดงเมื่อ booking มี VAT (ใบกำกับภาษี)
      if (booking.include_vat && booking.vat_amount > 0) {
        // คำนวณ: serviceFeeTotal = vatAmount × (100 + vatRate) / vatRate
        const vatRate = booking.vat_rate || 7;
        const serviceFeeTotal = booking.vat_amount * (100 + vatRate) / vatRate;
        const feeExcludeVAT = serviceFeeTotal - booking.vat_amount;
        const whtAmount = Math.round(feeExcludeVAT * 0.03 * 100) / 100;
        if (whtAmount > 0) {
          page.drawText(`กรุณาทำการหัก ณ ที่จ่าย ${numberFormat(whtAmount)} (ยอดหัก ณ ที่จ่าย) บาท`, {
            x: 20,
            y: height - height + 110,
            size: base,
            font: bold
          });
        }
      }

      // Signatures
      page.drawText('ผู้วางบิล', {
        x: 100,
        y: height - height + 50,
        size: base,
        font: normal
      });

      page.drawText(booking.sales_details?.first_name + ' ' + booking.sales_details?.last_name, {
        x: 100,
        y: height - height + 35,
        size: base,
        font: normal
      });

      page.drawText('วันที่/Date : ' + getDate(booking.created_at), {
        x: 100,
        y: height - height + 20,
        size: base,
        font: normal
      });

      page.drawText('ผู้อนุมัติ/Authorized Signature', {
        x: 400,
        y: height - height + 50,
        size: base,
        font: normal
      });

      page.drawText('นายอรรถพล พงภิทักษ์', {
        x: 400,
        y: height - height + 35,
        size: base,
        font: normal
      });

      page.drawText('วันที่/Date : ' + getDate(booking.created_at), {
        x: 400,
        y: height - height + 20,
        size: base,
        font: normal
      });

      // Save the PDFDocument to bytes
      const pdfBytes = await pdfDoc.save();
      const blob = new Blob([new Uint8Array(pdfBytes)], { type: 'application/pdf' });
      const url = URL.createObjectURL(blob);

      // Create a link element
      const link = document.createElement('a');
      link.href = url;

      const filename = `${booking.invoice_id} - บริษัท เที่ยวจีน จำกัด`;
      link.download = filename;

      document.body.appendChild(link);
      link.click();

      // Clean up
      document.body.removeChild(link);
      URL.revokeObjectURL(url);

    } catch (error) {
      console.error('Error creating PDF:', error);
      alert('เกิดข้อผิดพลาดในการสร้าง PDF');
    }
  }

  return (
    <Button
      variant={variant}
      size={size}
      onClick={(e) => {
        e.stopPropagation();
        createPdf();
      }}
      className={`cursor-pointer ${style === "column" ? "flex-col h-full" : ""}`}
    >
      <FileText className={`h-4 w-4 ${style === "column" ? "mb-1" : "mr-1"}`} />
      INVOICE
    </Button>
  );
}

export default PdfInvoiceNew;