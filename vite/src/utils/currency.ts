// src/utils/currency.ts
// Utility สำหรับแปลงตัวเลขเป็นตัวอักษรภาษาไทย

const thaiNumbers = [
  'ศูนย์', 'หนึ่ง', 'สอง', 'สาม', 'สี่', 'ห้า', 'หก', 'เจ็ด', 'แปด', 'เก้า'
];

const thaiUnits = [
  '', 'สิบ', 'ร้อย', 'พัน', 'หมื่น', 'แสน', 'ล้าน'
];

/**
 * แปลงตัวเลขเป็นตัวอักษรภาษาไทย
 * @param num - ตัวเลขที่ต้องการแปลง
 * @returns string - ตัวอักษรภาษาไทย
 */
export function numberToThaiText(num: number): string {
  if (num === 0) return 'ศูนย์';
  if (num < 0) return 'ลบ' + numberToThaiText(Math.abs(num));

  const numStr = num.toString();
  const length = numStr.length;

  if (length > 7) {
    // จัดการตัวเลขที่มากกว่า 7 หลัก (ล้าน)
    const millions = Math.floor(num / 1000000);
    const remainder = num % 1000000;

    let result = numberToThaiText(millions) + 'ล้าน';
    if (remainder > 0) {
      result += numberToThaiText(remainder);
    }
    return result;
  }

  let result = '';

  for (let i = 0; i < length; i++) {
    const digit = parseInt(numStr[i]);
    const position = length - i - 1;

    if (digit === 0) continue;

    // กรณีพิเศษสำหรับ "หนึ่งสิบ" = "สิบ"
    if (digit === 1 && position === 1 && length === 2) {
      result += 'สิบ';
      continue;
    }

    // กรณีพิเศษสำหรับ "หนึ่ง" ในหลักหน่วย เมื่อมีหลักสิบ
    if (digit === 1 && position === 0 && length > 1) {
      result += 'เอ็ด';
    } else if (digit === 2 && position === 1) {
      // กรณีพิเศษสำหรับ "ยี่สิบ"
      result += 'ยี่' + thaiUnits[position];
    } else {
      result += thaiNumbers[digit] + thaiUnits[position];
    }
  }

  return result;
}

/**
 * แปลงยอดเงินเป็นตัวอักษรภาษาไทย พร้อมคำว่า "บาท"
 * @param amount - จำนวนเงิน
 * @returns string - ตัวอักษรภาษาไทยพร้อมหน่วย "บาท"
 */
export function currencyToThaiText(amount: number): string {
  if (amount === 0) return 'ศูนย์บาท';

  // แยกส่วนจำนวนเต็มและทศนิยม
  const [integerPart, decimalPart] = amount.toString().split('.');
  const integerAmount = parseInt(integerPart);

  let result = numberToThaiText(integerAmount) + 'บาท';

  // จัดการส่วนทศนิยม (สตางค์)
  if (decimalPart && parseInt(decimalPart) > 0) {
    const satangAmount = parseInt(decimalPart.padEnd(2, '0').substring(0, 2));
    result += numberToThaiText(satangAmount) + 'สตางค์';
  }

  return result;
}

/**
 * แปลงยอดเงินเป็นตัวอักษรภาษาไทยแบบเต็ม (ถ้วน)
 * @param amount - จำนวนเงิน
 * @returns string - ตัวอักษรภาษาไทยพร้อมคำว่า "ถ้วน"
 */
export function currencyToThaiTextFull(amount: number): string {
  const text = currencyToThaiText(amount);

  // ถ้าไม่มีสตางค์ ให้เพิ่ม "ถ้วน"
  if (!text.includes('สตางค์')) {
    return text + 'ถ้วน';
  }

  return text;
}