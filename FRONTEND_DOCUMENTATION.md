# TIAOJIN Frontend Documentation (Vite React)

## Overview

ระบบ Frontend สำหรับ Tiaojin Tour Management System พัฒนาด้วย Vite + React + TypeScript

**Tech Stack:**
- React 19 + TypeScript
- Vite (Build Tool)
- React Router v7 (Routing)
- Zustand + React Query (State Management)
- Tailwind CSS (Styling)
- shadcn/ui + Radix UI (UI Components)
- Axios (HTTP Client)
- Recharts (Charts)

---

## 1. Project Structure (โครงสร้างโปรเจค)

```
vite/src/
├── assets/                          # Static files (images, icons)
├── components/                      # React Components
│   ├── admin/                       # Admin-specific components
│   ├── agent/                       # Agent-specific components
│   ├── auth/                        # Authentication (ProtectedRoute, PublicRoute)
│   ├── debug/                       # Debug components (WebSocket)
│   ├── globals/                     # Global components (DataTable, Loading)
│   ├── manage/                      # Management features
│   ├── page/                        # Page-level components
│   ├── providers/                   # Context providers
│   ├── tours-regular/               # Regular tours components
│   └── ui/                          # UI primitives (shadcn)
├── constants/                       # Constants
├── contexts/                        # React Contexts
├── hooks/                           # Custom React Hooks
├── layouts/                         # Layout components
│   ├── nav/                         # Navigation per role
│   ├── app-sidebar.tsx
│   ├── page.tsx
│   └── root-layout.tsx
├── page/                            # Page components (by role)
│   ├── admin/                       # Admin pages
│   ├── agent/                       # Agent pages
│   ├── globals/                     # Shared pages
│   ├── sale/                        # Sales pages
│   └── super/                       # Super admin pages
├── routes/                          # Route configuration
├── services/                        # API services
├── stores/                          # Zustand stores
├── theme/                           # Theme configuration
├── types/                           # TypeScript types
├── utils/                           # Utility functions
├── lib/                             # Libraries
├── App.tsx                          # Root component
└── main.tsx                         # Entry point
```

---

## 2. Pages by Role (หน้าตามบทบาท)

### 2.1 Super Admin Pages

| Path | Page | Description |
|------|------|-------------|
| `/super/system-settings` | SystemSettingsPage | ตั้งค่าระบบ, Expiration Tiers, Booking Limits |

### 2.2 Admin Pages

| Path | Page | Description |
|------|------|-------------|
| `/admin/dashboard` | AdminDashboard | แดชบอร์ด, สถิติ, กราฟ |
| `/admin/tours` | ToursPage | จัดการทัวร์ทั้งหมด |
| `/admin/tours-regular` | ToursRegularPage | ทัวร์ปกติ |
| `/admin/tours-special` | ToursSpecialPage | ทัวร์พิเศษ |
| `/admin/periods/:tourId/periods` | PeriodsPage | จัดการช่วงเวลาของทัวร์ |
| `/admin/bookings` | BookingsPage | จัดการการจอง |
| `/admin/cities` | CitiesPage | จัดการเมือง |
| `/admin/users` | UsersPage | จัดการผู้ใช้ |
| `/admin/sales` | SalesPage | จัดการบุคลากรขาย |
| `/admin/profiles` | ProfilesPage | จัดการโปรไฟล์บริษัท |
| `/admin/banners` | BannersPage | จัดการแบนเนอร์ |
| `/admin/notifications` | NotificationsPage | การแจ้งเตือน |
| `/admin/reports-tours` | ReportsToursPage | รายงานทัวร์ |
| `/admin/reports-daterange` | ReportsDaterangePage | รายงานตามช่วงวันที่ |

### 2.3 Agent Pages

