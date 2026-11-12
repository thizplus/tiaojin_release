import { useState, useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { useTours } from "@/hooks/useTours";
import type { Tour, UpdateTourRequest } from "@/types";
import { FileUpload } from "@/components/ui/file-upload";

import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
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
import { Button } from "@/components/ui/button";
import { Alert, AlertDescription } from "@/components/ui/alert";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { CitySelect } from "@/components/ui/city-select";

const editTourSchema = z.object({
  title: z.string().min(1, "กรุณาระบุชื่อทัวร์"),
  city_id: z.string().min(1, "กรุณาเลือกจังหวัด"),
  detail: z.string().min(1, "กรุณาระบุรายละเอียดทัวร์"),
  thumbnail: z.string().min(1, "กรุณาระบุรูปภาพ"),
  file_pdf: z.string().optional(),
  file_word: z.string().optional(),
  status: z.enum(["draft", "publish"]),
});

type EditTourFormData = z.infer<typeof editTourSchema>;

interface EditTourDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  tour: Tour | null;
  onSuccess?: () => void;
}

export function EditTourDialog({
  open,
  onOpenChange,
  tour,
  onSuccess,
}: EditTourDialogProps) {
  const { updateTour, isUpdateLoading, refetchToursList } = useTours();
  const [error, setError] = useState<string>("");

  const form = useForm<EditTourFormData>({
    resolver: zodResolver(editTourSchema),
    defaultValues: {
      title: "",
      city_id: "",
      detail: "",
      thumbnail: "",
      file_pdf: "",
      file_word: "",
      status: "draft",
    },
  });

  useEffect(() => {
    if (tour && open) {
      form.reset({
        title: tour.title,
        city_id: tour.city_id,
        detail: tour.detail,
        thumbnail: tour.thumbnail,
        file_pdf: tour.file_pdf || "",
        file_word: tour.file_word || "",
        status: tour.status,
      });
      setError("");
    }
  }, [tour, open, form]);

  const onSubmit = async (data: EditTourFormData) => {
    if (!tour) return;

    try {
      setError("");
      const payload: UpdateTourRequest = {
        title: data.title,
        city_id: data.city_id,
        detail: data.detail,
        thumbnail: data.thumbnail,
        file_pdf: data.file_pdf || undefined,
        file_word: data.file_word || undefined,
        status: data.status,
      };

      await updateTour(tour.id, payload);
      refetchToursList();
      onOpenChange(false);
      onSuccess?.();
    } catch (err) {
      setError(err instanceof Error ? err.message : "เกิดข้อผิดพลาดในการแก้ไขทัวร์");
    }
  };

  const handleClose = () => {
    onOpenChange(false);
    form.reset();
    setError("");
  };

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent className="sm:max-w-[600px] max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>แก้ไขข้อมูลทัวร์</DialogTitle>
          <DialogDescription>
            แก้ไขข้อมูลทัวร์ในระบบ กรุณากรอกข้อมูลให้ครบถ้วน
          </DialogDescription>
        </DialogHeader>

        {tour && (
          <div className="py-2">
            <p className="text-sm text-muted-foreground">
              รหัสทัวร์: <span className="font-mono">{tour.code}</span>
            </p>
          </div>
        )}

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            {error && (
              <Alert variant="destructive">
                <AlertDescription>{error}</AlertDescription>
              </Alert>
            )}

            <FormField
              control={form.control}
              name="title"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>ชื่อทัวร์ *</FormLabel>
                  <FormControl>
                    <Input
                      placeholder="เช่น เที่ยวกรุงเทพ 3 วัน 2 คืน"
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

<FormField
              control={form.control}
              name="detail"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>รายละเอียดทัวร์ *</FormLabel>
                  <FormControl>
                    <Textarea
                      placeholder="รายละเอียดและโปรแกรมการท่องเที่ยว..."
                      rows={4}
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="city_id"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>จังหวัด *</FormLabel>
                  <FormControl>
                    <CitySelect
                      value={field.value}
                      onValueChange={field.onChange}
                      placeholder="เลือกจังหวัดที่ทัวร์จะไป..."
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="thumbnail"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>รูปภาพทัวร์ *</FormLabel>
                  <FormControl>
                    <FileUpload
                      structuredOptions={{
                        category: 'tours',
                        file_type: 'image'
                      }}
                      accept="image/*"
                      maxSize={5 * 1024 * 1024} // 5MB
                      allowedTypes={['image/jpeg', 'image/png', 'image/webp']}
                      value={field.value}
                      onChange={(url) => {
                        console.log('FileUpload onChange called with URL:', url);
                        field.onChange(url || "");
                      }}
                      onUploadSuccess={(result) => {
                        console.log('Upload success:', result);
                      }}
                      onUploadError={(error) => {
                        console.error('Upload error:', error);
                        setError(`เกิดข้อผิดพลาดในการอัพโหลดไฟล์: ${error.message}`);
                      }}
                      variant="dropzone"
                      size="md"
                      placeholder="เลือกรูปภาพทัวร์หรือลากมาวาง"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

           

            <FormField
              control={form.control}
              name="file_pdf"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>ไฟล์ PDF (ถ้ามี)</FormLabel>
                  <FormControl>
                    <FileUpload
                      structuredOptions={{
                        category: 'tours',
                        file_type: 'document'
                      }}
                      accept=".pdf,application/pdf"
                      maxSize={50 * 1024 * 1024} // 50MB
                      allowedTypes={['application/pdf']}
                      value={field.value}
                      onChange={(url) => field.onChange(url || "")}
                      variant="dropzone"
                      size="md"
                      placeholder="เลือกไฟล์ PDF หรือลากมาวาง"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="file_word"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>ไฟล์ Word (ถ้ามี)</FormLabel>
                  <FormControl>
                    <FileUpload
                      structuredOptions={{
                        category: 'tours',
                        file_type: 'document'
                      }}
                      accept=".doc,.docx,application/msword,application/vnd.openxmlformats-officedocument.wordprocessingml.document"
                      maxSize={50 * 1024 * 1024} // 50MB
                      allowedTypes={['application/msword', 'application/vnd.openxmlformats-officedocument.wordprocessingml.document']}
                      value={field.value}
                      onChange={(url) => field.onChange(url || "")}
                      variant="dropzone"
                      size="md"
                      placeholder="เลือกไฟล์ Word หรือลากมาวาง"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="status"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>สถานะ</FormLabel>
                  <Select onValueChange={field.onChange} value={field.value}>
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder="เลือกสถานะ" />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      <SelectItem value="publish">เผยแพร่</SelectItem>
                      <SelectItem value="draft">ฉบับร่าง</SelectItem>
                    </SelectContent>
                  </Select>
                  <FormMessage />
                </FormItem>
              )}
            />

            <DialogFooter>
              <Button type="button" variant="outline" onClick={handleClose}>
                ยกเลิก
              </Button>
              <Button type="submit" disabled={isUpdateLoading}>
                {isUpdateLoading ? "กำลังบันทึก..." : "บันทึกการแก้ไข"}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}