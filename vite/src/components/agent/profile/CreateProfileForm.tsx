import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Save, Loader2, Building2, FileText } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { FileUpload } from "@/components/ui/file-upload";
import type { CreateAgentProfileRequest } from "@/types";

// Embedded validation constants
const VALIDATION = {
  name: {
    required: "กรุณากรอกชื่อบริษัท",
    maxLength: 100,
    maxLengthMessage: "ชื่อบริษัทต้องไม่เกิน 100 ตัวอักษร"
  },
  tax_id: {
    required: "กรุณากรอกเลขประจำตัวผู้เสียภาษี",
    pattern: /^\d{13}$/,
    patternMessage: "เลขประจำตัวผู้เสียภาษีต้องเป็นตัวเลข 13 หลัก"
  },
  phone: {
    required: "กรุณากรอกเบอร์โทรศัพท์"
  },
  address: {
    required: "กรุณากรอกที่อยู่"
  },
  province: {
    required: "กรุณากรอกจังหวัด"
  },
  zipcode: {
    required: "กรุณากรอกรหัสไปรษณีย์",
    pattern: /^\d{5}$/,
    patternMessage: "รหัสไปรษณีย์ต้องเป็นตัวเลข 5 หลัก"
  }
};

const LABELS = {
  name: 'ชื่อบริษัท',
  tax_id: 'เลขประจำตัวผู้เสียภาษี',
  phone: 'เบอร์โทรศัพท์',
  province: 'จังหวัด',
  address: 'ที่อยู่',
  zipcode: 'รหัสไปรษณีย์',
  atta: 'ใบอนุญาต ATTA',
  dbd: 'ทะเบียนพาณิชย์อิเล็กทรอนิกส์ (DBD)',
  pp20: 'ใบทะเบียนภาษีมูลค่าเพิ่ม (PP20)',
  logo: 'โลโก้บริษัท'
};

const PLACEHOLDERS = {
  name: 'เช่น บริษัท ท่องเที่ยวไทย จำกัด',
  tax_id: '1234567890123',
  phone: '02-123-4567',
  province: 'กรุงเทพมหานคร',
  address: '123 ถนนสุขุมวิท แขวงคลองเตย เขตคลองเตย',
  zipcode: '10110'
};

// Validation Schema
const createProfileSchema = z.object({
  name: z.string()
    .min(1, VALIDATION.name.required)
    .max(VALIDATION.name.maxLength, VALIDATION.name.maxLengthMessage),
  tax_id: z.string()
    .min(1, VALIDATION.tax_id.required)
    .regex(VALIDATION.tax_id.pattern, VALIDATION.tax_id.patternMessage),
  phone: z.string()
    .min(1, VALIDATION.phone.required),
  address: z.string()
    .min(1, VALIDATION.address.required),
  province: z.string()
    .min(1, VALIDATION.province.required),
  zipcode: z.string()
    .min(1, VALIDATION.zipcode.required)
    .regex(VALIDATION.zipcode.pattern, VALIDATION.zipcode.patternMessage),
  atta: z.string().optional(),
  dbd: z.string().optional(),
  pp20: z.string().optional(),
  logo: z.string().optional(),
});

type CreateProfileFormData = z.infer<typeof createProfileSchema>;

interface CreateProfileFormProps {
  onSubmit: (data: CreateAgentProfileRequest) => Promise<void>;
  isLoading?: boolean;
  className?: string;
}

