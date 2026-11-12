# Custom PDF Header Image — Implementation Plan

## Overview
Agent สามารถอัปโหลดภาพหัว PDF ของตัวเอง (custom header) ได้
- ถ้ามี custom image -> ใช้ภาพนั้นแทน header ที่ generate อัตโนมัติ
- ถ้าไม่มี -> fallback ใช้ header ที่ generate จากข้อมูลบริษัท (Logo + ชื่อ + ที่อยู่ + เลขภาษี)

---

## DB Migration

```sql
ALTER TABLE company_profiles ADD COLUMN pdf_header_image VARCHAR(500);
```

---

## Backend (18 แก้ไข ใน 8 ไฟล์)

### 1. Model — `gofiber/domain/models/company_profile.go`
เพิ่ม field หลัง `Logo` (line 22):
```go
PdfHeaderImage *string `gorm:"size:500;column:pdf_header_image" json:"pdf_header_image"`
```

### 2. DTO — `gofiber/domain/dto/company.go`
เพิ่ม field ใน 3 struct:

| Struct | Field | JSON Tag |
|--------|-------|----------|
| `CompanyProfileResponse` | `PdfHeaderImage *string` | `json:"pdfHeaderImage"` |
| `CreateCompanyProfileRequest` | `PdfHeaderImage *string` | `json:"pdfHeaderImage,omitempty"` |
| `UpdateCompanyProfileRequest` | `PdfHeaderImage *string` | `json:"pdfHeaderImage,omitempty"` |

เพิ่ม struct ใหม่:
```go
type UploadPdfHeaderImageRequest struct {
    UserID    uuid.UUID
    ImageData []byte
    FileName  string
    FileType  string
}
```

### 3. Repository Interface — `gofiber/domain/repositories/company_profile.go`
เพิ่ม method หลัง `UpdateLogo` (line 21):
```go
UpdatePdfHeaderImage(ctx context.Context, id uuid.UUID, imageURL string) error
```

### 4. Repository Impl — `gofiber/infrastructure/postgres/company_profile_repository_impl.go`
เพิ่ม method ใหม่ หลัง `UpdateLogo` (line 160):
```go
func (r *CompanyProfileRepositoryImpl) UpdatePdfHeaderImage(ctx context.Context, id uuid.UUID, imageURL string) error {
    result := r.db.WithContext(ctx).Model(&models.CompanyProfile{}).
        Where("id = ?", id).Update("pdf_header_image", imageURL)
    // ... error handling
}
```

### 5. Service Interface — `gofiber/domain/services/company.go`
เพิ่ม 2 methods หลัง `DeleteLogo` (line 19):
```go
UploadPdfHeaderImage(ctx context.Context, req dto.UploadPdfHeaderImageRequest) (*dto.CompanyProfileResponse, error)
DeletePdfHeaderImage(ctx context.Context, userID uuid.UUID) (*dto.CompanyProfileResponse, error)
```

### 6. Service Impl — `gofiber/application/serviceimpl/company_service_impl.go`
แก้ 6 จุด:

| จุดที่ | Method | การแก้ไข |
|--------|--------|----------|
| A | `toCompanyProfileResponse` (line 495) | เพิ่ม `PdfHeaderImage: cdn.GetFullFileURLPtr(profile.PdfHeaderImage)` |
| B | `CreateProfile` (line 90) | เพิ่ม `PdfHeaderImage: req.PdfHeaderImage` |
| C | `UpdateProfile` (line 212) | เพิ่ม `if req.PdfHeaderImage != nil { profile.PdfHeaderImage = req.PdfHeaderImage }` |
| D | `AdminUpdateProfile` (line 272) | เพิ่ม เหมือน C |
| E | เพิ่ม method ใหม่ | `UploadPdfHeaderImage` (copy pattern จาก `UploadLogo`) |
| F | เพิ่ม method ใหม่ | `DeletePdfHeaderImage` (copy pattern จาก `DeleteLogo`) |
| G | เพิ่ม helper ใหม่ | `generatePdfHeaderImagePath` (path: `pdf-headers/{userID}/{timestamp}.ext`) |

### 7. Handler — `gofiber/interfaces/api/handlers/company_handler.go`
เพิ่ม 2 handlers หลัง `DeleteLogo` (line 196):

| Handler | Method | Form Field | Max Size |
|---------|--------|------------|----------|
| `UploadPdfHeaderImage` | POST | `pdf_header_image` | 5MB |
| `DeletePdfHeaderImage` | DELETE | - | - |

Pattern เหมือน `UploadLogo`/`DeleteLogo` ทุกประการ

### 8. Routes — `gofiber/interfaces/api/routes/company_routes.go`
เพิ่ม 2 routes หลัง logo routes (line 27):
```go
company.Post("/me/pdf-header-image", h.CompanyHandler.UploadPdfHeaderImage)
company.Delete("/me/pdf-header-image", h.CompanyHandler.DeletePdfHeaderImage)
```

