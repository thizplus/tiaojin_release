# Frontend Structure (React + Vite)

## Path: `_main-system/vite`

---

## Tech Stack

| Category | Technology |
|----------|-----------|
| Framework | React 19.1.1 + TypeScript |
| Build | Vite 7.1.7 |
| Styling | Tailwind CSS v4 + shadcn/ui (40+ components) |
| State | Zustand 5 (18 stores) |
| Server State | TanStack React Query 5 |
| Routing | React Router DOM 7.9 |
| Forms | React Hook Form + Zod |
| HTTP | Axios |
| Charts | Recharts |
| PDF | pdf-lib |
| Icons | Lucide React |
| Toast | Sonner |
| Theme | next-themes (Dark/Light/System) |

---

## Directory Structure

```
vite/src/
├── main.tsx                             # Entry point
├── App.tsx                              # Root (QueryClient, Theme, Router)
├── index.css                            # Global styles + CSS variables
│
├── routes/
│   └── index.tsx                        # Route definitions (384 lines)
│
├── page/                                # Feature-based pages
│   ├── globals/                         # Public: login, register, home, profile
│   ├── admin/                           # Admin: dashboard, tours, bookings, users, reports
│   ├── agent/                           # Agent: dashboard, tours, bookings, reports
│   ├── sale/                            # Sales: dashboard, bookings
│   └── super/                           # Super: system settings
│
├── components/
│   ├── ui/                              # shadcn/ui (40+ Radix-based components)
│   ├── auth/                            # ProtectedRoute, PublicRoute, DashboardRedirect
│   ├── manage/                          # CRUD management (booking, tour, etc.)
│   ├── admin/                           # Admin: filters, selectors
│   ├── agent/                           # Agent-specific
│   ├── public/                          # Public-facing
│   ├── globals/                         # Image, LoadingFallback, DataTable
│   ├── debug/                           # WebSocket debug panel
│   ├── providers/                       # Context providers
│   └── page/                            # Page layout components
│
├── layouts/
│   ├── root-layout.tsx                  # Root wrapper
│   ├── page.tsx                         # Auth pages with sidebar
│   ├── app-sidebar.tsx                  # Main sidebar navigation
│   ├── nav-user.tsx                     # User dropdown
│   └── nav/                             # Nav components
│
├── hooks/                               # 45+ custom hooks
│   ├── useAuth.ts                       # Auth state & actions
│   ├── useAdminBooking.ts               # Admin booking CRUD
│   ├── useAgentBooking.ts               # Agent booking CRUD
│   ├── useSaleBooking.ts                # Sales booking CRUD
│   ├── useTours.ts                      # Tours operations
│   ├── useToursHot.ts                   # Hot deals tours
│   ├── useToursRegular.ts              # Regular tours
│   ├── useTourSearch.ts                 # Tour search & filter
│   ├── usePeriod.ts / usePeriodSearch.ts
│   ├── useCity.ts / useBanner.ts / useSales.ts
│   ├── useApiKey.ts / useUser.ts / useCompanyProfile.ts
│   ├── useSeatHold.ts                   # Seat hold management
│   ├── useAdminReport.ts / useAgentReport.ts / useSalesReport.ts
│   ├── useWebSocket.ts                  # WebSocket connection
│   ├── useFileUpload.ts                 # File upload
│   ├── useDebounce.ts                   # Debounce utility
│   ├── useUrlPagination.ts              # URL-based pagination
│   ├── useBankAccounts.ts               # Bank accounts
│   ├── useAgentLogos.ts                 # Agent logos
│   ├── use-mobile.ts                    # Mobile detection (768px)
│   └── use*Data.ts                      # React Query data hooks
│
├── stores/                              # 18 Zustand stores
│   ├── authStore.ts                     # Auth + token + localStorage persist
│   ├── adminBookingStore.ts             # Admin booking state
│   ├── agentBookingStore.ts             # Agent booking state
│   ├── tourStore.ts                     # Tours state
│   ├── configStore.ts                   # System config
│   ├── bannerStore.ts                   # Banner management
│   ├── cityStore.ts                     # City data
│   ├── userStore.ts                     # User management
│   ├── salesStore.ts                    # Sales agents
│   ├── companyProfileStore.ts           # Company profiles
│   ├── periodStore.ts                   # Tour periods
│   ├── apiKeyStore.ts                   # API keys
│   ├── toursHotStore.ts                 # Hot tours
│   ├── toursRegularStore.ts             # Regular tours
│   ├── bookingStore.ts                  # General bookings
│   └── saleBookingStore.ts              # Sales bookings
│
├── services/                            # 18+ API services
│   ├── apiService.ts                    # Axios wrapper (interceptors, auth, 401 logout)
│   ├── authService.ts                   # Login, register, token refresh
│   ├── config.ts                        # System config API
│   ├── adminBookingService.ts           # Admin booking API
│   └── ... (tour, city, banner, etc.)
│
├── constants/
│   ├── api.ts                           # ALL API endpoint constants
│   ├── agent.ts                         # Agent constants
│   └── seat-hold.ts                     # Seat hold constants
│
├── types/
│   ├── index.ts                         # Barrel exports
│   ├── booking.ts                       # Booking types
│   ├── config.ts                        # Config types
│   ├── tours-regular.ts                 # Tour types
│   ├── models/                          # Domain model interfaces
│   ├── request/                         # Request DTO interfaces
│   └── response/                        # Response DTO interfaces
│
├── theme/
│   ├── theme-provider.tsx               # Dark/Light/System mode
│   └── mode-toggle.tsx                  # Theme toggle button
│
├── contexts/                            # React Context providers
├── lib/
│   ├── websocket-manager.ts             # WebSocket lifecycle
│   └── utils.ts                         # Helpers (cn, etc.)
└── assets/                              # Images, SVGs
```

