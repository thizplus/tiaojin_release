# System Settings API Documentation

## Overview

ระบบ System Settings ใช้จัดการการตั้งค่าต่างๆ ของระบบ โดยเฉพาะ:
1. **Expiration Tiers** - กำหนดระยะเวลาหมดอายุการจองตามความใกล้ของวันเริ่มทัวร์
2. **Agent Booking Limits** - จำกัดจำนวนการจองของ Agent ต่อ Period

---

## 1. Expiration Tiers

### แนวคิด

Expiration Tiers คือระบบกำหนด **ระยะเวลาหมดอายุการจอง** แบบแบ่งระดับ:
- ทัวร์ **ใกล้จะเริ่ม** → ให้เวลาจ่ายเงิน **น้อย** (เพราะต้องรีบยืนยันที่นั่ง)
- ทัวร์ **ยังอีกนาน** → ให้เวลาจ่ายเงิน **นาน** (ลูกค้ามีเวลาตัดสินใจ)

### Flow การทำงาน

```
Agent สร้าง Booking
       ↓
ระบบคำนวณ: ทัวร์เริ่มอีกกี่วัน? (DaysUntilTour)
       ↓
หา Tier ที่ match (เรียงตาม SortOrder จากน้อยไปมาก)
       ↓
คำนวณ expires_at = created_at + ExpireTime ของ Tier นั้น
       ↓
ถ้าไม่จ่ายเงินภายใน expires_at → Booking หมดอายุ (ถูก cancel อัตโนมัติ)
```

### ตัวอย่าง Tiers

| Tier | เงื่อนไข (threshold_days) | หมดอายุใน | ตัวอย่าง |
|------|---------------------------|-----------|----------|
| Critical | ทัวร์เริ่มภายใน 3 วัน | 30 นาที | จองวันนี้ ทัวร์ออกพรุ่งนี้ |
| Urgent | ทัวร์เริ่มภายใน 7 วัน | 6 ชั่วโมง | จองวันจันทร์ ทัวร์ออกวันศุกร์ |
| Normal | ทัวร์เริ่มภายใน 14 วัน | 24 ชั่วโมง | จองล่วงหน้า 10 วัน |
| Standard | ทัวร์เริ่มภายใน 30 วัน | 48 ชั่วโมง | จองล่วงหน้า 3 สัปดาห์ |
| Default | ทัวร์อื่นๆ ทั้งหมด | 72 ชั่วโมง | จองล่วงหน้ามากกว่า 30 วัน |

### Data Model

```go
type ExpirationTier struct {
    Name           string  `json:"name"`            // ชื่อ tier เช่น "Critical", "Urgent"
    Description    string  `json:"description"`     // คำอธิบาย
    ThresholdDays  *int    `json:"threshold_days"`  // ถ้าทัวร์เริ่มภายในกี่วัน (null = default tier)
    ExpireMinutes  *int    `json:"expire_minutes"`  // หมดอายุในกี่นาที (ใช้อันใดอันหนึ่ง)
    ExpireHours    *int    `json:"expire_hours"`    // หมดอายุในกี่ชั่วโมง
    ExpireDays     *int    `json:"expire_days"`     // หมดอายุในกี่วัน
    SortOrder      int     `json:"sort_order"`      // ลำดับการประมวลผล (น้อย = ตรวจก่อน)
}
```

### API Endpoints

#### GET /api/v1/super/configs/expiration-tiers
ดึงข้อมูล Expiration Tiers ปัจจุบัน

