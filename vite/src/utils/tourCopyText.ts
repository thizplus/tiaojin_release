import { format } from 'date-fns';
import { th } from 'date-fns/locale';
import { getStorageUrl } from '@/lib/utils';
import type { Tour } from '@/types/tours-regular';
import type { CompanyProfile } from '@/types';

/** หาราคาต่ำสุดจาก rate items ของทุก period */
function getLowestPrice(tour: Tour): number {
  let lowest = Infinity;
  for (const period of tour.available_periods) {
    for (const rate of period.rate_items) {
      const price = rate.sale > 0 ? rate.sale : rate.regular;
      if (price > 0 && price < lowest) lowest = price;
    }
  }
  return lowest === Infinity ? 0 : lowest;
}

/** หาสายการบินจาก period แรก */
function getAirline(tour: Tour): string {
  const period = tour.available_periods[0];
  if (!period?.airline) return '';
  const airport = period.airport ? ` (${period.airport})` : '';
  return `${period.airline}${airport}`;
}

/** format วันที่เป็นภาษาไทย */
function formatDateTh(dateStr: string): string {
  return format(new Date(dateStr), 'd MMM yyyy', { locale: th });
}

/** หาช่วงเดือนจาก periods */
function getMonthRange(tour: Tour): string {
  if (tour.available_periods.length === 0) return '';
  const months = new Set<string>();
  for (const p of tour.available_periods) {
    months.add(format(new Date(p.start_date), 'MMMM yyyy', { locale: th }));
  }
  return Array.from(months).join(' - ');
}

/**
 * #8 Sale Script — ข้อความโปรโมทสั้นๆ สำหรับส่งลูกค้า
 */
export function generateSaleScript(tour: Tour, companyProfile?: CompanyProfile | null): string {
  const lowestPrice = getLowestPrice(tour);
  const airline = getAirline(tour);
  const monthRange = getMonthRange(tour);
  const pdfUrl = tour.file_pdf ? getStorageUrl(tour.file_pdf) : '';
  const wordUrl = tour.file_word ? getStorageUrl(tour.file_word) : '';

  // แยก detail เป็น highlights (แต่ละบรรทัด)
  const highlights = (tour.detail || '')
    .split('\n')
    .map(l => l.trim())
    .filter(l => l.length > 0)
    .slice(0, 10)
    .map(l => `💭${l}`)
    .join('\n');

  let text = `💐${tour.code} ${tour.title}\n\n`;

  if (airline) {
    text += `✈️ สายการบิน ${airline}\n\n`;
  }

  if (highlights) {
    text += `${highlights}\n\n`;
  }

  if (monthRange) {
    text += `📆 พีเรียดเดินทาง ${monthRange}\n`;
  }

  if (lowestPrice > 0) {
    text += `💵 ราคาเริ่มต้น ${lowestPrice.toLocaleString()}\n`;
  }

  if (pdfUrl) text += `\nPDF 💌 ${pdfUrl}`;
  if (wordUrl) text += `\nWORD 💌 ${wordUrl}`;

  if (companyProfile) {
    text += `\n\n🏢 สนใจติดต่อ ${companyProfile.name}`;
    if (companyProfile.phone) text += `\n📞 โทร ${companyProfile.phone}`;
  }

  return text;
}

/**
 * #9 Copy Details — ข้อมูลครบสำหรับ agent
 */
export function generateTourDetails(tour: Tour, companyProfile?: CompanyProfile | null): string {
  const lowestPrice = getLowestPrice(tour);
  const airline = getAirline(tour);
  const pdfUrl = tour.file_pdf ? getStorageUrl(tour.file_pdf) : '';
  const bannerUrl = tour.thumbnail ? getStorageUrl(tour.thumbnail) : '';

  let text = `📄 รหัส ${tour.code} ${tour.title}\n`;

  if (lowestPrice > 0) {
    text += `🔖 ราคาเริ่มต้น ${lowestPrice.toLocaleString()} บาท\n`;
  }

  if (airline) {
    text += `✈️ เดินทางโดย ${airline}\n`;
  }

  text += `--------------------------------------\n`;

  // แสดงทุก period
  if (tour.available_periods.length > 0) {
    text += `📍 ช่วงวันเดินทาง\n`;
    for (const period of tour.available_periods) {
      text += `${formatDateTh(period.start_date)} - ${formatDateTh(period.end_date)}`;
      if (period.available_seats > 0) {
        text += ` (ว่าง ${period.available_seats} ที่นั่ง)`;
      }
      text += '\n';
    }
    text += `--------------------------------------\n`;
  }

  if (pdfUrl) {
    text += `📰 ชมโปรแกรม PDF ${pdfUrl}\n`;
    text += `--------------------------------------\n`;
  }

  if (bannerUrl) {
    text += `🖼️ แบนเนอร์ ${bannerUrl}\n`;
    text += `--------------------------------------\n`;
  }

  if (companyProfile) {
    text += `🏢 สนใจติดต่อ ${companyProfile.name}\n`;
    if (companyProfile.phone) text += `📞 โทร ${companyProfile.phone}\n`;
    if (companyProfile.tax_id) text += `📋 เลขประจำตัวผู้เสียภาษี ${companyProfile.tax_id}\n`;
  }

  return text;
}
