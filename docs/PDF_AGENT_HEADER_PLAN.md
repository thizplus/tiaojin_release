# Item #7: PDF ดาวน์โหลดโปรแกรม แสดง Header ของ Agent

## สถานะปัจจุบัน

ปุ่มที่มีอยู่ใน TourCard:
- **PDF** — ดาวน์โหลดไฟล์ PDF โปรแกรมทัวร์ตรงๆ (ไม่มี header)
- **DOCX** — ดาวน์โหลดไฟล์ Word ตรงๆ
- **รูปภาพ** — ดูรูปภาพ

## สิ่งที่ต้องทำ

เพิ่มปุ่ม **"PDF + หัวบริษัท"** ข้างปุ่ม PDF เดิม

เมื่อกด:
1. สร้างหน้าปก (Cover Page) ด้วยข้อมูลบริษัท Agent
2. ดึงไฟล์ PDF โปรแกรมทัวร์จาก Storage
3. รวมหน้าปก + PDF เดิม เป็นไฟล์เดียว
4. ดาวน์โหลด

## Cover Page แสดงอะไรบ้าง

ข้อมูลจาก `AgentProfile` / `CompanyProfile`:
- โลโก้บริษัท (`logo`)
- ชื่อบริษัท (`name`)
- ที่อยู่ (`address`, `province`, `zipcode`)
- เลขประจำตัวผู้เสียภาษี (`tax_id`)
- เบอร์โทร (`phone`)

ข้อมูลทัวร์:
- ชื่อโปรแกรม (`tour.title`)
- รหัสทัวร์ (`tour.code`)
- เมือง (`tour.city_title`)

## Layout หน้าปก

```
+---------------------------------------+
|  [Logo Agent]                         |
|  ชื่อบริษัท Agent                      |
|  ที่อยู่                               |
|  เลขประจำตัวผู้เสียภาษี / เบอร์โทร      |
|                                       |
|  ─────────────────────────────────    |
|                                       |
|  โปรแกรมทัวร์:                         |
|  [ชื่อทัวร์]                           |
|  รหัส: [tour code]                    |
|  เมือง: [city]                        |
|                                       |
+---------------------------------------+
| หน้า 2+: PDF โปรแกรมทัวร์เดิม          |
+---------------------------------------+
```

## แนวทางทำ (Approach)

### ใช้ pdf-lib (มีในโปรเจคอยู่แล้ว)

```
1. Fetch agent profile → ได้ logo URL, name, address, etc.
2. Fetch logo image → convert เป็น bytes สำหรับ embed
3. Fetch PDF ต้นฉบับ จาก getStorageUrl(tour.file_pdf)
4. สร้าง cover page ด้วย pdf-lib:
   - PDFDocument.create()
   - addPage() → วาด logo, text
5. Load PDF ต้นฉบับ: PDFDocument.load(pdfBytes)
6. Copy ทุกหน้าจาก PDF ต้นฉบับ ไป document ใหม่
7. Download
```

### ข้อควรระวัง
- Logo ต้อง fetch เป็น ArrayBuffer ผ่าน getStorageUrl()
- PDF ต้นฉบับ อาจ CORS block → ใช้ backend proxy (getStorageUrl)
- Font ภาษาไทย → ใช้ font เดียวกับ pdfInvoiceNew.tsx (THSarabunNew)
- ถ้า PDF ต้นฉบับ encrypted → pdf-lib อาจ load ไม่ได้ (fallback: ดาวน์โหลดแยก)

## ไฟล์ที่ต้องสร้าง/แก้

| ไฟล์ | การเปลี่ยนแปลง |
|------|---------------|
| `vite/src/utils/pdfWithAgentHeader.ts` | **(สร้างใหม่)** ฟังก์ชัน generate cover page + merge PDF |
| `vite/src/components/tours-regular/TourCard.tsx` | เพิ่มปุ่ม "PDF + หัวบริษัท" ข้างปุ่ม PDF |

## Data Flow

```
ปุ่ม "PDF + หัวบริษัท" click
  → ดึง agent profile (จาก auth store หรือ useAgentProfile)
  → fetch logo image (getStorageUrl)
  → fetch tour PDF (getStorageUrl)
  → pdf-lib: สร้าง cover page + merge
  → download combined PDF
```

## เฉพาะ Role

- **Agent** — เห็นปุ่ม "PDF + หัวบริษัท" (ใช้ company profile ตัวเอง)
- **Admin/Super** — ไม่ต้องมี (หรือเลือก agent ได้ถ้าต้องการ)
- **Sale** — ใช้ company profile ของ agent ที่เป็นเจ้าของ booking