export function CreateProfileForm({
  onSubmit,
  isLoading = false,
  className = ""
}: CreateProfileFormProps) {
  const [error, setError] = useState<string>("");

  const form = useForm<CreateProfileFormData>({
    resolver: zodResolver(createProfileSchema),
    defaultValues: {
      name: "",
      tax_id: "",
      phone: "",
      address: "",
      province: "",
      zipcode: "",
      atta: "",
      dbd: "",
      pp20: "",
      logo: "",
    },
  });

  const handleSubmit = async (data: CreateProfileFormData) => {
    try {
      setError("");
      await onSubmit(data);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "เกิดข้อผิดพลาดในการสร้างโปรไฟล์";
      setError(errorMessage);
    }
  };

  const FormSection = ({ title, children, className: sectionClassName = "", icon }: {
    title: string;
    children: React.ReactNode;
    className?: string;
    icon?: React.ReactNode;
  }) => (
    <div className={`space-y-8 ${sectionClassName}`}>
      <div className="flex items-center gap-4 pb-6 border-b border-border">
        {icon && (
          <div className="w-10 h-10 rounded-lg bg-secondary border border-border flex items-center justify-center">
            {icon}
          </div>
        )}
        <h3 className="text-xl font-bold text-foreground">
          {title}
        </h3>
      </div>
      {children}
    </div>
  );

  return (
    <div className={`bg-card border border-border rounded-xl p-10 shadow-lg ${className}`}>
     

      <Form {...form}>
        <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-8">
          {error && (
            <div className="bg-muted border-2 border-destructive p-6 rounded-lg mb-8">
              <div className="flex items-center gap-4">
                <div className="w-12 h-12 rounded-lg bg-destructive text-destructive-foreground flex items-center justify-center">
                  <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                </div>
                <div>
                  <h4 className="text-lg font-semibold text-foreground mb-1">เกิดข้อผิดพลาด</h4>
                  <p className="text-base text-muted-foreground">{error}</p>
                </div>
              </div>
            </div>
          )}

          {/* Basic Information */}
          <FormSection
            title="ข้อมูลพื้นฐาน"
            icon={<Building2 className="w-5 h-5 text-secondary-foreground" />}
          >
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{LABELS.name} *</FormLabel>
                  <FormControl>
                    <Input placeholder={PLACEHOLDERS.name} {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <FormField
                control={form.control}
                name="tax_id"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{LABELS.tax_id} *</FormLabel>
                    <FormControl>
                      <Input
                        placeholder={PLACEHOLDERS.tax_id}
                        maxLength={13}
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="phone"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{LABELS.phone} *</FormLabel>
                    <FormControl>
                      <Input placeholder={PLACEHOLDERS.phone} {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>
          </FormSection>

          {/* Address Information */}
          <FormSection
            title="ที่อยู่"
            icon={<svg className="w-5 h-5 text-secondary-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" /><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" /></svg>}
          >
            <FormField
              control={form.control}
              name="address"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{LABELS.address} *</FormLabel>
                  <FormControl>
                    <Textarea
                      placeholder={PLACEHOLDERS.address}
                      rows={3}
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <FormField
                control={form.control}
                name="province"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{LABELS.province} *</FormLabel>
                    <FormControl>
                      <Input placeholder={PLACEHOLDERS.province} {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="zipcode"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{LABELS.zipcode} *</FormLabel>
                    <FormControl>
                      <Input
                        placeholder={PLACEHOLDERS.zipcode}
                        maxLength={5}
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>
          </FormSection>

          {/* Documents Upload */}
          <FormSection
            title="เอกสารใบอนุญาต"
            icon={<FileText className="w-5 h-5 text-secondary-foreground" />}
          >
            <div className="grid gap-6">
              {/* ATTA License */}
              <FormField
                control={form.control}
                name="atta"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{LABELS.atta}</FormLabel>
                    <FormControl>
                      <FileUpload
                        structuredOptions={{
                          category: 'profiles',
                          file_type: 'document'
                        }}
                        accept=".pdf,.jpg,.jpeg,.png"
                        maxSize={5 * 1024 * 1024} // 5MB
                        allowedTypes={['application/pdf', 'image/jpeg', 'image/png']}
                        value={field.value}
                        onChange={(url) => field.onChange(url || "")}
                        variant="dropzone"
                        size="sm"
                        placeholder="เลือกไฟล์ใบอนุญาต ATTA หรือลากมาวาง"
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {/* DBD License */}
              <FormField
                control={form.control}
                name="dbd"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{LABELS.dbd}</FormLabel>
                    <FormControl>
                      <FileUpload
                        structuredOptions={{
                          category: 'profiles',
                          file_type: 'document'
                        }}
                        accept=".pdf,.jpg,.jpeg,.png"
                        maxSize={5 * 1024 * 1024} // 5MB
                        allowedTypes={['application/pdf', 'image/jpeg', 'image/png']}
                        value={field.value}
                        onChange={(url) => field.onChange(url || "")}
                        variant="dropzone"
                        size="sm"
                        placeholder="เลือกไฟล์ใบอนุญาต DBD หรือลากมาวาง"
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {/* PP20 License */}
              <FormField
                control={form.control}
                name="pp20"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{LABELS.pp20}</FormLabel>
                    <FormControl>
                      <FileUpload
                        structuredOptions={{
                          category: 'profiles',
                          file_type: 'document'
                        }}
                        accept=".pdf,.jpg,.jpeg,.png"
                        maxSize={5 * 1024 * 1024} // 5MB
                        allowedTypes={['application/pdf', 'image/jpeg', 'image/png']}
                        value={field.value}
                        onChange={(url) => field.onChange(url || "")}
                        variant="dropzone"
                        size="sm"
                        placeholder="เลือกไฟล์ใบอนุญาต PP.20 หรือลากมาวาง"
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>
          </FormSection>

          {/* Logo Upload */}
          <FormSection
            title="โลโก้บริษัท"
            icon={<svg className="w-5 h-5 text-secondary-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" /></svg>}
          >
            <FormField
              control={form.control}
              name="logo"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{LABELS.logo}</FormLabel>
                  <FormControl>
                    <FileUpload
                      structuredOptions={{
                        category: 'profiles',
                        file_type: 'image'
                      }}
                      accept="image/*"
                      maxSize={2 * 1024 * 1024} // 2MB
                      allowedTypes={['image/jpeg', 'image/png', 'image/webp']}
                      value={field.value}
                      onChange={(url) => field.onChange(url || "")}
                      variant="dropzone"
                      size="md"
                      placeholder="เลือกโลโก้บริษัทหรือลากมาวาง"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
          </FormSection>

          {/* Submit Button */}
          <div className="flex justify-center pt-10 border-t border-border">
            <Button
              type="submit"
              disabled={isLoading}
              size="lg"
              className="px-16 py-4 text-lg font-bold rounded-lg min-w-[280px]"
            >
              {isLoading ? (
                <>
                  <Loader2 className="h-6 w-6 mr-4 animate-spin" />
                  กำลังสร้างโปรไฟล์...
                </>
              ) : (
                <>
                  <Save className="h-6 w-6 mr-4" />
                  สร้างโปรไฟล์บริษัท
                </>
              )}
            </Button>
          </div>

          {/* Helper Text */}
          <div className="text-center text-base mt-6 text-muted-foreground">
            <p>หลังจากสร้างโปรไฟล์แล้ว คุณสามารถแก้ไขข้อมูลได้ทุกเมื่อ</p>
          </div>
        </form>
      </Form>
    </div>
  );
}