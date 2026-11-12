# แผนการทำหน้า Home Page (Public)

## ปัญหาปัจจุบัน
- หน้า `/` ถูกครอบด้วย `ProtectedRoute` บังคับให้ต้อง login ก่อน
- ต้องการให้หน้า `/` เป็น public page แสดงรายการโปรแกรมทัวร์

---

## แผนการแก้ไข

### ขั้นตอนที่ 1: สร้างหน้า HomePage

**ไฟล์ที่สร้างใหม่:** `vite/src/page/globals/home/page.tsx`

```tsx
// หน้า public แสดงรายการทัวร์
- ใช้ ToursRegularList component ที่มีอยู่แล้ว
- ใช้ TourSearchForm สำหรับค้นหา
- แสดง Banner (ถ้าต้องการ)
- ปุ่ม Login สำหรับ agent ที่ต้องการจอง
- canBook = false (ไม่แสดงปุ่มจอง เพราะยังไม่ได้ login)
```

**โครงสร้างหน้า:**
```
┌─────────────────────────────────────────────┐
│  Header: Logo + ปุ่ม Login/Register          │
├─────────────────────────────────────────────┤
│  Banner Carousel (optional)                  │
├─────────────────────────────────────────────┤
│  Search Form                                 │
├─────────────────────────────────────────────┤
│  Tour List (ToursRegularList) ห้ามแสดงค่า commission │
│  - แสดงทัวร์ที่ publish                      │
│  - ไม่มีปุ่มจอง (view only)                  │
│  - แสดง periods, ราคา, รายละเอียด            │
└─────────────────────────────────────────────┘
```

---

### ขั้นตอนที่ 2: แก้ไข Routes

**ไฟล์ที่แก้ไข:** `vite/src/routes/index.tsx`

**ก่อนแก้ไข:**
```tsx
<Route
  path="/"
  element={
    <ProtectedRoute>
      <DashboardRedirect />
    </ProtectedRoute>
  }
/>
```

**หลังแก้ไข:**
```tsx
// Import หน้า HomePage
const HomePage = lazy(() => import('@/page/globals/home/page'));

// เปลี่ยนเป็น PublicRoute
<Route
  path="/"
  element={
    <PublicRoute>
      <HomePage />
    </PublicRoute>
  }
/>
```

---

### ขั้นตอนที่ 3: (Optional) สร้าง Public Layout

ถ้าต้องการ layout แยกสำหรับหน้า public (มี header, footer แตกต่างจากหน้า dashboard)

**ไฟล์ที่สร้างใหม่:** `vite/src/layouts/public-layout.tsx`

```tsx
// Layout สำหรับหน้า public
- Header: Logo, Navigation, Login/Register buttons
- Main Content
- Footer (optional)
```

---

## สรุปไฟล์ที่ต้องแก้ไข/สร้าง

| ไฟล์ | Action | รายละเอียด |
|------|--------|------------|
| `vite/src/page/globals/home/page.tsx` | **สร้างใหม่** | หน้า HomePage แสดงทัวร์ |
| `vite/src/routes/index.tsx` | **แก้ไข** | เปลี่ยน `/` เป็น PublicRoute |
| `vite/src/layouts/public-layout.tsx` | **สร้างใหม่ (optional)** | Layout สำหรับหน้า public |

---

## Components ที่ใช้ (มีอยู่แล้ว)

| Component | Path | ใช้ทำอะไร |
|-----------|------|----------|
| `ToursRegularList` | `components/tours-regular/ToursRegularList.tsx` | แสดงรายการทัวร์ |
| `TourSearchForm` | `components/tours-regular/TourSearchForm.tsx` | ฟอร์มค้นหาทัวร์ |
| `TourCard` | `components/tours-regular/TourCard.tsx` | การ์ดแสดงทัวร์ |
| `PublicRoute` | `components/auth/PublicRoute.tsx` | Route สำหรับหน้า public |

---

## Props สำหรับ ToursRegularList (หน้า Public)

```tsx
<ToursRegularList
  canBook={false}           // ไม่แสดงปุ่มจอง
  userStatus="guest"        // หรือไม่ส่ง
  initialParams={{
    status: 'publish',      // แสดงเฉพาะทัวร์ที่ publish
    available_only: true,   // แสดงเฉพาะทัวร์ที่มี period ว่าง
  }}
/>
```

---

## Flow หลังแก้ไข

```
User เข้า / (home)
    │
    ▼
┌─────────────────────┐
│  HomePage (Public)   │ ◄─── ไม่ต้อง login
│  - ดูรายการทัวร์     │
│  - ค้นหาทัวร์        │
│  - ดูราคา/periods    │
└─────────────────────┘
    │
    │ คลิก Login
    ▼
┌─────────────────────┐
│  Login Page          │
└─────────────────────┘
    │
    │ Login สำเร็จ
    ▼
┌─────────────────────┐
│  Dashboard (ตาม role)│
│  - Admin Dashboard   │
│  - Agent Dashboard   │
│  - Sales Dashboard   │
└─────────────────────┘
```

---

## ข้อควรพิจารณา

1. **SEO**: หน้า public ควรมี meta tags สำหรับ SEO
2. **Performance**: ใช้ lazy loading สำหรับ component
3. **Mobile**: ตรวจสอบว่า responsive ดี
4. **Navigation**: ผู้ใช้ที่ login แล้วควรมีปุ่มไป dashboard
5. **ใช้ shadcn ui  ในการสร้างนะครับ **

---

## พร้อมดำเนินการเมื่อได้รับอนุมัติ
