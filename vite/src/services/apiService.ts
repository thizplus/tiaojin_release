// src/services/apiService.ts
import axios, { type AxiosInstance, type AxiosRequestConfig } from 'axios';
import { useAuthStore, forceLogout } from '@/stores/authStore';

// API Base URL
const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

/**
 * สร้าง axios instance สำหรับเรียกใช้ API
 */
const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 30000, // 30 seconds timeout
});

/**
 * Interceptor สำหรับเพิ่ม token ใน header ทุกครั้งที่มีการเรียก API
 */
apiClient.interceptors.request.use(
  (config) => {
    // เข้าถึง token จาก localStorage แบบถูกต้อง
    const authData = localStorage.getItem('auth-storage');
    if (authData) {
      try {
        const parsed = JSON.parse(authData);
        const token = parsed.state?.token;

        if (token && config.headers) {
          config.headers['Authorization'] = `Bearer ${token}`;
        }
      } catch (e) {
        console.error('Failed to parse auth data:', e);
      }
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

/**
 * Interceptor สำหรับจัดการ response และ error handling
 */
apiClient.interceptors.response.use(
  (response) => {
    return response;
  },
  async (error) => {
    // Log error สำหรับ debugging
    console.error('API Error:', {
      status: error.response?.status,
      statusText: error.response?.statusText,
      data: error.response?.data,
      config: error.config
    });

    // ถ้าเป็น error 401 (Unauthorized) หรือ 404 บน profile endpoint (user ถูกลบ)
    const isUnauthorized = error.response?.status === 401;
    const isProfileNotFound = error.response?.status === 404 &&
      error.config?.url?.includes('/users/profile');

    if (isUnauthorized || isProfileNotFound) {
      const authState = useAuthStore.getState();

      // ถ้ากำลัง logout อยู่แล้ว ไม่ต้องทำซ้ำ
      if (authState.isLoggingOut) {
        return Promise.reject(error);
      }

      // ตรวจสอบว่ามี token และ authenticated อยู่หรือไม่
      if (authState.token && authState.isAuthenticated) {
        // เช็คว่าไม่ใช่หน้า public ที่ไม่ต้อง auth
        const publicPaths = ['/', '/login', '/register', '/unauthorized'];
        const currentPath = window.location.pathname;

        if (!publicPaths.includes(currentPath)) {
          // forceLogout จะ cancel queries และ set isLoggingOut = true
          forceLogout();

          // ใช้ setTimeout เพื่อให้ state update เสร็จก่อน redirect
          setTimeout(() => {
            window.location.href = '/login';
          }, 100);
        }
      }
    }

    // จัดการ error response จาก backend
    if (error.response?.data) {
      const errorData = error.response.data;

      // สร้าง enhanced error object ที่มี backend error message
      const enhancedError = {
        ...errorData,
        status: error.response.status,
        statusText: error.response.statusText,
        originalError: error
      };

      return Promise.reject(enhancedError);
    }

    // ถ้าไม่มี response data (เช่น network error)
    return Promise.reject({
      success: false,
      message: 'เกิดข้อผิดพลาดในการเชื่อมต่อเครือข่าย',
      error: error.message,
      status: 0,
      originalError: error
    });
  }
);

/**
 * Service สำหรับเรียกใช้ API
 */
const apiService = {
  /**
   * ส่งคำขอ GET
   * @param url - URL ปลายทาง
   * @param params - พารามิเตอร์สำหรับ query string
   * @param config - ค่า config เพิ่มเติมสำหรับ axios
   */
  get: async <T>(url: string, params?: unknown, config?: AxiosRequestConfig): Promise<T> => {
    const response = await apiClient.get<T>(url, { params, ...config });
    return response.data;
  },

  /**
   * ส่งคำขอ POST
   * @param url - URL ปลายทาง
   * @param data - ข้อมูลที่จะส่งไปยัง API
   * @param config - ค่า config เพิ่มเติมสำหรับ axios
   */
  post: async <T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> => {
    const response = await apiClient.post<T>(url, data, config);
    return response.data;
  },

  /**
   * ส่งคำขอ PUT
   * @param url - URL ปลายทาง
   * @param data - ข้อมูลที่จะส่งไปยัง API
   * @param config - ค่า config เพิ่มเติมสำหรับ axios
   */
  put: async <T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> => {
    const response = await apiClient.put<T>(url, data, config);
    return response.data;
  },

  putFormData: async <T>(url: string, formData: FormData): Promise<T> => {
    const response = await apiClient.put<T>(url, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  },

  /**
   * ส่งคำขอ PATCH
   * @param url - URL ปลายทาง
   * @param data - ข้อมูลที่จะส่งไปยัง API
   * @param config - ค่า config เพิ่มเติมสำหรับ axios
   */
  patch: async <T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> => {
    const response = await apiClient.patch<T>(url, data, config);
    return response.data;
  },

  /**
   * ส่งคำขอ PATCH พร้อม FormData (สำหรับการอัพโหลดไฟล์)
   * @param url - URL ปลายทาง
   * @param formData - FormData ที่มีข้อมูลและไฟล์
   */
  patchFormData: async <T>(url: string, formData: FormData): Promise<T> => {
    const response = await apiClient.patch<T>(url, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  },

  /**
   * ส่งคำขอ DELETE
   * @param url - URL ปลายทาง
   * @param config - ค่า config เพิ่มเติมสำหรับ axios
   */
  delete: async <T>(url: string, config?: AxiosRequestConfig): Promise<T> => {
    const response = await apiClient.delete<T>(url, config);
    return response.data;
  },

  /**
   * ส่งคำขอ DELETE พร้อม data ใน request body
   * @param url - URL ปลายทาง
   * @param data - ข้อมูลที่จะส่งไปยัง API
   * @param config - ค่า config เพิ่มเติมสำหรับ axios
   */
  deleteWithData: async <T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> => {
    const response = await apiClient.delete<T>(url, { data, ...config });
    return response.data;
  },

  /**
   * ส่งคำขอแบบกำหนดเองได้เต็มรูปแบบ
   * @param config - ค่า config สำหรับ axios request
   */
  request: async <T>(config: AxiosRequestConfig): Promise<T> => {
    const response = await apiClient.request<T>(config);
    return response.data;
  },

  /**
   * ส่งคำขอ POST พร้อม FormData (สำหรับการอัพโหลดไฟล์)
   * @param url - URL ปลายทาง
   * @param formData - FormData ที่มีข้อมูลและไฟล์
   */
  postFile: async <T>(url: string, formData: FormData): Promise<T> => {
    const response = await apiClient.post<T>(url, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  },

  /**
   * ส่งคำขอ Upload ไฟล์
   * @param url - URL ปลายทาง
   * @param formData - FormData ที่มีไฟล์ที่จะอัปโหลด
   * @param onProgress - callback สำหรับติดตามความคืบหน้าในการอัปโหลด
   */
  upload: async <T>(
    url: string,
    formData: FormData,
    onProgress?: (percentage: number) => void
  ): Promise<T> => {
    const response = await apiClient.post<T>(url, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
      onUploadProgress: (progressEvent) => {
        if (onProgress && progressEvent.total) {
          const percentage = Math.round((progressEvent.loaded * 100) / progressEvent.total);
          onProgress(percentage);
        }
      },
    });
    return response.data;
  },
};

export default apiService;