---

## Frontend (9 แก้ไข ใน 9 ไฟล์)

### 9. Type Model — `vite/src/types/models/companyProfile.ts`
เพิ่มหลัง `logo` (line 13):
```typescript
pdf_header_image: string | null;
```

### 10. Type Model (Agent) — `vite/src/types/models/agentProfile.ts`
เพิ่ม `pdfHeaderImage?: string | null` ใน `AgentProfile`, `CreateAgentProfileRequest`, `UpdateAgentProfileRequest`

### 11. Auth Response — `vite/src/types/response/auth.ts`
เพิ่มหลัง `logo` (line 56):
```typescript
pdf_header_image: string | null;
```

### 12. API Constants — `vite/src/constants/api.ts`
เพิ่มใน `AGENT_PROFILE_API` (line 130):
```typescript
UPLOAD_PDF_HEADER: '/company-profiles/me/pdf-header-image',
DELETE_PDF_HEADER: '/company-profiles/me/pdf-header-image',
```

### 13. Service — `vite/src/services/agent/agentProfileService.ts`
เพิ่ม 2 methods:
```typescript
uploadPdfHeaderImage(file: File)   // POST FormData
deletePdfHeaderImage()             // DELETE
```

### 14. Store — `vite/src/stores/agent/agentProfileStore.ts`
เพิ่ม 2 mutations:
```typescript
useUploadPdfHeaderImageMutation()
useDeletePdfHeaderImageMutation()
```

### 15. PDF Generator — `vite/src/utils/pdfWithAgentHeader.ts`
แก้ logic header:
```
if (companyProfile.pdf_header_image) {
  // ใช้ภาพ custom เต็มขอบด้านบน
} else {
  // ใช้ header ที่ generate (code ปัจจุบัน)
}
```

### 16. Profile UI — `vite/src/components/agent/profile/ProfileInfoCard.tsx`
เพิ่ม section "ภาพหัว PDF" พร้อมปุ่ม Upload/Delete + Preview

### 17. Agent Profile Page — `vite/src/page/agent/profile/page.tsx`
Wire mutations → `ProfileInfoCard` props

---

## สรุปไฟล์ทั้งหมด (17 ไฟล์)

### Backend (8 ไฟล์)
| # | ไฟล์ | สิ่งที่ทำ |
|---|------|----------|
| 1 | `gofiber/domain/models/company_profile.go` | เพิ่ม field `PdfHeaderImage` |
| 2 | `gofiber/domain/dto/company.go` | เพิ่ม field ใน DTOs + struct ใหม่ |
| 3 | `gofiber/domain/repositories/company_profile.go` | เพิ่ม method interface |
| 4 | `gofiber/infrastructure/postgres/company_profile_repository_impl.go` | เพิ่ม method impl |
| 5 | `gofiber/domain/services/company.go` | เพิ่ม 2 methods interface |
| 6 | `gofiber/application/serviceimpl/company_service_impl.go` | เพิ่ม 2 methods + แก้ 4 จุด |
| 7 | `gofiber/interfaces/api/handlers/company_handler.go` | เพิ่ม 2 handlers |
| 8 | `gofiber/interfaces/api/routes/company_routes.go` | เพิ่ม 2 routes |

### Frontend (9 ไฟล์)
| # | ไฟล์ | สิ่งที่ทำ |
|---|------|----------|
| 9 | `vite/src/types/models/companyProfile.ts` | เพิ่ม field |
| 10 | `vite/src/types/models/agentProfile.ts` | เพิ่ม field |
| 11 | `vite/src/types/response/auth.ts` | เพิ่ม field |
| 12 | `vite/src/constants/api.ts` | เพิ่ม endpoints |
| 13 | `vite/src/services/agent/agentProfileService.ts` | เพิ่ม 2 methods |
| 14 | `vite/src/stores/agent/agentProfileStore.ts` | เพิ่ม 2 mutations |
| 15 | `vite/src/utils/pdfWithAgentHeader.ts` | เพิ่ม custom image logic |
| 16 | `vite/src/components/agent/profile/ProfileInfoCard.tsx` | เพิ่ม UI section |
| 17 | `vite/src/page/agent/profile/page.tsx` | Wire mutations |

### DB
| # | สิ่งที่ทำ |
|---|----------|
| 1 | `ALTER TABLE company_profiles ADD COLUMN pdf_header_image VARCHAR(500);` |

---

## ลำดับการทำงาน (แนะนำ)

1. DB Migration (เพิ่ม column)
2. Backend: Model → DTO → Repo → Service → Handler → Routes
3. Frontend: Types → Constants → Service → Store → PDF Generator → UI
4. ทดสอบ: อัปโหลดภาพ → generate PDF → ตรวจว่า fallback ทำงาน
5. Deploy: scp → docker compose build ทั้ง gofiber + vite
