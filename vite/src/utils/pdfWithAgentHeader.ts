import { PDFDocument, rgb } from 'pdf-lib';
import fontkit from '@pdf-lib/fontkit';
import { getStorageUrl } from '@/lib/utils';
import type { CompanyProfile } from '@/types';

interface GenerateOptions {
  companyProfile: CompanyProfile;
  originalPdfUrl: string;
  filename?: string;
  tourTitle?: string;
  tourCode?: string;
  cityTitle?: string;
  tourThumbnail?: string | null;
}

/**
 * สร้าง PDF หน้าปกบริษัท Agent (แนวนอน) + ภาพปกทัวร์ + แนบ PDF โปรแกรมทัวร์ต่อท้าย
 * Layout: Logo ซ้าย | ข้อมูลบริษัทกลาง | เลขภาษีขวา → เส้นทอง → ภาพปกทัวร์ → ข้อมูลทัวร์
 */
export async function generatePdfWithAgentHeader({
  companyProfile,
  originalPdfUrl,
  filename = 'program.pdf',
  tourTitle,
  tourCode,
  cityTitle,
  tourThumbnail,
}: GenerateOptions): Promise<void> {
  // 1. Load fonts
  const [fontNormalBytes, fontBoldBytes] = await Promise.all([
    fetch('/fonts/CSPraJad.otf').then(r => r.arrayBuffer()),
    fetch('/fonts/CSPraJad-bold.otf').then(r => r.arrayBuffer()),
  ]);

  // 2. Load original PDF
  const originalPdfBytes = await fetch(getStorageUrl(originalPdfUrl)).then(r => r.arrayBuffer());
  const originalPdf = await PDFDocument.load(originalPdfBytes, { ignoreEncryption: true });

  // 3. Create new PDF with cover page
  const pdfDoc = await PDFDocument.create();
  pdfDoc.registerFontkit(fontkit);

  const normal = await pdfDoc.embedFont(fontNormalBytes);
  const bold = await pdfDoc.embedFont(fontBoldBytes);

  // 4. Create cover page (A4: 595 x 842)
  const page = pdfDoc.addPage([595, 842]);
  const { width, height } = page.getSize();

  // === สี ===
  const goldColor = rgb(0.76, 0.6, 0.2);
  const darkText = rgb(0.12, 0.12, 0.12);
  const grayText = rgb(0.35, 0.35, 0.35);
  const lightGray = rgb(0.55, 0.55, 0.55);
  const headerBg = rgb(0.97, 0.97, 0.97);

  // === Layout ===
  const margin = 40;
  const headerTop = height - 30;
  const headerHeight = 130;
  const headerBottom = headerTop - headerHeight;
  const logoAreaWidth = 120;
  const taxAreaWidth = 140;
  const centerStartX = margin + logoAreaWidth + 10;
  const centerEndX = width - margin - taxAreaWidth;

  // -- พื้นหลัง Header --
  page.drawRectangle({
    x: margin,
    y: headerBottom,
    width: width - margin * 2,
    height: headerHeight,
    color: headerBg,
    borderColor: rgb(0.88, 0.88, 0.88),
    borderWidth: 0.5,
  });

  // === ฝั่งซ้าย: Logo ===
  if (companyProfile.logo) {
    try {
      const logoBytes = await fetch(getStorageUrl(companyProfile.logo)).then(r => r.arrayBuffer());
      const logoUint8 = new Uint8Array(logoBytes);

      let logoImage;
      if (logoUint8[0] === 0x89 && logoUint8[1] === 0x50) {
        logoImage = await pdfDoc.embedPng(logoBytes);
      } else {
        logoImage = await pdfDoc.embedJpg(logoBytes);
      }

      const maxW = 100, maxH = 90;
      const scale = Math.min(maxW / logoImage.width, maxH / logoImage.height, 1);
      const logoW = logoImage.width * scale;
      const logoH = logoImage.height * scale;

      const logoX = margin + (logoAreaWidth - logoW) / 2;
      const logoY = headerBottom + (headerHeight - logoH) / 2 + 5;

      page.drawImage(logoImage, {
        x: logoX,
        y: logoY,
        width: logoW,
        height: logoH,
      });
    } catch {
      // logo load ไม่ได้ ข้ามไป
    }
  }

  // ATTA ใบอนุญาต (ใต้ logo area)
  if (companyProfile.atta) {
    const attaText = `ใบอนุญาตเลขที่: ${companyProfile.atta}`;
    const attaWidth = normal.widthOfTextAtSize(attaText, 8);
    const attaX = margin + (logoAreaWidth - attaWidth) / 2;
    page.drawText(attaText, {
      x: attaX,
      y: headerBottom + 8,
      size: 8,
      font: normal,
      color: lightGray,
    });
  }

  // === ฝั่งกลาง: ชื่อบริษัท + ที่อยู่ + เบอร์โทร ===
  const centerMidX = (centerStartX + centerEndX) / 2;

  // ชื่อบริษัท (bold ใหญ่)
  const companyName = companyProfile.name;
  const nameSize = companyName.length > 25 ? 18 : 22;
  const nameWidth = bold.widthOfTextAtSize(companyName, nameSize);
  page.drawText(companyName, {
    x: centerMidX - nameWidth / 2,
    y: headerTop - 35,
    size: nameSize,
    font: bold,
    color: darkText,
  });

  // ที่อยู่ บรรทัดที่ 1
  if (companyProfile.address) {
    const addr1Width = normal.widthOfTextAtSize(companyProfile.address, 10);
    page.drawText(companyProfile.address, {
      x: centerMidX - addr1Width / 2,
      y: headerTop - 58,
      size: 10,
      font: normal,
      color: grayText,
    });
  }

  // ที่อยู่ บรรทัดที่ 2 (จังหวัด + รหัสไปรษณีย์)
  const addressLine2 = [companyProfile.province, companyProfile.zipcode].filter(Boolean).join(' ');
  if (addressLine2) {
    const addr2Width = normal.widthOfTextAtSize(addressLine2, 10);
    page.drawText(addressLine2, {
      x: centerMidX - addr2Width / 2,
      y: headerTop - 73,
      size: 10,
      font: normal,
      color: grayText,
    });
  }

  // เบอร์โทร
  if (companyProfile.phone) {
    const phoneText = `Tel: ${companyProfile.phone}`;
    const phoneWidth = normal.widthOfTextAtSize(phoneText, 10);
    page.drawText(phoneText, {
      x: centerMidX - phoneWidth / 2,
      y: headerTop - 90,
      size: 10,
      font: normal,
      color: grayText,
    });
  }

  // === ฝั่งขวา: เลขประจำตัวผู้เสียภาษี ===
  if (companyProfile.tax_id) {
    const taxLabelText = 'เลขประจำตัวผู้เสียภาษี';
    const taxValueText = companyProfile.tax_id;
    const rightX = width - margin - taxAreaWidth;

    const taxLabelWidth = normal.widthOfTextAtSize(taxLabelText, 8);
    page.drawText(taxLabelText, {
      x: rightX + (taxAreaWidth - taxLabelWidth) / 2,
      y: headerTop - 55,
      size: 8,
      font: normal,
      color: lightGray,
    });

    const taxValueWidth = bold.widthOfTextAtSize(taxValueText, 11);
    page.drawText(taxValueText, {
      x: rightX + (taxAreaWidth - taxValueWidth) / 2,
      y: headerTop - 70,
      size: 11,
      font: bold,
      color: darkText,
    });
  }

  // === เส้นทอง คั่น ===
  const goldLineY = headerBottom - 8;
  page.drawRectangle({
    x: margin,
    y: goldLineY,
    width: width - margin * 2,
    height: 3,
    color: goldColor,
  });

  // === ภาพปกทัวร์ (Thumbnail) — เต็มความกว้างไม่มี padding ===
  let nextY = goldLineY - 10;

  if (tourThumbnail) {
    try {
      const thumbBytes = await fetch(getStorageUrl(tourThumbnail)).then(r => r.arrayBuffer());
      const thumbUint8 = new Uint8Array(thumbBytes);

      let thumbImage;
      if (thumbUint8[0] === 0x89 && thumbUint8[1] === 0x50) {
        thumbImage = await pdfDoc.embedPng(thumbBytes);
      } else {
        thumbImage = await pdfDoc.embedJpg(thumbBytes);
      }

      // เต็มความกว้างหน้ากระดาษ (ไม่มี margin) สูงตามสัดส่วน
      const imgW = width;
      const imgScale = imgW / thumbImage.width;
      const imgH = thumbImage.height * imgScale;
      // จำกัดไม่เกินพื้นที่ว่างที่เหลือ (เผื่อที่ให้ข้อมูลทัวร์ด้านล่าง)
      const maxImgH = nextY - 120;
      const finalH = Math.min(imgH, maxImgH);
      const finalW = finalH < imgH ? (finalH / imgH) * imgW : imgW;

      const imgX = width / 2 - finalW / 2;
      const imgY = nextY - finalH;

      page.drawImage(thumbImage, {
        x: imgX,
        y: imgY,
        width: finalW,
        height: finalH,
      });

      nextY = imgY - 15;
    } catch {
      // thumbnail load ไม่ได้ ข้ามไป
    }
  }

  // === ข้อมูลทัวร์ ===
  let infoY = nextY - 10;

  // "โปรแกรมทัวร์"
  const programLabel = 'โปรแกรมทัวร์';
  const programWidth = bold.widthOfTextAtSize(programLabel, 24);
  page.drawText(programLabel, {
    x: width / 2 - programWidth / 2,
    y: infoY,
    size: 24,
    font: bold,
    color: darkText,
  });

  // ชื่อทัวร์ (word wrap)
  if (tourTitle) {
    const maxCharsPerLine = 45;
    const lines: string[] = [];
    let remaining = tourTitle;
    while (remaining.length > maxCharsPerLine) {
      const breakAt = remaining.lastIndexOf(' ', maxCharsPerLine);
      const cutAt = breakAt > 0 ? breakAt : maxCharsPerLine;
      lines.push(remaining.substring(0, cutAt));
      remaining = remaining.substring(cutAt).trim();
    }
    if (remaining) lines.push(remaining);

    infoY -= 40;
    for (const line of lines) {
      const lineWidth = bold.widthOfTextAtSize(line, 16);
      page.drawText(line, {
        x: width / 2 - lineWidth / 2,
        y: infoY,
        size: 16,
        font: bold,
        color: darkText,
      });
      infoY -= 24;
    }
  }

  // รหัสทัวร์
  if (tourCode) {
    infoY -= 5;
    const codeText = `รหัสทัวร์: ${tourCode}`;
    const codeWidth = normal.widthOfTextAtSize(codeText, 14);
    page.drawText(codeText, {
      x: width / 2 - codeWidth / 2,
      y: infoY,
      size: 14,
      font: normal,
      color: grayText,
    });
    infoY -= 24;
  }

  // เมือง
  if (cityTitle) {
    const cityWidth = normal.widthOfTextAtSize(cityTitle, 14);
    page.drawText(cityTitle, {
      x: width / 2 - cityWidth / 2,
      y: infoY,
      size: 14,
      font: normal,
      color: grayText,
    });
  }

  // 5. Copy all pages from original PDF
  const copiedPages = await pdfDoc.copyPages(originalPdf, originalPdf.getPageIndices());
  for (const copiedPage of copiedPages) {
    pdfDoc.addPage(copiedPage);
  }

  // 6. Download
  const pdfBytes = await pdfDoc.save();
  const blob = new Blob([pdfBytes], { type: 'application/pdf' });
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = filename;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  URL.revokeObjectURL(url);
}
