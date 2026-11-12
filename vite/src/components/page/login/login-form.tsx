import { useState } from "react"
import { useNavigate } from "react-router-dom"
import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { useAuth } from "@/hooks/useAuth"
import { toast } from "sonner"
import { Eye, EyeOff } from "lucide-react"

export function LoginForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [showPassword, setShowPassword] = useState(false)
  const navigate = useNavigate()

  const { login, isLoginLoading, isAuthenticated } = useAuth()

  // Redirect if already authenticated
  if (isAuthenticated) {
    navigate('/dashboard')
    return null
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!email || !password) {
      toast.error("กรุณากรอกชื่อผู้ใช้/อีเมล และรหัสผ่าน")
      return
    }

    try {
      const result = await login({ email, password })
      toast.success(`เข้าสู่ระบบสำเร็จ ${result.user?.email || ''}`)
      navigate('/dashboard')
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "เข้าสู่ระบบไม่สำเร็จ")
    }
  }

  const handleForgotPassword = () => {
    navigate('/forgot-password')
  }

  const handleRegister = () => {
    navigate('/register')
  }

  return (
    <div className={cn("flex flex-col gap-6", className)} {...props}>
      <Card className="overflow-hidden p-0">
        <CardContent className="grid p-0 md:grid-cols-2">
          <form onSubmit={handleSubmit} className="p-6 md:p-8">
            <div className="flex flex-col gap-6">
              <div className="flex flex-col items-center text-center">
                <h1 className="text-2xl font-bold">{import.meta.env.VITE_APP_TITLE || 'APPLICATION'}</h1>
                <p className="text-muted-foreground text-balance">
                  กรุณาเข้าสู่ระบบด้วยบัญชีผู้ใช้ของคุณ
                </p>
              </div>
              <div className="grid gap-3">
                <Label htmlFor="email">ชื่อผู้ใช้หรืออีเมล</Label>
                <Input
                  id="email"
                  type="text"
                  placeholder="username หรือ email@example.com"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  disabled={isLoginLoading}
                  required
                />
              </div>
              <div className="grid gap-3">
                <div className="flex items-center">
                  <Label htmlFor="password">รหัสผ่าน</Label>
                  <button
                    type="button"
                    onClick={handleForgotPassword}
                    className="ml-auto text-sm underline-offset-2 hover:underline"
                  >
                    ลืมรหัสผ่าน?
                  </button>
                </div>
                <div className="relative">
                  <Input
                    id="password"
                    type={showPassword ? "text" : "password"}
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    disabled={isLoginLoading}
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
              <Button type="submit" className="w-full" disabled={isLoginLoading}>
                {isLoginLoading ? "กำลังเข้าสู่ระบบ..." : "เข้าสู่ระบบ"}
              </Button>
             
            
              <div className="text-center text-sm">
                ยังไม่มีบัญชี?{" "}
                <button
                  type="button"
                  onClick={handleRegister}
                  className="underline underline-offset-4 hover:no-underline"
                >
                  ลงทะเบียน
                </button>
              </div>
            </div>
          </form>
          <div className="bg-muted relative hidden md:block">
            <img
              src="/images/background/mascot.png"
              alt="Image"
              className="absolute inset-0 h-full w-full object-cover"
            />
          </div>
        </CardContent>
      </Card>
      <div className="text-muted-foreground *:[a]:hover:text-primary text-center text-xs text-balance *:[a]:underline *:[a]:underline-offset-4">
      กรุณาอ่านและยอมรับ <a href="#">ข้อตกลงการให้บริการ</a>{" "}
        และ <a href="#">นโยบายความเป็นส่วนตัว</a>.
      </div>
    </div>
  )
}
