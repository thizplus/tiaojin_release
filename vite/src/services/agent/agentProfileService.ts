import API from '@/constants/api';
import apiService from '@/services/apiService';
import type {

  CreateAgentProfileRequest,
  UpdateAgentProfileRequest,

  ProfileCheckApiResponse,
  ProfileApiResponse,
  ProfileCreateApiResponse,
  ProfileUpdateApiResponse,
  ApiResponse
} from '@/types';

class AgentProfileService {
  /**
   * ตรวจสอบว่า Agent มี profile หรือยัง
   */
  async checkProfile(): Promise<ProfileCheckApiResponse> {
    try {
      const response = await apiService.get<ProfileCheckApiResponse>(
        API.AGENT_PROFILE.CHECK_PROFILE
      );
      return response;
    } catch (error) {
      console.error('Failed to check profile:', error);
      throw error;
    }
  }

  /**
   * ดูข้อมูล profile ของ Agent
   */
  async getProfile(): Promise<ProfileApiResponse> {
    try {
      const response = await apiService.get<ProfileApiResponse>(
        API.AGENT_PROFILE.PROFILE
      );
      return response;
    } catch (error) {
      console.error('Failed to get profile:', error);
      throw error;
    }
  }

  /**
   * สร้าง profile ใหม่สำหรับ Agent
   */
  async createProfile(data: CreateAgentProfileRequest): Promise<ProfileCreateApiResponse> {
    try {
      const response = await apiService.post<ProfileCreateApiResponse>(
        API.AGENT_PROFILE.CREATE_PROFILE,
        data
      );
      return response;
    } catch (error) {
      console.error('Failed to create profile:', error);
      throw error;
    }
  }

  /**
   * แก้ไขข้อมูล profile ของ Agent
   */
  async updateProfile(data: UpdateAgentProfileRequest): Promise<ProfileUpdateApiResponse> {
    try {
      const response = await apiService.put<ProfileUpdateApiResponse>(
        API.AGENT_PROFILE.PROFILE,
        data
      );
      return response;
    } catch (error) {
      console.error('Failed to update profile:', error);
      throw error;
    }
  }

  /**
   * อัพโหลดไฟล์สำหรับ profile (logo, atta, dbd, pp20)
   */
  async uploadFile(
    file: File,
    category: 'agent_profiles' | 'agent_licenses',
    fileType: 'logo' | 'atta' | 'dbd' | 'pp20'
  ): Promise<ApiResponse<{ url: string; fileName: string }>> {
    try {
      const formData = new FormData();
      formData.append('file', file);
      formData.append('category', category);
      formData.append('file_type', fileType);

      const response = await apiService.upload<ApiResponse<{ url: string; fileName: string }>>(
        API.FILE.UPLOAD,
        formData
      );
      return response;
    } catch (error) {
      console.error(`Failed to upload ${fileType} file:`, error);
      throw error;
    }
  }

  /**
   * อัพโหลด logo สำหรับ profile (legacy method for backward compatibility)
   */
  async uploadLogo(file: File): Promise<ApiResponse<{ url: string; fileName: string }>> {
    return this.uploadFile(file, 'agent_profiles', 'logo');
  }

  /**
   * อัพโหลดใบอนุญาต ATTA
   */
  async uploadAtta(file: File): Promise<ApiResponse<{ url: string; fileName: string }>> {
    return this.uploadFile(file, 'agent_licenses', 'atta');
  }

  /**
   * อัพโหลดทะเบียน DBD
   */
  async uploadDbd(file: File): Promise<ApiResponse<{ url: string; fileName: string }>> {
    return this.uploadFile(file, 'agent_licenses', 'dbd');
  }

  /**
   * อัพโหลดใบทะเบียน PP20
   */
  async uploadPp20(file: File): Promise<ApiResponse<{ url: string; fileName: string }>> {
    return this.uploadFile(file, 'agent_licenses', 'pp20');
  }

