// src/services/index.ts

// นำเข้า services ทั้งหมด
import apiService from './apiService';
import bookingService from './bookingService';

// Export ทั้งหมดเพื่อให้นำเข้าได้ง่ายขึ้น
export {
  apiService,
  bookingService,
};

// Export default เป็น object ที่รวม services ทั้งหมด
const services = {
  api: apiService,
  booking: bookingService,
};

export default services;