**Response:**
```json
{
  "success": true,
  "data": {
    "tiers": [
      {
        "name": "Critical",
        "description": "ทัวร์ที่ใกล้จะเริ่มมาก ต้องยืนยันด่วน",
        "threshold_days": 3,
        "expire_minutes": 30,
        "expire_hours": null,
        "expire_days": null,
        "sort_order": 0,
        "example_message": "หมดอายุใน 30 นาที"
      },
      {
        "name": "Urgent",
        "description": "ทัวร์ที่ใกล้จะเริ่ม",
        "threshold_days": 7,
        "expire_minutes": null,
        "expire_hours": 6,
        "expire_days": null,
        "sort_order": 1,
        "example_message": "หมดอายุใน 6 ชั่วโมง"
      },
      {
        "name": "Normal",
        "description": "ทัวร์ปกติ",
        "threshold_days": 14,
        "expire_minutes": null,
        "expire_hours": 24,
        "expire_days": null,
        "sort_order": 2,
        "example_message": "หมดอายุใน 24 ชั่วโมง"
      },
      {
        "name": "Standard",
        "description": "ทัวร์ล่วงหน้า",
        "threshold_days": 30,
        "expire_minutes": null,
        "expire_hours": null,
        "expire_days": 2,
        "sort_order": 3,
        "example_message": "หมดอายุใน 2 วัน"
      },
      {
        "name": "Default",
        "description": "ทัวร์ล่วงหน้ามาก",
        "threshold_days": null,
        "expire_minutes": null,
        "expire_hours": null,
        "expire_days": 3,
        "sort_order": 4,
        "example_message": "หมดอายุใน 3 วัน"
      }
    ],
    "metadata": {
      "version": "1.0",
      "last_updated": "2024-01-15T10:30:00Z",
      "updated_by": "admin@example.com"
    }
  }
}
```

#### PUT /api/v1/super/configs/expiration-tiers
อัพเดท Expiration Tiers

**Request Body:**
```json
{
  "tiers": [
    {
      "name": "Critical",
      "description": "ทัวร์ที่ใกล้จะเริ่มมาก ต้องยืนยันด่วน",
      "threshold_days": 3,
      "expire_minutes": 30,
      "sort_order": 0
    },
    {
      "name": "Urgent",
      "description": "ทัวร์ที่ใกล้จะเริ่ม",
      "threshold_days": 7,
      "expire_hours": 6,
      "sort_order": 1
    },
    {
      "name": "Normal",
      "description": "ทัวร์ปกติ",
      "threshold_days": 14,
      "expire_hours": 24,
      "sort_order": 2
    },
    {
      "name": "Standard",
      "description": "ทัวร์ล่วงหน้า",
      "threshold_days": 30,
      "expire_days": 2,
      "sort_order": 3
    },
    {
      "name": "Default",
      "description": "ทัวร์ล่วงหน้ามาก (default tier)",
      "threshold_days": null,
      "expire_days": 3,
      "sort_order": 4
    }
  ]
}
```

**Validation Rules:**
- `name`: required, ไม่ซ้ำกัน
- `threshold_days`: ต้องเป็น positive หรือ null (สำหรับ default tier)
- ต้องระบุ **เพียงหนึ่งอย่าง** จาก: `expire_minutes`, `expire_hours`, `expire_days`
- ค่า expire ต้องเป็น positive
- ต้องมี **อย่างน้อย 1 tier** ที่ `threshold_days = null` (default tier)

---

## 2. Agent Booking Limits

### แนวคิด

จำกัดจำนวนการจองของ Agent แต่ละรายต่อ Period เพื่อ:
- ป้องกัน Agent รายใดรายหนึ่งกว้านซื้อที่นั่งทั้งหมด
- กระจายโอกาสให้ Agent อื่นๆ

### Data Model

