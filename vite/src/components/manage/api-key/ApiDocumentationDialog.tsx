import { memo } from "react";
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Badge } from "@/components/ui/badge";
import { Code, CheckCircle2, XCircle, AlertTriangle, Globe, Lock } from "lucide-react";

interface ApiDocumentationDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export const ApiDocumentationDialog = memo(function ApiDocumentationDialog({
  open,
  onOpenChange,
}: ApiDocumentationDialogProps) {
  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent side="right" className="w-full sm:max-w-2xl overflow-y-auto p-6">
        <SheetHeader className="p-0">
          <SheetTitle className="flex items-center gap-2">
            <Globe className="h-5 w-5" />
            การใช้งาน Public API
          </SheetTitle>
          <SheetDescription>
            คู่มือการใช้งาน API สำหรับดึงข้อมูลทัวร์
          </SheetDescription>
        </SheetHeader>

        <Tabs defaultValue="tours-list" className="w-full mt-6">
          <TabsList className="grid w-full grid-cols-4">
            <TabsTrigger value="tours-list" className="text-xs">ดึงรายการทัวร์</TabsTrigger>
            <TabsTrigger value="tour-detail" className="text-xs">ดึงทัวร์เดียว</TabsTrigger>
            <TabsTrigger value="security" className="text-xs">ความปลอดภัย</TabsTrigger>
            <TabsTrigger value="examples" className="text-xs">ตัวอย่างโค้ด</TabsTrigger>
          </TabsList>

          {/* Tab 1: ดึงรายการทัวร์ */}
          <TabsContent value="tours-list" className="space-y-4">
            <div className="space-y-3">
              <div className="flex items-center gap-2">
                <Badge variant="default" className="font-mono">GET</Badge>
                <code className="text-sm">/api/v1/public/tours</code>
              </div>

              <div>
                <h4 className="font-semibold mb-2">Headers:</h4>
                <div className="bg-muted p-3 rounded-lg space-y-1">
                  <code className="text-sm block">X-API-Key: tiaojin_your_api_key_here</code>
                  <div className="text-xs text-muted-foreground">หรือ</div>
                  <code className="text-sm block">Authorization: ApiKey tiaojin_your_api_key_here</code>
                </div>
              </div>

              <div>
                <h4 className="font-semibold mb-2">Query Parameters:</h4>
                <div className="space-y-2">
                  <div className="flex gap-2">
                    <Badge variant="outline">offset</Badge>
                    <span className="text-sm">เริ่มต้นที่รายการที่ (default: 0)</span>
                  </div>
                  <div className="flex gap-2">
                    <Badge variant="outline">limit</Badge>
                    <span className="text-sm">จำนวนรายการต่อหน้า (default: 20, max: 100)</span>
                  </div>
                  <div className="flex gap-2">
                    <Badge variant="outline">search</Badge>
                    <span className="text-sm">ค้นหาจาก title หรือ code</span>
                  </div>
                  <div className="flex gap-2">
                    <Badge variant="outline">city_id</Badge>
                    <span className="text-sm">กรองตาม city</span>
                  </div>
                </div>
              </div>

              <div>
                <h4 className="font-semibold mb-2">ตัวอย่าง Request:</h4>
                <div className="bg-muted p-3 rounded-lg overflow-x-auto">
                  <pre className="text-xs">
{`curl -X GET "http://localhost:8080/api/v1/public/tours?limit=10&search=ญี่ปุ่น" \\
  -H "X-API-Key: tiaojin_your_api_key_here"`}
                  </pre>
                </div>
              </div>

              <div>
                <h4 className="font-semibold mb-2">Response Example:</h4>
                <div className="bg-muted p-3 rounded-lg overflow-x-auto max-h-64">
                  <pre className="text-xs">
{`{
  "success": true,
  "message": "Tours retrieved successfully",
  "data": {
    "tours": [
      {
        "id": "123e4567-e89b-12d3-a456-426614174000",
        "title": "ทัวร์ญี่ปุ่น โตเกียว โอซาก้า 5 วัน 3 คืน",
        "code": "JP001",
        "city_id": "550e8400-e29b-41d4-a716-446655440000",
        "city_title": "โตเกียว",
        "detail": "เที่ยวญี่ปุ่นแบบจัดเต็ม...",
        "thumbnail": "https://example.com/image.jpg",
        "status": "publish",
        "created_at": "2025-09-01T10:00:00Z",
        "min_start_date": "2025-10-15",
        "available_periods": [
          {
            "id": "456e7890-e89b-12d3-a456-426614174001",
            "start_date": "2025-10-15",
            "end_date": "2025-10-19",
            "seat": 30,
            "booking": 15,
            "available_seats": 12,
            "airline": "Thai Airways",
            "airport": "สุวรรณภูมิ",
            "sector1": "BKK-NRT TG640 10:00-18:00",
            "sector2": "KIX-BKK TG623 18:00-22:00",
            "status": "publish"
          }
        ]
      }
    ],
    "pagination": {
      "total": 45,
      "offset": 0,
      "limit": 10,
      "has_more": true
    }
  }
}`}
                  </pre>
                </div>
              </div>
            </div>
          </TabsContent>

          {/* Tab 2: ดึงทัวร์เดียว */}
          <TabsContent value="tour-detail" className="space-y-4">
            <div className="space-y-3">
              <div className="flex items-center gap-2">
                <Badge variant="default" className="font-mono">GET</Badge>
                <code className="text-sm">/api/v1/public/tours/:id</code>
              </div>

              <div>
                <h4 className="font-semibold mb-2">Headers:</h4>
                <div className="bg-muted p-3 rounded-lg">
                  <code className="text-sm">X-API-Key: tiaojin_your_api_key_here</code>
                </div>
              </div>

              <div>
                <h4 className="font-semibold mb-2">ตัวอย่าง Request:</h4>
                <div className="bg-muted p-3 rounded-lg overflow-x-auto">
                  <pre className="text-xs">
{`curl -X GET "http://localhost:8080/api/v1/public/tours/123e4567-e89b-12d3-a456-426614174000" \\
  -H "X-API-Key: tiaojin_your_api_key_here"`}
                  </pre>
                </div>
              </div>

              <div>
                <h4 className="font-semibold mb-2">Response Example:</h4>
                <div className="bg-muted p-3 rounded-lg overflow-x-auto max-h-64">
                  <pre className="text-xs">
{`{
  "success": true,
  "message": "Tour retrieved successfully",
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "title": "ทัวร์ญี่ปุ่น โตเกียว โอซาก้า 5 วัน 3 คืน",
    "code": "JP001",
    "city_id": "550e8400-e29b-41d4-a716-446655440000",
    "city_title": "โตเกียว",
    "detail": "เที่ยวญี่ปุ่นแบบจัดเต็ม...",
    "thumbnail": "https://example.com/image.jpg",
    "status": "publish",
    "created_at": "2025-09-01T10:00:00Z",
    "min_start_date": "2025-10-15",
    "available_periods": [...]
  }
}`}
                  </pre>
                </div>
              </div>
            </div>
          </TabsContent>

          {/* Tab 3: ความปลอดภัย */}
          <TabsContent value="security" className="space-y-4">
            <div className="space-y-4">
              <div>
                <h4 className="font-semibold mb-3 flex items-center gap-2">
                  <XCircle className="h-4 w-4 text-destructive" />
                  ข้อมูลที่ซ่อนไว้ (ไม่แสดงใน Public API)
                </h4>
                <div className="space-y-2 ml-6">
                  <div className="flex items-center gap-2 text-sm text-muted-foreground">
                    <XCircle className="h-3 w-3 text-destructive" />
                    <span>ราคามัดจำ (deposit, deposit_type)</span>
                  </div>
                  <div className="flex items-center gap-2 text-sm text-muted-foreground">
                    <XCircle className="h-3 w-3 text-destructive" />
                    <span>รายการราคา (rate_items)</span>
                  </div>
                  <div className="flex items-center gap-2 text-sm text-muted-foreground">
                    <XCircle className="h-3 w-3 text-destructive" />
                    <span>ค่าคอมมิชชั่น (commission_items)</span>
                  </div>
                  <div className="flex items-center gap-2 text-sm text-muted-foreground">
                    <XCircle className="h-3 w-3 text-destructive" />
                    <span>ราคาห้องพัก (room_types)</span>
                  </div>
                  <div className="flex items-center gap-2 text-sm text-muted-foreground">
                    <XCircle className="h-3 w-3 text-destructive" />
                    <span>บริการเสริม (additional_services)</span>
                  </div>
                </div>
              </div>

              <div>
                <h4 className="font-semibold mb-3 flex items-center gap-2">
                  <CheckCircle2 className="h-4 w-4 text-chart-2" />
                  ข้อมูลที่แสดง
                </h4>
                <div className="space-y-2 ml-6">
                  <div className="flex items-center gap-2 text-sm">
                    <CheckCircle2 className="h-3 w-3 text-chart-2" />
                    <span>ข้อมูลทัวร์ทั่วไป (title, code, detail, thumbnail)</span>
                  </div>
                  <div className="flex items-center gap-2 text-sm">
                    <CheckCircle2 className="h-3 w-3 text-chart-2" />
                    <span>จำนวนที่นั่ง (seat, booking, available_seats)</span>
                  </div>
                  <div className="flex items-center gap-2 text-sm">
                    <CheckCircle2 className="h-3 w-3 text-chart-2" />
                    <span>วันเดินทาง (start_date, end_date)</span>
                  </div>
                  <div className="flex items-center gap-2 text-sm">
                    <CheckCircle2 className="h-3 w-3 text-chart-2" />
                    <span>สายการบิน, สนามบิน</span>
                  </div>
                  <div className="flex items-center gap-2 text-sm">
                    <CheckCircle2 className="h-3 w-3 text-chart-2" />
                    <span>Sector เที่ยวบิน</span>
                  </div>
                  <div className="flex items-center gap-2 text-sm">
                    <CheckCircle2 className="h-3 w-3 text-chart-2" />
                    <span>ไฟล์เอกสาร</span>
                  </div>
                </div>
              </div>

              <div className="bg-accent/50 border border-accent p-4 rounded-lg">
                <h4 className="font-semibold mb-2 flex items-center gap-2">
                  <Lock className="h-4 w-4" />
                  การจำกัดสิทธิ์
                </h4>
                <ul className="space-y-1 ml-6 text-sm">
                  <li>• แสดงเฉพาะทัวร์และ period ที่ status = "publish"</li>
                  <li>• ต้องมี API Key ที่ active และไม่หมดอายุ</li>
                  <li>• ตรวจสอบ IP whitelist (ถ้ามีการกำหนด)</li>
                </ul>
              </div>

              <div>
                <h4 className="font-semibold mb-3 flex items-center gap-2">
                  <AlertTriangle className="h-4 w-4 text-destructive" />
                  Error Responses
                </h4>
                <div className="space-y-3">
                  <div>
                    <Badge variant="destructive">401 Unauthorized</Badge>
                    <div className="bg-muted p-3 rounded-lg mt-2 text-xs">
                      <pre>{`{
  "success": false,
  "message": "Invalid or missing API key",
  "error": null
}`}</pre>
                    </div>
                  </div>
                  <div>
                    <Badge variant="destructive">404 Not Found</Badge>
                    <div className="bg-muted p-3 rounded-lg mt-2 text-xs">
                      <pre>{`{
  "success": false,
  "message": "Tour not found",
  "error": null
}`}</pre>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </TabsContent>

          {/* Tab 4: ตัวอย่างโค้ด */}
          <TabsContent value="examples" className="space-y-4">
            <div className="space-y-4">
              <div>
                <h4 className="font-semibold mb-2 flex items-center gap-2">
                  <Code className="h-4 w-4" />
                  JavaScript / Fetch
                </h4>
                <div className="bg-muted p-3 rounded-lg overflow-x-auto">
                  <pre className="text-xs">
{`const apiKey = 'tiaojin_your_api_key_here';

// ดึงทัวร์ทั้งหมด
fetch('http://localhost:8080/api/v1/public/tours?limit=20', {
  headers: {
    'X-API-Key': apiKey
  }
})
  .then(res => res.json())
  .then(data => console.log(data));

// ดึงทัวร์เฉพาะรายการ
fetch('http://localhost:8080/api/v1/public/tours/123e4567-e89b-12d3-a456-426614174000', {
  headers: {
    'X-API-Key': apiKey
  }
})
  .then(res => res.json())
  .then(data => console.log(data));`}
                  </pre>
                </div>
              </div>

              <div>
                <h4 className="font-semibold mb-2 flex items-center gap-2">
                  <Code className="h-4 w-4" />
                  PHP / cURL
                </h4>
                <div className="bg-muted p-3 rounded-lg overflow-x-auto">
                  <pre className="text-xs">
{`<?php
$apiKey = 'tiaojin_your_api_key_here';

$ch = curl_init('http://localhost:8080/api/v1/public/tours?limit=20');
curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
curl_setopt($ch, CURLOPT_HTTPHEADER, [
    'X-API-Key: ' . $apiKey
]);

$response = curl_exec($ch);
$data = json_decode($response, true);

print_r($data);
curl_close($ch);
?>`}
                  </pre>
                </div>
              </div>

              <div>
                <h4 className="font-semibold mb-2 flex items-center gap-2">
                  <Code className="h-4 w-4" />
                  Python / Requests
                </h4>
                <div className="bg-muted p-3 rounded-lg overflow-x-auto">
                  <pre className="text-xs">
{`import requests

api_key = 'tiaojin_your_api_key_here'
headers = {'X-API-Key': api_key}

# ดึงทัวร์ทั้งหมด
response = requests.get(
    'http://localhost:8080/api/v1/public/tours',
    headers=headers,
    params={'limit': 20, 'search': 'ญี่ปุ่น'}
)
print(response.json())`}
                  </pre>
                </div>
              </div>
            </div>
          </TabsContent>
        </Tabs>

        <div className="bg-primary/10 border border-primary/20 p-4 rounded-lg mt-4 mb-6">
          <p className="text-sm">
            <strong>สรุป:</strong> Agent สามารถสร้าง API Key ผ่านระบบ จากนั้นนำ key ไปใช้เรียก Public API
            เพื่อดึงข้อมูลทัวร์ที่เปิดขายโดยไม่แสดงข้อมูลราคา เหมาะสำหรับแสดงบนเว็บไซต์ภายนอก
            หรือ integration กับระบบอื่นๆ
          </p>
        </div>
      </SheetContent>
    </Sheet>
  );
});