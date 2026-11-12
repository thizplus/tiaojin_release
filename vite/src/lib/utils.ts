import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

/**
 * Storage URL utilities สำหรับ private storage
 * ใช้ proxy ผ่าน backend API เพื่อเข้าถึงไฟล์ใน private S3 bucket
 */

// ดึง API base URL จาก env (ตัด /api/v1 ออกเพื่อเอาแค่ domain)
const getApiBaseUrl = (): string => {
  const apiUrl = import.meta.env.VITE_API_URL || ''
  // ถ้ามี VITE_API_URL ให้ตัด /api/v1 ออก เช่น https://xxx.host/api/v1 -> https://xxx.host
  if (apiUrl) {
    return apiUrl.replace(/\/api\/v1\/?$/, '')
  }
  return ''
}

/**
 * แปลง storage path เป็น proxy URL สำหรับแสดงผล (images, PDFs)
 * @param path - path ของไฟล์ เช่น "profiles/logo/abc.jpg" หรือ full URL หรือ proxy URL
 * @returns proxy URL ที่ใช้งานได้
 */
export function getStorageUrl(path: string | undefined | null): string {
  if (!path) return ''

  const baseUrl = getApiBaseUrl()

  // ถ้าเป็น proxy URL อยู่แล้ว (/api/v1/storage/...)
  if (path.startsWith('/api/v1/storage/')) {
    // ถ้ามี baseUrl ให้ใช้ absolute URL
    return baseUrl ? `${baseUrl}${path}` : path
  }

  // ถ้าเป็น full URL ที่เริ่มต้นด้วย http แล้ว ให้แปลงเป็น proxy path
  if (path.startsWith('http://') || path.startsWith('https://')) {
    // ดึงเฉพาะ path จาก URL
    // เช่น https://bucket.e2.example.com/profiles/logo/abc.jpg -> profiles/logo/abc.jpg
    try {
      const url = new URL(path)
      const storagePath = url.pathname.startsWith('/') ? url.pathname.slice(1) : url.pathname
      const proxyPath = `/api/v1/storage/${storagePath}`
      return baseUrl ? `${baseUrl}${proxyPath}` : proxyPath
    } catch {
      // ถ้า parse URL ไม่ได้ ให้ใช้ path ตรงๆ
      const proxyPath = `/api/v1/storage/${path}`
      return baseUrl ? `${baseUrl}${proxyPath}` : proxyPath
    }
  }

  // ถ้าเป็น path ปกติ ให้เติม proxy prefix
  const cleanPath = path.startsWith('/') ? path.slice(1) : path
  const proxyPath = `/api/v1/storage/${cleanPath}`
  return baseUrl ? `${baseUrl}${proxyPath}` : proxyPath
}

/**
 * ตรวจสอบว่าไฟล์เป็น image หรือไม่
 */
export function isImageFile(path: string | undefined | null): boolean {
  if (!path) return false
  const ext = path.toLowerCase().split('.').pop()
  return ['jpg', 'jpeg', 'png', 'gif', 'webp', 'svg'].includes(ext || '')
}

/**
 * ตรวจสอบว่าไฟล์เป็น PDF หรือไม่
 */
export function isPdfFile(path: string | undefined | null): boolean {
  if (!path) return false
  return path.toLowerCase().endsWith('.pdf')
}

/**
 * ตรวจสอบว่าไฟล์เป็น document (doc, docx, xls, xlsx) หรือไม่
 */
export function isDocumentFile(path: string | undefined | null): boolean {
  if (!path) return false
  const ext = path.toLowerCase().split('.').pop()
  return ['doc', 'docx', 'xls', 'xlsx'].includes(ext || '')
}