```go
type AgentBookingLimitConfig struct {
    GlobalSettings       AgentBookingGlobalSettings    `json:"global_settings"`
    PeriodSpecificLimits []AgentBookingPeriodLimit     `json:"period_specific_limits"`
    AgentGlobalBypasses  []string                      `json:"agent_global_bypasses"`
    Metadata             AgentBookingLimitMetadata     `json:"metadata"`
}

type AgentBookingGlobalSettings struct {
    Enabled                bool     `json:"enabled"`                  // เปิด/ปิดระบบ
    DefaultLimitPerPeriod  int      `json:"default_limit_per_period"` // จำกัดกี่ที่นั่งต่อ period
    CountAllStatuses       []string `json:"count_all_statuses"`       // นับ status ไหนบ้าง
}

type AgentBookingPeriodLimit struct {
    PeriodID        string   `json:"period_id"`        // UUID ของ period
    Limit           int      `json:"limit"`            // จำกัดกี่ที่นั่ง
    BypassedAgents  []string `json:"bypassed_agents"`  // Agent ที่ไม่ถูกจำกัดใน period นี้
}
```

### API Endpoints

#### GET /api/v1/super/configs/agent-booking-limits
ดึงข้อมูล Agent Booking Limits ปัจจุบัน

**Response:**
```json
{
  "success": true,
  "data": {
    "global_settings": {
      "enabled": true,
      "default_limit_per_period": 10,
      "count_all_statuses": ["pending", "confirmed", "paid"]
    },
    "period_specific_limits": [
      {
        "period_id": "550e8400-e29b-41d4-a716-446655440000",
        "limit": 5,
        "bypassed_agents": ["agent-123", "agent-456"]
      },
      {
        "period_id": "660e8400-e29b-41d4-a716-446655440001",
        "limit": 20,
        "bypassed_agents": []
      }
    ],
    "agent_global_bypasses": [
      "vip-agent-001",
      "partner-agent-002"
    ],
    "metadata": {
      "version": "1.0",
      "last_updated": "2024-01-15T10:30:00Z",
      "updated_by": "admin@example.com"
    }
  }
}
```

#### PUT /api/v1/super/configs/agent-booking-limits
อัพเดท Agent Booking Limits

**Request Body:**
```json
{
  "global_settings": {
    "enabled": true,
    "default_limit_per_period": 10,
    "count_all_statuses": ["pending", "confirmed", "paid"]
  },
  "period_specific_limits": [
    {
      "period_id": "550e8400-e29b-41d4-a716-446655440000",
      "limit": 5,
      "bypassed_agents": ["agent-123"]
    }
  ],
  "agent_global_bypasses": [
    "vip-agent-001"
  ]
}
```

**Field Descriptions:**

| Field | Description |
|-------|-------------|
| `enabled` | เปิด/ปิดระบบจำกัดการจอง |
| `default_limit_per_period` | จำกัดการจองต่อ Agent ต่อ Period (default) |
| `count_all_statuses` | นับ booking ที่มี status เหล่านี้ เช่น pending, confirmed |
| `period_specific_limits` | กำหนด limit เฉพาะบาง period (override default) |
| `agent_global_bypasses` | Agent IDs ที่ไม่ถูกจำกัด (VIP agents) |

---

## 3. System Configs (General)

### API Endpoints

#### GET /api/v1/super/configs
ดึงรายการ System Configs ทั้งหมด

**Query Parameters:**
- `config_type`: filter by type (general, booking, payment, notification, expiration)
- `is_active`: filter by active status
- `search`: search by config_key
- `offset`: pagination offset
- `limit`: pagination limit

**Response:**
```json
{
  "success": true,
  "data": {
    "configs": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "config_key": "booking.expiration_tiers",
        "config_value": { "tiers": [...], "metadata": {...} },
        "config_type": "expiration",
        "description": "Expiration tier configuration",
        "is_active": true,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-15T10:30:00Z"
      }
    ],
    "pagination": {
      "total": 5,
      "offset": 0,
      "limit": 10
    }
  }
}
```

#### POST /api/v1/super/configs
สร้าง System Config ใหม่

**Request Body:**
```json
{
  "config_key": "notification.email_enabled",
  "config_value": {
    "enabled": true,
    "smtp_host": "smtp.example.com",
    "smtp_port": 587
  },
  "config_type": "notification",
  "description": "Email notification settings",
  "is_active": true
}
```

---