| Path | Page | Description |
|------|------|-------------|
| `/agent/dashboard` | AgentDashboard | แดชบอร์ด, แบนเนอร์, เมือง |
| `/agent/tours` | ToursPage | ดูทัวร์ที่จองได้ |
| `/agent/cities/:cityId` | ToursByCityPage | ทัวร์ตามเมือง |
| `/agent/tours-hot` | ToursHotPage | ทัวร์ฮอต |
| `/agent/bookings` | BookingsPage | การจองของตัวเอง |
| `/agent/profile` | ProfilePage | โปรไฟล์ |
| `/agent/notifications` | NotificationsPage | การแจ้งเตือน |
| `/agent/reports` | ReportsPage | รายงาน |
| `/agent/api-keys` | APIKeysPage | จัดการ API Keys |

### 2.4 Sales Pages

| Path | Page | Description |
|------|------|-------------|
| `/sale/dashboard` | SaleDashboard | แดชบอร์ด |
| `/sale/bookings` | BookingsPage | การจองของลูกค้า |

### 2.5 Global Pages

| Path | Page | Description |
|------|------|-------------|
| `/login` | LoginPage | เข้าสู่ระบบ |
| `/register` | RegisterPage | ลงทะเบียน |
| `/profile` | ProfilePage | โปรไฟล์ผู้ใช้ |
| `/unauthorized` | UnauthorizedPage | เข้าถึงไม่ได้ |

---

## 3. State Management (การจัดการ State)

### 3.1 Zustand Stores

| Store | Description |
|-------|-------------|
| `authStore.ts` | Authentication state + mutations |
| `adminBookingStore.ts` | Admin booking data |
| `agentBookingStore.ts` | Agent booking data |
| `apiKeyStore.ts` | API Key management |
| `bannerStore.ts` | Banner data |
| `bookingStore.ts` | General booking data |
| `cityStore.ts` | City data |
| `companyProfileStore.ts` | Company profiles |
| `configStore.ts` | System configuration |
| `periodStore.ts` | Tour periods |
| `saleBookingStore.ts` | Sales booking data |
| `salesStore.ts` | Sales staff data |
| `toursHotStore.ts` | Hot tours |
| `toursRegularStore.ts` | Regular tours |
| `tourStore.ts` | General tour data |
| `userStore.ts` | User management |

### 3.2 Auth Store Example

```typescript
export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      token: null,
      isAuthenticated: false,
      setAuth: (token, isAuthenticated) => set({ token, isAuthenticated }),
      setUser: (user) => set({ user }),
      logout: () => set({ user: null, token: null, isAuthenticated: false }),
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({
        token: state.token,
        user: state.user,
        isAuthenticated: state.isAuthenticated
      }),
    }
  )
);
```

### 3.3 React Query Integration

| Hook | Description |
|------|-------------|
| `useUserProfileQuery()` | ดึงข้อมูลโปรไฟล์ผู้ใช้ |
| `useLoginMutation()` | Login |
| `useRegisterMutation()` | Register |
| `useLogoutMutation()` | Logout |
| `useValidateTokenQuery()` | ตรวจสอบ token |
| `useRefreshTokenMutation()` | Refresh token |

---

## 4. Custom Hooks

### 4.1 Authentication Hooks

| Hook | Description |
|------|-------------|
| `useAuth()` | Main auth hook (login, logout, register) |
| `useUserProfileQuery()` | Fetch user profile |

### 4.2 Data Hooks

| Hook | Description |
|------|-------------|
| `useTours()` | Tours management |
| `useToursData()` | Tours listing |
| `useTourSearch()` | Tour search |
| `useCities()` | Cities listing |
| `useCitiesData()` | Cities data |
| `useCity()` | Single city |
| `useBanner()` | Banner management |
| `useBannersData()` | Banners listing |
| `usePublicBanners()` | Public banners |
| `useUser()` | User management |
| `useUsersData()` | Users listing |
| `useSales()` | Sales staff |
| `useSalesData()` | Sales listing |
| `useApiKey()` | API Key management |
| `useApiKeysData()` | API Keys listing |
| `useCompanyProfile()` | Company profile |
| `useCompanyProfilesData()` | Company profiles listing |
| `usePeriod()` | Period management |
| `usePeriodsData()` | Periods listing |
| `useSeatHold()` | Seat hold management |
| `useBooking()` | Booking management |
| `useAdminBooking()` | Admin booking |
| `useAgentBooking()` | Agent booking |
| `useSaleBooking()` | Sale booking |

