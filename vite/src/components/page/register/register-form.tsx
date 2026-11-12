import { useState } from "react"
import { useNavigate } from "react-router-dom"
import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { useAuth } from "@/hooks/useAuth"
import { toast } from "sonner"
import { Eye, EyeOff } from "lucide-react"
import type { RegisterRequest, UserRole } from "@/types"

interface SimpleRegisterRequest {
  email: string;
  password: string;
  role: UserRole;
}

export function RegisterForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
  const navigate = useNavigate()
  const { register, isRegisterLoading, isAuthenticated } = useAuth()

  // Form state
  const [formData, setFormData] = useState({
    email: "",
    password: "",
    confirmPassword: ""
  })
  const [showPassword, setShowPassword] = useState(false)
  const [showConfirmPassword, setShowConfirmPassword] = useState(false)

  // Redirect if already authenticated
  if (isAuthenticated) {
    navigate('/dashboard')
    return null
  }

  const handleInputChange = (field: string, value: string) => {
    setFormData(prev => ({
      ...prev,
      [field]: value
    }))
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    // Basic validation
    if (!formData.email || !formData.password || !formData.confirmPassword) {
      toast.error("กรุณากรอกข้อมูลให้ครบถ้วน")
      return
    }

    // Email validation
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    if (!emailRegex.test(formData.email)) {
      toast.error("รูปแบบอีเมลไม่ถูกต้อง")
      return
    }

    if (formData.password !== formData.confirmPassword) {
      toast.error("รหัสผ่านไม่ตรงกัน")
      return
    }

    if (formData.password.length < 6) {
      toast.error("รหัสผ่านต้องมีอย่างน้อย 6 ตัวอักษร")
      return
    }

    const registerData: SimpleRegisterRequest = {
      email: formData.email,
      password: formData.password,
      role: "agent"
    }

    try {
      await register(registerData as RegisterRequest) // Cast to match the expected type temporarily
      toast.success("ลงทะเบียนสำเร็จ กรุณาเข้าสู่ระบบและสร้างโปรไฟล์")
      navigate('/login')
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "ลงทะเบียนไม่สำเร็จ")
    }
  }

  const handleLogin = () => {
    navigate('/login')
  }


  return (
    <div className={cn("flex flex-col gap-6", className)} {...props}>
      <Card>
        <CardHeader>
          <CardTitle>สร้างบัญชีเอเจนต์ใหม่</CardTitle>
          <CardDescription>
            กรุณากรอกข้อมูลเพื่อสร้างบัญชีเอเจนต์ (สร้างโปรไฟล์ภายหลัง)
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit}>
            <div className="flex flex-col gap-6">
              <div className="grid gap-3">
                <Label htmlFor="email">อีเมล</Label>
                <Input
                  id="email"
                  type="email"
                  placeholder="m@example.com"
                  value={formData.email}
                  onChange={(e) => handleInputChange('email', e.target.value)}
                  disabled={isRegisterLoading}
                  required
                />
              </div>

              <div className="grid gap-3">
                <Label htmlFor="password">รหัสผ่าน</Label>
                <div className="relative">
                  <Input
                    id="password"
                    type={showPassword ? "text" : "password"}
                    placeholder="รหัสผ่านอย่างน้อย 6 ตัวอักษร"
                    value={formData.password}
                    onChange={(e) => handleInputChange('password', e.target.value)}
                    disabled={isRegisterLoading}
                    required
                    className="pr-10"
                  />
                  <button
                    type="button"
                    onClick={() => setShowPassword(!showPassword)}
                    className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
                    tabIndex={-1}
                  >
                    {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                  </button>
                </div>
              </div>

              <div className="grid gap-3">
                <Label htmlFor="confirmPassword">ยืนยันรหัสผ่าน</Label>
                <div className="relative">
                  <Input
                    id="confirmPassword"
                    type={showConfirmPassword ? "text" : "password"}
                    placeholder="ยืนยันรหัสผ่าน"
                    value={formData.confirmPassword}
                    onChange={(e) => handleInputChange('confirmPassword', e.target.value)}
                    disabled={isRegisterLoading}
                    required
                    className="pr-10"
                  />
                  <button
                    type="button"
                    onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                    className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
                    tabIndex={-1}
                  >
                    {showConfirmPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                  </button>
                </div>
              </div>


              <Button type="submit" className="w-full" disabled={isRegisterLoading}>
                {isRegisterLoading ? "กำลังลงทะเบียน..." : "ลงทะเบียน"}
              </Button>

              <div className="text-center text-sm">
                มีบัญชีแล้ว?{" "}
                <button
                  type="button"
                  onClick={handleLogin}
                  className="underline underline-offset-4 hover:no-underline"
                >
                  เข้าสู่ระบบ
                </button>
              </div>

              <div className="text-center text-xs text-muted-foreground">
               หลังจากลงทะเบียนแล้ว กรุณาเข้าสู่ระบบและสร้างโปรไฟล์ของคุณ
              </div>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}