---

## Routing

### Public
| Path | หน้า |
|------|------|
| `/` | Home |
| `/login` | เข้าสู่ระบบ |
| `/register` | สมัครสมาชิก |

### Admin (role: admin, super)
| Path | หน้า |
|------|------|
| `/admin/dashboard` | Dashboard & Analytics |
| `/admin/tours` | จัดการทัวร์ |
| `/admin/tours-regular` | ทัวร์ปกติ |
| `/admin/tours-special` | ทัวร์พิเศษ |
| `/admin/tours/:tourId/periods` | จัดการ Period |
| `/admin/bookings` | จัดการ Booking |
| `/admin/cities` | จัดการเมือง |
| `/admin/users` | จัดการผู้ใช้ |
| `/admin/sales` | จัดการพนักงานขาย |
| `/admin/profiles` | Company profiles |
| `/admin/banners` | จัดการแบนเนอร์ |
| `/admin/notifications` | การแจ้งเตือน |
| `/admin/reports-tours` | รายงานทัวร์ |
| `/admin/reports-daterange` | รายงานตามช่วงเวลา |

### Agent (role: agent)
| Path | หน้า |
|------|------|
| `/agent/dashboard` | Dashboard |
| `/agent/tours` | ดูทัวร์ |
| `/agent/tours-hot` | ทัวร์ Hot deals |
| `/agent/cities/:cityId` | ทัวร์ตามเมือง |
| `/agent/bookings` | Booking ของตัวเอง |
| `/agent/profile` | โปรไฟล์ |
| `/agent/notifications` | การแจ้งเตือน |
| `/agent/reports` | รายงาน |
| `/agent/api-keys` | จัดการ API Key |

### Sales (role: sales)
| Path | หน้า |
|------|------|
| `/sale/dashboard` | Dashboard |
| `/sale/bookings` | Booking ของตัวเอง |

### Super (role: super)
| Path | หน้า |
|------|------|
| `/super/system-settings` | ตั้งค่าระบบ |

---

## Data Flow Pattern

```
Page Component
    ↓
Custom Hook (useXXX)
    ↓
Zustand Store (state + mutations)
    ↓
Service (API calls)
    ↓
apiClient (Axios + interceptors)
    ↓
Backend API (:8080)
```

---

## State Management

| ประเภทข้อมูล | เครื่องมือ | ที่เก็บ |
|-------------|-----------|--------|
| Auth/User | Zustand + localStorage persist | `stores/authStore.ts` |
| Server Data (fetch) | React Query | `hooks/use*Data.ts` |
| Feature State (CRUD) | Zustand | `stores/*.ts` |
| UI State (local) | useState | Component |

---

## Auth Flow

```
1. Login → JWT token → localStorage
2. Axios interceptor → auto-add Bearer token
3. 401 response → forceLogout → clear queries → redirect /login
4. ProtectedRoute → check role → render or redirect
5. DashboardRedirect → route to role-specific dashboard
```

---

## Key Booking Components

| Component | หน้าที่ |
|-----------|---------|
| BookingCardList / BookingGridCard | แสดงรายการ Booking |
| BookingDetailSheet | รายละเอียด Booking (sidebar) |
| CreatePaymentDialog | สร้างรายการชำระ |
| PaymentSlipUploadSheet | อัพโหลดสลิป |
| PaymentSlipVerifyDialog | Admin ตรวจสอบสลิป |
| PaymentSlipReverseDialog | ยกเลิกการชำระ |
| CancelBookingDialog | ยกเลิก Booking |
| DeleteBookingDialog | ลบ Booking |
| PassportUploadSheet | อัพโหลดพาสปอร์ต + OCR |
| PaymentDeadlineNotice | แจ้งเตือนหมดอายุ |
| UpdateBookingStatusDialog | เปลี่ยนสถานะ |

---

## Styling

- **Tailwind CSS v4** + CSS variables (dark/light)
- **shadcn/ui** components (Radix-based)
- **Semantic CSS classes**: `status-pending`, `status-success`, `status-danger`, `status-info`
- **Thai fonts**: CSPraJad (custom)
- **Base font size**: 18px
- **Responsive**: `useIsMobile()` hook (breakpoint: 768px)
- **Theme**: Dark/Light/System via `next-themes`

---

## Dev Config

```
Dev server:     port 5173
API proxy:      /api → http://localhost:8080
Path alias:     @/ → ./src/
Allowed hosts:  ngrok URL, localhost
```

---

## Special Features

| Feature | รายละเอียด |
|---------|-----------|
| WebSocket | Real-time booking updates & notifications |
| PDF Generation | Invoice สร้างจาก pdf-lib + Thai font |
| File Upload | Multi-file, progress tracking, FormData |
| Data Tables | TanStack Table + pagination, filter, sort |
| Form Validation | React Hook Form + Zod (Thai error messages) |
| Image Lightbox | fslightbox-react |
| Charts | Recharts (dashboard analytics) |