### 4.3 Report Hooks

| Hook | Description |
|------|-------------|
| `useAdminReport()` | Admin dashboard report |
| `useAgentReport()` | Agent report |
| `useSalesReport()` | Sales report |

### 4.4 Utility Hooks

| Hook | Description |
|------|-------------|
| `useDebounce()` | Debounce values |
| `useFileUpload()` | File upload |
| `use-mobile()` | Mobile detection |
| `use-toast()` | Toast notifications |
| `useWebSocket()` | WebSocket connection |
| `useWebSocketHooks()` | WebSocket utilities |

---

## 5. UI Components

### 5.1 shadcn/ui Components

| Component | Description |
|-----------|-------------|
| `accordion.tsx` | Accordion |
| `alert-dialog.tsx` | Alert dialogs |
| `avatar.tsx` | User avatars |
| `badge.tsx` | Status badges |
| `button.tsx` | Buttons |
| `calendar.tsx` | Date picker calendar |
| `card.tsx` | Card containers |
| `chart.tsx` | Recharts integration |
| `checkbox.tsx` | Checkboxes |
| `command.tsx` | Command palette |
| `data-pagination.tsx` | Data table pagination |
| `date-picker.tsx` | Date picker |
| `dialog.tsx` | Modal dialogs |
| `dropdown-menu.tsx` | Dropdown menus |
| `file-upload.tsx` | File upload |
| `form.tsx` | React Hook Form integration |
| `input.tsx` | Input fields |
| `pagination.tsx` | Pagination controls |
| `popover.tsx` | Popovers |
| `progress.tsx` | Progress bars |
| `scroll-area.tsx` | Scrollable areas |
| `search-input.tsx` | Search inputs |
| `select.tsx` | Select dropdowns |
| `sidebar.tsx` | Sidebar layout |
| `sheet.tsx` | Sheet dialogs |
| `switch.tsx` | Toggle switches |
| `table.tsx` | Data tables |
| `tabs.tsx` | Tabbed content |
| `textarea.tsx` | Text areas |
| `toggle.tsx` | Toggle buttons |
| `tooltip.tsx` | Tooltips |

### 5.2 Layout Components

| Component | Description |
|-----------|-------------|
| `RootLayout` | ครอบคลุมทุกหน้าในแอปพลิเคชัน |
| `PageLayout` | Sidebar + Content area + Breadcrumb |
| `AppSidebar` | Sidebar ตาม role |

### 5.3 Navigation Components

| Component | Role |
|-----------|------|
| `NavRoleAdmin` | Admin navigation menu |
| `NavRoleAgent` | Agent navigation menu |
| `NavRoleSuper` | Super admin navigation menu |
| `NavRoleSale` | Sales staff navigation menu |

---

## 6. Management Components (หน้าจัดการ)

### 6.1 Tour Management

| Component | Description |
|-----------|-------------|
| `TourSearchAndTable` | Search and display tours |
| `CreateTourDialog` | สร้างทัวร์ใหม่ |
| `EditTourDialog` | แก้ไขทัวร์ |
| `DeleteTourDialog` | ลบทัวร์ |
| `TourColumns` | Table columns definition |

### 6.2 Booking Management

| Component | Description |
|-----------|-------------|
| `BookingCardList` | Card list view |
| `BookingDetailSheet` | View booking details |
| `CreatePaymentDialog` | สร้างการชำระเงิน |
| `DeleteBookingDialog` | ลบการจอง |

### 6.3 User Management

