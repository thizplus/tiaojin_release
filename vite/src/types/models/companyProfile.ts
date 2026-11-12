export interface CompanyProfile {
  id: string;
  user_id: string;
  name: string;
  tax_id: string;
  phone: string;
  address: string;
  province: string;
  zipcode: string;
  atta: string;
  dbd: string;
  pp20: string;
  logo: string | null;
  created_at: string;
  updated_at: string;
}

export interface PendingCompanyProfile extends CompanyProfile {
  // เฉพาะ profiles ของ users ที่มี status = pending
}