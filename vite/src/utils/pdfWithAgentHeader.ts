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
 * Layout: Logo ซ้าย | ข้อมูลบริษัทกลาง | เลขภาษีขวา → เส้นทอง → ภาพปกทัวร์ (เต็มขอบ)
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
  const redColor = rgb(0.78, 0.12, 0.12);
  const darkText = rgb(0.12, 0.12, 0.12);
  const grayText = rgb(0.3, 0.3, 0.3);
  const headerBg = rgb(0.97, 0.97, 0.97);

  // ตัวแปรบอกตำแหน่งล่างสุดของ header (ใช้ต่อในส่วน thumbnail/tour info)
  let headerBottomY: number;

  if (companyProfile.pdf_header_image) {
    // === Custom PDF Header Image — วาดเต็มความกว้างหน้ากระดาษ ===
    try {
      const headerImgBytes = await fetch(getStorageUrl(companyProfile.pdf_header_image)).then(r => r.arrayBuffer());
      const headerImgUint8 = new Uint8Array(headerImgBytes);

      let headerImage;
      if (headerImgUint8[0] === 0x89 && headerImgUint8[1] === 0x50) {
        headerImage = await pdfDoc.embedPng(headerImgBytes);
      } else {
        headerImage = await pdfDoc.embedJpg(headerImgBytes);
      }

      // Scale to full page width, cap height at 200pt
      const imgScale = width / headerImage.width;
      const scaledH = headerImage.height * imgScale;
      const finalH = Math.min(scaledH, 200);
      const finalW = finalH < scaledH ? (finalH / scaledH) * width : width;
      const imgX = (width - finalW) / 2;
      const imgY = height - finalH;

      page.drawImage(headerImage, {
        x: imgX,
        y: imgY,
        width: finalW,
        height: finalH,
      });

      headerBottomY = imgY;
    } catch {
      // ภาพหัว PDF โหลดไม่ได้ — fallback ใช้ generated header
      headerBottomY = height;
    }
  } else {
    // === Generated Header (โค้ดเดิม) ===
    const textPadding = 15;
    const headerTop = height;
    const headerHeight = 145;
    const headerBottom = headerTop - headerHeight;
    const logoAreaWidth = 130;
    const taxAreaWidth = 140;
    const centerStartX = logoAreaWidth;

    // -- พื้นหลัง Header (เต็มขอบ) --
    page.drawRectangle({
      x: 0,
      y: headerBottom,
      width: width,
      height: headerHeight,
      color: headerBg,
    });

    // เส้นขอบล่าง header
    page.drawLine({
      start: { x: 0, y: headerBottom },
      end: { x: width, y: headerBottom },
      thickness: 0.5,
      color: rgb(0.85, 0.85, 0.85),
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

        const maxW = 110, maxH = 95;
        const scale = Math.min(maxW / logoImage.width, maxH / logoImage.height, 1);
        const logoW = logoImage.width * scale;
        const logoH = logoImage.height * scale;

        const logoX = textPadding + (logoAreaWidth - textPadding - logoW) / 2;
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

    // หมายเหตุ: companyProfile.atta เก็บเป็น URL ไฟล์ PDF ไม่ใช่เลขใบอนุญาต จึงไม่แสดงใน header

    // === ฝั่งกลาง: ชื่อบริษัท + ที่อยู่ + เบอร์โทร (ชิดซ้าย) ===
    const infoX = centerStartX + 5;

    // ชื่อบริษัท (bold ใหญ่)
    const companyName = companyProfile.name;
    const nameSize = companyName.length > 20 ? 22 : 26;
    page.drawText(companyName, {
      x: infoX,
      y: headerTop - 38,
      size: nameSize,
      font: bold,
      color: darkText,
    });

    // ที่อยู่ บรรทัดที่ 1
    if (companyProfile.address) {
      page.drawText(companyProfile.address, {
        x: infoX,
        y: headerTop - 62,
        size: 12,
        font: normal,
        color: grayText,
      });
    }

    // ที่อยู่ บรรทัดที่ 2 (จังหวัด + รหัสไปรษณีย์)
    const addressLine2 = [companyProfile.province, companyProfile.zipcode].filter(Boolean).join(' ');
    if (addressLine2) {
      page.drawText(addressLine2, {
        x: infoX,
        y: headerTop - 80,
        size: 12,
        font: normal,
        color: grayText,
      });
    }

    // เบอร์โทร
    if (companyProfile.phone) {
      page.drawText(`โทรศัพท์: ${companyProfile.phone}`, {
        x: infoX,
        y: headerTop - 98,
        size: 12,
        font: bold,
        color: grayText,
      });
    }

    // === ฝั่งขวา: พื้นแดงเฉียง + เลขประจำตัวผู้เสียภาษี (ขาว) ===
    const whiteText = rgb(1, 1, 1);
    const slant = 20;

    // วาด parallelogram แดง ด้วย SVG path
    const redW = taxAreaWidth + slant;
    const svgPath = `M ${slant} 0 L ${redW} 0 L ${redW} ${headerHeight} L 0 ${headerHeight} Z`;
    page.drawSvgPath(svgPath, {
      x: width - redW,
      y: headerTop,
      color: redColor,
    });

    if (companyProfile.tax_id) {
      const taxLabelText = 'เลขประจำตัวผู้เสียภาษี';
      const taxValueText = companyProfile.tax_id;
      const redCenterX = width - taxAreaWidth / 2;

      const taxLabelWidth = normal.widthOfTextAtSize(taxLabelText, 10);
      page.drawText(taxLabelText, {
        x: redCenterX - taxLabelWidth / 2,
        y: headerTop - 55,
        size: 10,
        font: normal,
        color: whiteText,
      });

      const taxValueWidth = bold.widthOfTextAtSize(taxValueText, 13);
      page.drawText(taxValueText, {
        x: redCenterX - taxValueWidth / 2,
        y: headerTop - 78,
        size: 13,
        font: bold,
        color: whiteText,
      });
    }

    // === เส้นแดง คั่น (เต็มขอบ) ===
    const lineY = headerBottom - 3;
    page.drawRectangle({
      x: 0,
      y: lineY,
      width: width,
      height: 3,
      color: redColor,
    });

    headerBottomY = lineY;
  }

  // === ภาพปกทัวร์ (Thumbnail) — เต็มขอบไม่มี padding ===
  let nextY = headerBottomY;

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

      // เต็มความกว้างหน้ากระดาษ สูงตามสัดส่วน
      const imgScale = width / thumbImage.width;
      const imgH = thumbImage.height * imgScale;
      // จำกัดไม่เกินพื้นที่ว่าง (เผื่อข้อมูลทัวร์ด้านล่าง)
      const maxImgH = nextY - 100;
      const finalH = Math.min(imgH, maxImgH);
      const finalW = finalH < imgH ? (finalH / imgH) * width : width;

      const imgX = (width - finalW) / 2;
      const imgY = nextY - finalH;

      page.drawImage(thumbImage, {
        x: imgX,
        y: imgY,
        width: finalW,
        height: finalH,
      });

      nextY = imgY - 10;
    } catch {
      // thumbnail load ไม่ได้ ข้ามไป
    }
  }

  // === ข้อมูลทัวร์ ===
  let infoY = nextY - 15;

  // "โปรแกรมทัวร์"
  const programLabel = 'โปรแกรมทัวร์';
  const programWidth = bold.widthOfTextAtSize(programLabel, 22);
  page.drawText(programLabel, {
    x: width / 2 - programWidth / 2,
    y: infoY,
    size: 22,
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

    infoY -= 35;
    for (const line of lines) {
      const lineWidth = bold.widthOfTextAtSize(line, 15);
      page.drawText(line, {
        x: width / 2 - lineWidth / 2,
        y: infoY,
        size: 15,
        font: bold,
        color: darkText,
      });
      infoY -= 22;
    }
  }

  // รหัสทัวร์
  if (tourCode) {
    infoY -= 3;
    const codeText = `รหัสทัวร์: ${tourCode}`;
    const codeWidth = normal.widthOfTextAtSize(codeText, 13);
    page.drawText(codeText, {
      x: width / 2 - codeWidth / 2,
      y: infoY,
      size: 13,
      font: normal,
      color: grayText,
    });
    infoY -= 20;
  }

  // เมือง
  if (cityTitle) {
    const cityWidth = normal.widthOfTextAtSize(cityTitle, 13);
    page.drawText(cityTitle, {
      x: width / 2 - cityWidth / 2,
      y: infoY,
      size: 13,
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