  /**
   * ตรวจสอบความถูกต้องของข้อมูลก่อนส่ง
   */
  validateProfileData(data: CreateAgentProfileRequest | UpdateAgentProfileRequest): {
    isValid: boolean;
    errors: Record<string, string>;
  } {
    const errors: Record<string, string> = {};

    // ตรวจสอบชื่อบริษัท
    if (!data.name?.trim()) {
      errors.name = 'กรุณาระบุชื่อบริษัท';
    } else if (data.name.length > 200) {
      errors.name = 'ชื่อบริษัทต้องไม่เกิน 200 ตัวอักษร';
    }

    // ตรวจสอบเลขผู้เสียภาษี (สำหรับ create เท่านั้น)
    if ('tax_id' in data) {
      if (!data.tax_id?.trim()) {
        errors.tax_id = 'กรุณาระบุเลขประจำตัวผู้เสียภาษี';
      } else if (!/^\d{13}$/.test(data.tax_id)) {
        errors.tax_id = 'เลขประจำตัวผู้เสียภาษีต้องเป็นตัวเลข 13 หลัก';
      }
    }

    // ตรวจสอบเบอร์โทรศัพท์
    if (!data.phone?.trim()) {
      errors.phone = 'กรุณาระบุเบอร์โทรศัพท์';
    }

    // ตรวจสอบที่อยู่
    if (!data.address?.trim()) {
      errors.address = 'กรุณาระบุที่อยู่';
    }

    // ตรวจสอบจังหวัด
    if (!data.province?.trim()) {
      errors.province = 'กรุณาระบุจังหวัด';
    }

    // ตรวจสอบรหัสไปรษณีย์
    if (!data.zipcode?.trim()) {
      errors.zipcode = 'กรุณาระบุรหัสไปรษณีย์';
    } else if (!/^\d{5}$/.test(data.zipcode)) {
      errors.zipcode = 'รหัสไปรษณีย์ต้องเป็นตัวเลข 5 หลัก';
    }

    return {
      isValid: Object.keys(errors).length === 0,
      errors
    };
  }

  /**
   * ตรวจสอบความถูกต้องของไฟล์
   */
  validateFile(
    file: File,
    fileType: 'logo' | 'atta' | 'dbd' | 'pp20'
  ): {
    isValid: boolean;
    error?: string;
  } {
    const settings = {
      logo: {
        maxSize: 5 * 1024 * 1024, // 5MB
        allowedTypes: ['image/jpeg', 'image/png', 'image/webp'],
        typeName: 'โลโก้บริษัท'
      },
      atta: {
        maxSize: 10 * 1024 * 1024, // 10MB
        allowedTypes: ['application/pdf', 'image/jpeg', 'image/png'],
        typeName: 'ใบอนุญาต ATTA'
      },
      dbd: {
        maxSize: 10 * 1024 * 1024, // 10MB
        allowedTypes: ['application/pdf', 'image/jpeg', 'image/png'],
        typeName: 'ทะเบียน DBD'
      },
      pp20: {
        maxSize: 10 * 1024 * 1024, // 10MB
        allowedTypes: ['application/pdf', 'image/jpeg', 'image/png'],
        typeName: 'ใบทะเบียน PP20'
      }
    };

    const setting = settings[fileType];

    // ตรวจสอบขนาดไฟล์
    if (file.size > setting.maxSize) {
      return {
        isValid: false,
        error: `${setting.typeName}มีขนาดใหญ่เกินไป (สูงสุด ${Math.round(setting.maxSize / 1024 / 1024)}MB)`
      };
    }

    // ตรวจสอบประเภทไฟล์
    if (!setting.allowedTypes.includes(file.type)) {
      return {
        isValid: false,
        error: `ประเภทไฟล์ไม่ถูกต้องสำหรับ${setting.typeName}`
      };
    }

    return { isValid: true };
  }
}

const agentProfileService = new AgentProfileService();
export default agentProfileService;