| Component | Description |
|-----------|-------------|
| `CreateUserDialog` | สร้างผู้ใช้ |
| `EditUserDialog` | แก้ไขผู้ใช้ |
| `DeleteUserDialog` | ลบผู้ใช้ |
| `ChangeRoleDialog` | เปลี่ยนบทบาท |
| `ChangeStatusDialog` | เปลี่ยนสถานะ |

### 6.4 Period Management

| Component | Description |
|-----------|-------------|
| `PeriodForm` | ฟอร์มครบถ้วน |
| `BasicInfoSection` | ข้อมูลพื้นฐาน |
| `FlightSectorsSection` | Flight sectors |
| `RoomTypesSection` | ประเภทห้อง |
| `RateItemsSection` | รายการราคา |
| `CommissionSection` | Commission settings |
| `AdditionalServicesSection` | บริการเพิ่มเติม |

### 6.5 Other Management

| Feature | Components |
|---------|------------|
| City | Create, Edit, Delete dialogs |
| Banner | Create, Edit, Delete dialogs |
| Sales | Create, Edit, Delete dialogs |
| API Key | Create, Edit, Delete dialogs |
| Seat Hold | Holding dialogs and cards |
| Company Profile | View and approve dialogs |
| System Config | Expiration tiers, booking limits |

---

## 7. Authentication & Authorization

### 7.1 ProtectedRoute Component

```typescript
<ProtectedRoute requiredRoles={['admin', 'super']}>
  <AdminDashboardPage />
</ProtectedRoute>
```

**Features:**
- ตรวจสอบ authentication
- ตรวจสอบ authorization (roles)
- Redirect ไปยัง login ถ้าไม่ได้ auth
- Redirect ไปยัง unauthorized ถ้าไม่มีสิทธิ์

### 7.2 useAuth Hook

```typescript
const {
  user,
  token,
  isAuthenticated,
  isLoading,
  login,
  logout,
  register
} = useAuth();
```

### 7.3 API Interceptor

```typescript
// Request interceptor: เพิ่ม Bearer token
api.interceptors.request.use((config) => {
  const token = useAuthStore.getState().token;
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});
```

---

## 8. WebSocket Integration

### 8.1 WebSocketProvider

```typescript
<WebSocketProvider>
  <App />
</WebSocketProvider>
```

**Features:**
- Real-time seat updates
- Connection status management
- Automatic reconnection (max 10 attempts)
- Notification system

### 8.2 Seat Update Events

| Event | Description |
|-------|-------------|
| `seat_held` | ที่นั่งถูกจอง |
| `seat_released` | ที่นั่งถูกปล่อย |
| `seat_hold_expired` | การจองหมดอายุ |
| `booking_created` | มีการจองใหม่ |
| `booking_updated` | อัพเดทการจอง |
| `booking_cancelled` | ยกเลิกการจอง |
| `booking_expired` | การจองหมดอายุ |

---

## 9. TypeScript Types

### 9.1 User Types

```typescript
type UserRole = 'super' | 'admin' | 'agent' | 'sales';
type UserStatus = 'pending' | 'active' | 'inactive';

interface User {
  id: string;
  email: string;
  role: UserRole;
  status: UserStatus;
  created_at: string;
  updated_at: string;
  company_profile?: CompanyProfile | null;
}
```

### 9.2 Tour Types

```typescript
type TourStatus = 'draft' | 'publish';

interface Tour {
  id: string;
  title: string;
  code: string;
  city_id: string;
  detail: string;
  thumbnail: string;
  status: TourStatus;
  created_at: string;
  updated_at: string;
}
```

### 9.3 Type Files

| File | Description |
|------|-------------|
| `models/user.ts` | User models |
| `models/tour.ts` | Tour models |
| `models/period.ts` | Period models |
| `models/booking.ts` | Booking models |
| `models/sales.ts` | Sales models |
| `models/city.ts` | City models |
| `models/banner.ts` | Banner models |
| `models/apiKey.ts` | API Key models |
| `models/companyProfile.ts` | Company profile models |
| `models/agentProfile.ts` | Agent profile models |
| `models/seat-hold.ts` | Seat hold models |

