# Admin แก้ไข Company Profile — แผนงาน

## สถานะปัจจุบัน

### มีแล้ว (60%)
- Backend: `UpdateProfile()` service + repository ใช้งานได้ (สำหรับ user แก้ไขของตัวเอง)
- Backend: `GetProfileByID()` handler + route
- Frontend: `ViewCompanyProfileDialog` (ดูอย่างเดียว read-only)
- Frontend: Table มีปุ่ม View + Approve

### ขาด
- Backend: ไม่มี route `PUT /admin/companies/:id` สำหรับ admin แก้ไข
- Backend: ไม่มี handler method สำหรับ admin update
- Backend: `UpdateProfile()` เช็ค ownership (user แก้ได้แค่ของตัวเอง)
- Frontend: ไม่มี `EditCompanyProfileDialog`
- Frontend: ไม่มี updateCompanyProfile ใน service/store/hook
- Frontend: ไม่มีปุ่ม "แก้ไข" ใน table actions

## สิ่งที่ต้องทำ

### Backend (3 ไฟล์)

| ไฟล์ | การเปลี่ยนแปลง |
|------|---------------|
| `handlers/company_handler.go` | เพิ่ม `AdminUpdateProfile()` handler — รับ profile ID จาก URL, update โดยไม่เช็ค ownership |
| `routes/admin_routes.go` | เพิ่ม `PUT /admin/companies/:id` route |
| `serviceimpl/company_service_impl.go` | เพิ่ม `AdminUpdateProfile(ctx, profileID, req)` — เหมือน UpdateProfile แต่ใช้ profileID แทน userID |

### Frontend (6 ไฟล์)

| ไฟล์ | การเปลี่ยนแปลง |
|------|---------------|
| `services/admin/companyProfileService.ts` | เพิ่ม `updateCompanyProfile(id, data)` — `PUT /admin/companies/:id` |
| `stores/companyProfileStore.ts` | เพิ่ม `useUpdateCompanyProfileMutation()` |
| `hooks/useCompanyProfile.ts` | เพิ่ม `updateCompanyProfile` action + loading state |
| `components/manage/profile/EditCompanyProfileDialog.tsx` | **(สร้างใหม่)** form แก้ไข name, phone, address, province, zipcode |
| `components/manage/profile/company-profile-columns.tsx` | เพิ่มปุ่ม "แก้ไข" ใน dropdown menu |
| `page/admin/profiles/page.tsx` | เพิ่ม state + handler สำหรับ Edit dialog |

## Fields ที่ Admin แก้ไขได้

| Field | แก้ไขได้ | หมายเหตุ |
|-------|---------|---------|
| `name` | ✅ | ชื่อบริษัท |
| `phone` | ✅ | เบอร์โทร |
| `address` | ✅ | ที่อยู่ |
| `province` | ✅ | จังหวัด |
| `zipcode` | ✅ | รหัสไปรษณีย์ |
| `tax_id` | ✅ | เลขประจำตัวผู้เสียภาษี (validate uniqueness) |
| `logo` | ✅ | โลโก้ (FileUpload) |
| `atta` | ❌ | เอกสารอนุญาต — ไม่ควรให้ admin แก้ |
| `dbd` | ❌ | หนังสือรับรอง — ไม่ควรให้ admin แก้ |
| `pp20` | ❌ | ภ.พ.20 — ไม่ควรให้ admin แก้ |

## Data Flow

```
Admin กดปุ่ม "แก้ไข" ใน table
  → เปิด EditCompanyProfileDialog (pre-fill ข้อมูลเดิม)
  → Admin แก้ไข fields
  → กด "บันทึก"
  → PUT /api/v1/admin/companies/:id
  → AdminUpdateProfile handler
  → AdminUpdateProfile service (ไม่เช็ค ownership)
  → Repository.Update()
  → Refetch list
  → Toast "แก้ไขสำเร็จ"
```