## 4. Expiration Statistics

#### GET /api/v1/super/configs/expiration-stats
ดึงสถิติการหมดอายุของ Booking

**Response:**
```json
{
  "success": true,
  "data": {
    "total_bookings_with_expiration": 150,
    "expired_bookings": 12,
    "near_expiring_bookings": 8,
    "bookings_by_tier": {
      "Critical": 5,
      "Urgent": 15,
      "Normal": 45,
      "Standard": 50,
      "Default": 35
    },
    "average_expiration_time": "18 hours",
    "expiration_trends": [
      {
        "date": "2024-01-14",
        "created": 25,
        "expired": 3,
        "near_expiring": 5
      },
      {
        "date": "2024-01-15",
        "created": 30,
        "expired": 2,
        "near_expiring": 4
      }
    ]
  }
}
```

---

## 5. Process Expired Bookings

#### POST /api/v1/super/configs/process-expired
ประมวลผล Booking ที่หมดอายุ (Cancel อัตโนมัติ)

**Request Body:**
```json
{
  "dry_run": true
}
```

- `dry_run: true` → แค่ดูว่ามี booking ไหนหมดอายุ (ไม่ cancel จริง)
- `dry_run: false` → Cancel booking ที่หมดอายุจริง

**Response:**
```json
{
  "success": true,
  "data": {
    "processed_at": "2024-01-15T10:30:00Z",
    "dry_run": true,
    "total_checked": 150,
    "expired_found": 12,
    "processed_count": 12,
    "failed_count": 0,
    "processed_bookings": [
      {
        "booking_id": "550e8400-e29b-41d4-a716-446655440000",
        "invoice_id": "INV-2024-001",
        "customer_name": "John Doe",
        "expires_at": "2024-01-15T08:00:00Z",
        "processed_at": "2024-01-15T10:30:00Z",
        "action": "canceled",
        "success": true
      }
    ],
    "errors": []
  }
}
```

---

## Database Schema

### system_configs table

```sql
CREATE TABLE system_configs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    config_key VARCHAR(100) UNIQUE NOT NULL,
    config_value JSONB NOT NULL,
    config_type VARCHAR(50) NOT NULL DEFAULT 'general',
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Config types: general, booking, payment, notification, expiration
```

### Initial Data (Seed)

```sql
INSERT INTO system_configs (config_key, config_value, config_type, description) VALUES
(
  'booking.expiration_tiers',
  '{
    "tiers": [
      {"name": "Critical", "description": "ทัวร์ที่ใกล้จะเริ่มมาก", "threshold_days": 3, "expire_minutes": 30, "sort_order": 0},
      {"name": "Urgent", "description": "ทัวร์ที่ใกล้จะเริ่ม", "threshold_days": 7, "expire_hours": 6, "sort_order": 1},
      {"name": "Normal", "description": "ทัวร์ปกติ", "threshold_days": 14, "expire_hours": 24, "sort_order": 2},
      {"name": "Standard", "description": "ทัวร์ล่วงหน้า", "threshold_days": 30, "expire_days": 2, "sort_order": 3},
      {"name": "Default", "description": "ทัวร์ล่วงหน้ามาก", "threshold_days": null, "expire_days": 3, "sort_order": 4}
    ],
    "metadata": {"version": "1.0", "last_updated": "2024-01-01T00:00:00Z", "updated_by": "system"}
  }',
  'expiration',
  'Booking expiration tier configuration'
),
(
  'booking.agent_period_limits',
  '{
    "global_settings": {
      "enabled": false,
      "default_limit_per_period": 10,
      "count_all_statuses": ["pending", "confirmed", "paid"]
    },
    "period_specific_limits": [],
    "agent_global_bypasses": [],
    "metadata": {"version": "1.0", "last_updated": "2024-01-01T00:00:00Z", "updated_by": "system"}
  }',
  'booking',
  'Agent booking limit configuration per period'
);
```

---

## Frontend Page