---

## 10. Services (API Integration)

### 10.1 Service Structure

```
services/
├── apiService.ts                 # Axios instance + interceptors
├── authService.ts                # Authentication
├── bookingService.ts             # Booking management
├── cityService.ts                # City management
├── bannerService.ts              # Banner management
├── apiKeyService.ts              # API Key management
├── reportService.ts              # Reports
├── admin/                        # Admin-specific services
├── agent/                        # Agent-specific services
├── sale/                         # Sale-specific services
└── super/                        # Super admin services
```

### 10.2 API Service Example

```typescript
const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
});
```

---

## 11. Utility Functions

| File | Description |
|------|-------------|
| `bookingStatus.tsx` | Booking status helper |
| `chartConfig.ts` | Recharts configuration |
| `currency.ts` | Currency formatting |
| `exportBookingsCSV.ts` | Export to CSV |
| `navigation-admin.ts` | Admin navigation paths |
| `navigation-agent.ts` | Agent navigation paths |
| `navigation-sale.ts` | Sales navigation paths |
| `navigation-super.ts` | Super admin navigation paths |
| `pdfGenerator.ts` | PDF generation |
| `status-helpers.ts` | Status helper functions |

---

## 12. Configuration

### 12.1 Vite Config

```typescript
export default defineConfig({
  plugins: [react(), tailwindcss(), imagetools()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  server: {
    host: '0.0.0.0',
    port: 5173,
    allowedHosts: ['localhost', '127.0.0.1'],
  },
});
```

### 12.2 Environment Variables

```env
VITE_API_URL=http://localhost:8080/api/v1
VITE_APP_TITLE=Tiaojin
VITE_APP_VERSION=1.0.0
```

---

## 13. Dependencies

### 13.1 Core Dependencies

```json
{
  "react": "^19.1.1",
  "react-dom": "^19.1.1",
  "react-router-dom": "^7.9.1",
  "@tanstack/react-query": "^5.90.1",
  "@tanstack/react-table": "^8.21.3",
  "zustand": "^5.0.8",
  "axios": "^1.12.2"
}
```

### 13.2 UI Dependencies

```json
{
  "@radix-ui/react-*": "^latest",
  "tailwindcss": "^4.1.13",
  "lucide-react": "^0.544.0",
  "class-variance-authority": "^0.7.1"
}
```

### 13.3 Form & Validation

```json
{
  "react-hook-form": "^7.63.0",
  "@hookform/resolvers": "^5.2.2",
  "zod": "^4.1.11"
}
```

### 13.4 Utilities

```json
{
  "recharts": "^2.15.4",
  "date-fns": "^4.1.0",
  "pdf-lib": "^1.17.1",
  "next-themes": "^0.4.6",
  "sonner": "^2.0.7"
}
```

---

## 14. Key Features

1. **Multi-role Authentication** - Super, Admin, Agent, Sales
2. **Real-time Updates** - WebSocket seat updates
3. **Complete Tour Management** - Tours, periods, pricing
4. **Booking Management** - With payment tracking
5. **Report Generation** - Charts with Recharts
6. **City & Banner Management**
7. **User & Sales Staff Management**
8. **API Key Management** - For agents
9. **Dark/Light Theme** - Theme support
10. **Responsive Design** - Mobile-friendly
11. **PDF Generation** - For documents
12. **CSV Export** - For bookings
13. **File Uploads** - Images, documents
14. **Type Safety** - Full TypeScript

---

## 15. File Statistics

| Category | Count |
|----------|-------|
| Source Files | 150+ |
| Roles | 4 |
| Custom Hooks | 30+ |
| Zustand Stores | 16 |
| UI Components | 60+ |
| Pages | 25+ |
