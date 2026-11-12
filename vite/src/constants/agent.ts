// Agent Profile Constants

// Validation Messages
export const AGENT_PROFILE_MESSAGES = {
  // Success Messages
  PROFILE_CREATED: 'โปรไฟล์บริษัทได้รับการสร้างเรียบร้อยแล้ว',
  PROFILE_UPDATED: 'ข้อมูลโปรไฟล์บริษัทได้รับการอัพเดทเรียบร้อยแล้ว',
  PROFILE_RETRIEVED: 'ดึงข้อมูลโปรไฟล์สำเร็จ',
  FILE_UPLOADED: 'ไฟล์ได้รับการอัพโหลดเรียบร้อยแล้ว',

  // Error Messages
  PROFILE_EXISTS: 'คุณมีโปรไฟล์บริษัทอยู่แล้ว',
  PROFILE_NOT_FOUND: 'ไม่พบข้อมูลโปรไฟล์บริษัท',
  VALIDATION_FAILED: 'ข้อมูลไม่ถูกต้อง กรุณาตรวจสอบอีกครั้ง',
  CREATE_FAILED: 'ไม่สามารถสร้างโปรไฟล์ได้',
  UPDATE_FAILED: 'ไม่สามารถแก้ไขโปรไฟล์ได้',
  FETCH_FAILED: 'ไม่สามารถโหลดข้อมูลโปรไฟล์ได้',
  FILE_UPLOAD_FAILED: 'ไม่สามารถอัพโหลดไฟล์ได้',

  // Loading Messages
  CREATING: 'กำลังสร้างโปรไฟล์...',
  UPDATING: 'กำลังอัพเดทข้อมูล...',
  LOADING: 'กำลังโหลดข้อมูล...',
  UPLOADING: 'กำลังอัพโหลดไฟล์...',
} as const;

// Field Labels
export const AGENT_PROFILE_LABELS = {
  name: 'ชื่อบริษัท',
  tax_id: 'เลขประจำตัวผู้เสียภาษี',
  phone: 'เบอร์โทรศัพท์',
  address: 'ที่อยู่',
  province: 'จังหวัด',
  zipcode: 'รหัสไปรษณีย์',
  atta: 'ใบอนุญาต ATTA',
  dbd: 'ทะเบียนพาณิชย์อิเล็กทรอนิกส์ (DBD)',
  pp20: 'ใบทะเบียนภาษีมูลค่าเพิ่ม (PP20)',
  logo: 'โลโก้บริษัท',
} as const;

// Field Placeholders
export const AGENT_PROFILE_PLACEHOLDERS = {
  name: 'เช่น บริษัท ABC ทัวร์ จำกัด',
  tax_id: '1234567890123',
  phone: '02-123-4567',
  address: '123 ถนนสุขุมวิท แขวงคลองเตย เขตคลองเตย',
  province: 'กรุงเทพมหานคร',
  zipcode: '10110',
  atta: 'เลือกไฟล์ใบอนุญาต ATTA',
  dbd: 'เลือกไฟล์ทะเบียน DBD',
  pp20: 'เลือกไฟล์ใบทะเบียน PP20',
  logo: 'เลือกไฟล์โลโก้บริษัท',
} as const;

// Validation Rules
export const AGENT_PROFILE_VALIDATION = {
  name: {
    required: 'กรุณาระบุชื่อบริษัท',
    maxLength: 200,
    maxLengthMessage: 'ชื่อบริษัทต้องไม่เกิน 200 ตัวอักษร',
  },
  tax_id: {
    required: 'กรุณาระบุเลขประจำตัวผู้เสียภาษี',
    pattern: /^\d{13}$/,
    patternMessage: 'เลขประจำตัวผู้เสียภาษีต้องเป็นตัวเลข 13 หลัก',
  },
  phone: {
    required: 'กรุณาระบุเบอร์โทรศัพท์',
  },
  address: {
    required: 'กรุณาระบุที่อยู่',
  },
  province: {
    required: 'กรุณาระบุจังหวัด',
  },
  zipcode: {
    required: 'กรุณาระบุรหัสไปรษณีย์',
    pattern: /^\d{5}$/,
    patternMessage: 'รหัสไปรษณีย์ต้องเป็นตัวเลข 5 หลัก',
  },
} as const;


// Page States
export const PROFILE_STATES = {
  LOADING: 'loading',
  NO_PROFILE: 'no_profile',
  HAS_PROFILE: 'has_profile',
  CREATING: 'creating',
  EDITING: 'editing',
  ERROR: 'error',
} as const;

export type ProfileState = typeof PROFILE_STATES[keyof typeof PROFILE_STATES];

// Section Names for Navigation/UI
export const PROFILE_SECTIONS = {
  OVERVIEW: 'overview',
  BASIC_INFO: 'basic_info',
  CONTACT: 'contact',
  ADDRESS: 'address',
  DOCUMENTS: 'documents',
  LOGO: 'logo',
} as const;

export type ProfileSection = typeof PROFILE_SECTIONS[keyof typeof PROFILE_SECTIONS];