**Path:** `/super/system-settings`

### Tabs:
| Tab | Component | Description |
|-----|-----------|-------------|
| การหมดอายุ | `ExpirationTiersCard` | ดู/แก้ไข Expiration Tiers |
| จำกัดการจอง | `BookingLimitsCard` | ตั้งค่า Agent Booking Limits |
| การตั้งค่า | `ConfigManagementCard` | จัดการ System Configs ทั่วไป |
| บำรุงรักษา | `ProcessExpiredCard` | ประมวลผล Booking ที่หมดอายุ |

---

## 6. Config Examples (ตัวอย่างการตั้งค่า)

### รูปแบบ Config Key

```
category.setting_name
```

| Category | ใช้สำหรับ |
|----------|----------|
| `booking` | การจอง, การหมดอายุ, การจำกัด |
| `payment` | การชำระเงิน, ธนาคาร |
| `notification` | การแจ้งเตือน (Email, LINE, SMS) |
| `security` | ความปลอดภัย, Rate Limit |
| `integration` | เชื่อมต่อบริการภายนอก |
| `general` | ทั่วไป |

---

### ตัวอย่างที่ 1: Expiration Tiers (การหมดอายุการจอง)

| Field | ค่า |
|-------|-----|
| **Config Key** | `booking.expiration_tiers` |
| **ประเภท** | การหมดอายุ (expiration) |
| **คำอธิบาย** | กำหนดระยะเวลาหมดอายุการจองตามความใกล้ของวันเริ่มทัวร์ |

**JSON:**
```json
{
  "tiers": [
    {
      "name": "Critical",
      "description": "ทัวร์ที่ใกล้จะเริ่มมาก ต้องยืนยันด่วน",
      "threshold_days": 3,
      "expire_minutes": 30,
      "sort_order": 0
    },
    {
      "name": "Urgent",
      "description": "ทัวร์ที่ใกล้จะเริ่ม",
      "threshold_days": 7,
      "expire_hours": 6,
      "sort_order": 1
    },
    {
      "name": "Normal",
      "description": "ทัวร์ปกติ",
      "threshold_days": 14,
      "expire_hours": 24,
      "sort_order": 2
    },
    {
      "name": "Standard",
      "description": "ทัวร์ล่วงหน้า",
      "threshold_days": 30,
      "expire_days": 2,
      "sort_order": 3
    },
    {
      "name": "Default",
      "description": "ทัวร์ล่วงหน้ามาก (catch-all)",
      "threshold_days": null,
      "expire_days": 3,
      "sort_order": 4
    }
  ],
  "metadata": {
    "version": "1.0",
    "last_updated": "2024-01-15T10:30:00Z",
    "updated_by": "admin@company.com"
  }
}
```

**คำอธิบาย Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | ชื่อ tier (ต้องไม่ซ้ำกัน) |
| `description` | string | คำอธิบายสำหรับแสดงใน UI |
| `threshold_days` | int/null | ทัวร์เริ่มภายในกี่วัน (null = default tier) |
| `expire_minutes` | int/null | หมดอายุในกี่นาที |
| `expire_hours` | int/null | หมดอายุในกี่ชั่วโมง |
| `expire_days` | int/null | หมดอายุในกี่วัน |
| `sort_order` | int | ลำดับการประมวลผล (0 = ตรวจก่อน) |

**กฎสำคัญ:**
- ต้องระบุ **เพียง 1 อย่าง** จาก: `expire_minutes`, `expire_hours`, `expire_days`
- ต้องมี **อย่างน้อย 1 tier** ที่ `threshold_days = null` (default tier สำหรับจับทุกกรณี)
- `sort_order` น้อย = ตรวจก่อน (tier ที่เฉพาะเจาะจงกว่าควรมี sort_order น้อยกว่า)

---

### ตัวอย่างที่ 2: Agent Booking Limits (จำกัดการจอง)

