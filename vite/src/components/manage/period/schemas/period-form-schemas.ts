import { z } from 'zod';

export const flightSectorSchema = z.object({
  date: z.string().min(1, 'กรุณาระบุวันที่'),
  flight: z.string().min(1, 'กรุณาระบุรหัสเที่ยวบิน'),
  route: z.string().min(1, 'กรุณาระบุเส้นทาง'),
});

export const rateItemSchema = z.object({
  title: z.string().min(1, 'กรุณาระบุชื่อราคา'),
  regular: z.number().min(0, 'ราคาต้องไม่ติดลบ'),
  sale: z.number().min(0, 'ราคาต้องไม่ติดลบ'),
  exclude: z.boolean().optional().default(false),
});

export const commissionItemSchema = z.object({
  regular: z.number().min(0, 'คอมมิชชั่นต้องไม่ติดลบ'),
  sale: z.number().min(0, 'คอมมิชชั่นต้องไม่ติดลบ'),
});

export const roomTypeSchema = z.object({
  title: z.string().min(1, 'กรุณาระบุชื่อประเภทห้อง'),
  price: z.number(),
});

export const additionalServiceSchema = z.object({
  title: z.string().min(1, 'กรุณาระบุชื่อบริการ'),
  price: z.number().min(0, 'ราคาต้องไม่ติดลบ'),
  rate_type: z.enum(['per_unit', 'per_booking']),
  unit_name: z.string().min(1, 'กรุณาระบุหน่วย'),
});

export const periodFileSchema = z.object({
  name: z.string().min(1, 'กรุณาระบุชื่อไฟล์').max(255, 'ชื่อไฟล์ต้องไม่เกิน 255 ตัวอักษร'),
  type: z.string().min(1, 'กรุณาระบุประเภทไฟล์').max(100, 'ประเภทไฟล์ต้องไม่เกิน 100 ตัวอักษร'),
  url: z.string().url('กรุณาระบุ URL ที่ถูกต้อง').max(500, 'URL ต้องไม่เกิน 500 ตัวอักษร'),
});

export const periodFormSchema = z.object({
  tour_id: z.string().min(1, 'กรุณาเลือกทัวร์'),
  start_date: z.string().min(1, 'กรุณาระบุวันที่เริ่มต้น'),
  end_date: z.string().min(1, 'กรุณาระบุวันที่สิ้นสุด'),
  seat: z.number().min(1, 'จำนวนที่นั่งต้องมากกว่า 0'),
  deposit: z.number().min(0, 'เงินมัดจำต้องไม่ติดลบ'),
  deposit_type: z.enum(['per_unit', 'percent']),
 // single_stay: z.number().min(0, 'ค่าที่พักเดี่ยวต้องไม่ติดลบ'),
  airline: z.string().min(1, 'กรุณาระบุสายการบิน'),
  airport: z.string().min(1, 'กรุณาระบุสนามบิน'),
  note: z.string().optional(),
  sector_1: z.array(flightSectorSchema).min(1, 'กรุณาเพิ่มเที่ยวบินไปอย่างน้อย 1 เที่ยว'),
  sector_2: z.array(flightSectorSchema).min(1, 'กรุณาเพิ่มเที่ยวบินกลับอย่างน้อย 1 เที่ยว'),
  rate_type: z.enum(['regular', 'sale']),
  rate_items: z.array(rateItemSchema).min(1, 'กรุณาเพิ่มรายการราคาอย่างน้อย 1 รายการ'),
  commission_type: z.enum(['regular', 'sale']),
  commission_items: z.array(commissionItemSchema).min(1, 'กรุณาเพิ่มคอมมิชชั่นอย่างน้อย 1 รายการ'),
  room_types: z.array(roomTypeSchema).optional(),
  additional_services: z.array(additionalServiceSchema).optional(),
  files: z.array(periodFileSchema).optional(),
  status: z.enum(['draft', 'publish']).default('draft'),
  service_fee_per_pax: z.number().min(0, 'ค่าบริการต้องไม่ติดลบ').default(0),

  // Expiration Override - สำหรับทัวร์ไฟไหม้
  expiration_override: z.boolean().optional().default(false),
  expiration_unit: z.enum(['minutes', 'hours', 'days']).optional().default('hours'),
  expiration_value: z.number().min(1).max(10080).optional(), // max 7 days in minutes
  expiration_message: z.string().max(255).optional(),
});

export type PeriodFormData = z.infer<typeof periodFormSchema>;