| Field | ค่า |
|-------|-----|
| **Config Key** | `booking.agent_period_limits` |
| **ประเภท** | ทั่วไป (general) |
| **คำอธิบาย** | จำกัดจำนวนการจองของ Agent ต่อ Period |

**JSON:**
```json
{
  "global_settings": {
    "enabled": true,
    "default_limit_per_period": 10,
    "count_all_statuses": ["pending", "confirmed", "paid"]
  },
  "period_specific_limits": [
    {
      "period_id": "550e8400-e29b-41d4-a716-446655440000",
      "limit": 5,
      "bypassed_agents": ["agent-vip-001"]
    }
  ],
  "agent_global_bypasses": [
    "partner-agent-001",
    "vip-agent-002"
  ],
  "metadata": {
    "version": "1.0",
    "last_updated": "2024-01-15T10:30:00Z",
    "updated_by": "admin@company.com"
  }
}
```

---

### ตัวอย่างที่ 3: Payment Settings (การชำระเงิน)

| Field | ค่า |
|-------|-----|
| **Config Key** | `payment.bank_accounts` |
| **ประเภท** | การชำระเงิน (payment) |
| **คำอธิบาย** | บัญชีธนาคารสำหรับรับชำระเงิน |

**JSON:**
```json
{
  "enabled": true,
  "accounts": [
    {
      "bank_code": "KBANK",
      "bank_name": "ธนาคารกสิกรไทย",
      "account_number": "123-4-56789-0",
      "account_name": "บริษัท ทัวร์ จำกัด",
      "branch": "สาขาสยาม",
      "is_primary": true
    },
    {
      "bank_code": "SCB",
      "bank_name": "ธนาคารไทยพาณิชย์",
      "account_number": "987-6-54321-0",
      "account_name": "บริษัท ทัวร์ จำกัด",
      "branch": "สาขาเซ็นทรัล",
      "is_primary": false
    }
  ],
  "promptpay": {
    "enabled": true,
    "id": "0123456789012",
    "name": "บริษัท ทัวร์ จำกัด",
    "qr_code_url": "https://example.com/qr/promptpay.png"
  }
}
```

---

### ตัวอย่างที่ 4: Email Notification (แจ้งเตือนอีเมล)

| Field | ค่า |
|-------|-----|
| **Config Key** | `notification.email_settings` |
| **ประเภท** | การแจ้งเตือน (notification) |
| **คำอธิบาย** | การตั้งค่าส่งอีเมลแจ้งเตือน |

**JSON:**
```json
{
  "enabled": true,
  "smtp": {
    "host": "smtp.gmail.com",
    "port": 587,
    "secure": true,
    "user": "noreply@company.com",
    "password_secret_key": "SMTP_PASSWORD"
  },
  "from": {
    "name": "Tour Booking System",
    "email": "noreply@company.com"
  },
  "templates": {
    "booking_confirmation": {
      "enabled": true,
      "subject": "ยืนยันการจอง #{invoice_id}"
    },
    "payment_reminder": {
      "enabled": true,
      "subject": "แจ้งเตือน: กรุณาชำระเงินภายใน {hours} ชั่วโมง",
      "send_before_hours": [24, 6, 1]
    },
    "booking_expired": {
      "enabled": true,
      "subject": "การจองหมดอายุ #{invoice_id}"
    }
  }
}
```

---

### ตัวอย่างที่ 5: LINE Notification (แจ้งเตือน LINE)

| Field | ค่า |
|-------|-----|
| **Config Key** | `notification.line_settings` |
| **ประเภท** | การแจ้งเตือน (notification) |
| **คำอธิบาย** | การตั้งค่า LINE Notify |

**JSON:**
```json
{
  "enabled": true,
  "channel_access_token_secret_key": "LINE_CHANNEL_TOKEN",
  "notifications": {
    "new_booking": {
      "enabled": true,
      "notify_admin": true,
      "notify_agent": true
    },
    "payment_received": {
      "enabled": true,
      "notify_admin": true,
      "notify_agent": true
    },
    "booking_expired": {
      "enabled": true,
      "notify_admin": false,
      "notify_agent": true
    }
  },
  "admin_group_ids": [
    "Cxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  ]
}
```

---

### ตัวอย่างที่ 6: Rate Limit (จำกัด Request)

| Field | ค่า |
|-------|-----|
| **Config Key** | `security.rate_limit` |
| **ประเภท** | ความปลอดภัย (security) |
| **คำอธิบาย** | จำกัดจำนวน request ต่อนาที |

**JSON:**
```json
{
  "enabled": true,
  "default": {
    "max_requests_per_minute": 100,
    "max_requests_per_hour": 1000
  },
  "endpoints": {
    "/api/v1/auth/login": {
      "max_requests_per_minute": 10,
      "lockout_minutes": 15
    },
    "/api/v1/bookings": {
      "max_requests_per_minute": 30
    }
  },
  "whitelist_ips": [
    "192.168.1.0/24",
    "10.0.0.0/8"
  ],
  "blacklist_ips": []
}
```

---

### ตัวอย่างที่ 7: Company Info (ข้อมูลบริษัท)

| Field | ค่า |
|-------|-----|
| **Config Key** | `general.company_info` |
| **ประเภท** | ทั่วไป (general) |
| **คำอธิบาย** | ข้อมูลบริษัทสำหรับแสดงในเอกสาร |

**JSON:**
```json
{
  "name": {
    "th": "บริษัท ทัวร์ จำกัด",
    "en": "Tour Company Ltd."
  },
  "tax_id": "0123456789012",
  "address": {
    "th": "123 ถ.สุขุมวิท แขวงคลองเตย เขตคลองเตย กรุงเทพฯ 10110",
    "en": "123 Sukhumvit Rd., Klongtoey, Bangkok 10110"
  },
  "contact": {
    "phone": "02-123-4567",
    "mobile": "081-234-5678",
    "email": "info@company.com",
    "line_id": "@tourcompany",
    "website": "https://www.company.com"
  },
  "license": {
    "tat_license": "11/12345",
    "dbd_registration": "0123456789012"
  },
  "logo_url": "https://example.com/logo.png"
}
```

---

### ตัวอย่างที่ 8: Commission Rates (อัตราค่าคอมมิชชั่น)

| Field | ค่า |
|-------|-----|
| **Config Key** | `booking.commission_rates` |
| **ประเภท** | ทั่วไป (general) |
| **คำอธิบาย** | อัตราค่าคอมมิชชั่นสำหรับ Agent |

**JSON:**
```json
{
  "enabled": true,
  "default_rate_percent": 10,
  "rates_by_agent_tier": {
    "standard": 10,
    "gold": 12,
    "platinum": 15,
    "diamond": 18
  },
  "rates_by_tour_type": {
    "domestic": 8,
    "international": 12,
    "premium": 15
  },
  "minimum_commission": 100,
  "payment_schedule": "monthly",
  "payment_day": 15
}
```

---

## Summary

| Feature | Config Key | Purpose |
|---------|-----------|---------|
| Expiration Tiers | `booking.expiration_tiers` | กำหนดระยะเวลาหมดอายุการจอง |
| Booking Limits | `booking.agent_period_limits` | จำกัดการจองของ Agent ต่อ Period |
| Bank Accounts | `payment.bank_accounts` | บัญชีธนาคารสำหรับรับชำระ |
| Email Settings | `notification.email_settings` | การตั้งค่าส่งอีเมล |
| LINE Settings | `notification.line_settings` | การตั้งค่า LINE Notify |
| Rate Limit | `security.rate_limit` | จำกัดจำนวน request |
| Company Info | `general.company_info` | ข้อมูลบริษัท |
| Commission | `booking.commission_rates` | อัตราค่าคอมมิชชั่น |
