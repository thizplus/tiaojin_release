import { useState } from 'react'
import { useAuth } from '@/hooks/useAuth'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Label } from '@/components/ui/label'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { Loader2, Eye, EyeOff, CheckCircle2, Clock, XCircle } from 'lucide-react'
import { toast } from 'sonner'
import { cn } from '@/lib/utils'

export default function ProfilePage() {
  const { user, updatePassword, isUpdatePasswordLoading } = useAuth()
  const [oldPassword, setOldPassword] = useState('')
  const [newPassword, setNewPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [showOldPassword, setShowOldPassword] = useState(false)
  const [showNewPassword, setShowNewPassword] = useState(false)
  const [showConfirmPassword, setShowConfirmPassword] = useState(false)

  // Get role display text
  const getRoleDisplay = () => {
    switch (user?.role) {
      case 'super':
        return 'ผู้ดูแลระบบสูงสุด'
      case 'admin':
        return 'ผู้ดูแลระบบ'
      case 'agent':
        return 'เอเจนต์'
      case 'sales':
        return 'เซลส์'
      default:
        return user?.role || ''
    }
  }

  // Get status display
  const getStatusDisplay = () => {
    switch (user?.status) {
      case 'active':
        return {
          label: 'อนุมัติแล้ว',
          icon: CheckCircle2,
          color: 'text-green-500',
          bgColor: 'bg-green-500'
        }
      case 'pending':
        return {
          label: 'รอการอนุมัติ',
          icon: Clock,
          color: 'text-yellow-500',
          bgColor: 'bg-yellow-500'
        }
      case 'inactive':
        return {
          label: 'ระงับการใช้งาน',
          icon: XCircle,
          color: 'text-red-500',
          bgColor: 'bg-red-500'
        }
      default:
        return {
          label: user?.status || '',
          icon: CheckCircle2,
          color: 'text-muted-foreground',
          bgColor: 'bg-muted-foreground'
        }
    }
  }

  // Generate user initials for avatar fallback
  const getInitials = (email?: string) => {
    if (!email) return 'NA'
    return email.substring(0, 2).toUpperCase()
  }

  // Get display name
  const getDisplayName = () => {
    if (user?.first_name && user?.last_name) {
      return `${user.first_name} ${user.last_name}`
    }
    return user?.email || ''
  }

  const handleUpdatePassword = async (e: React.FormEvent) => {
    e.preventDefault()

    // Validation
    if (!oldPassword || !newPassword || !confirmPassword) {
      toast.error('กรุณากรอกข้อมูลให้ครบถ้วน')
      return
    }

    if (newPassword !== confirmPassword) {
      toast.error('รหัสผ่านใหม่ไม่ตรงกัน')
      return
    }

    if (newPassword.length < 6) {
      toast.error('รหัสผ่านใหม่ต้องมีความยาวอย่างน้อย 6 ตัวอักษร')
      return
    }

    try {
      await updatePassword({
        id: user?.id || '',
        old_password: oldPassword,
        new_password: newPassword
      })

      toast.success('เปลี่ยนรหัสผ่านสำเร็จ')

      // Clear form
      setOldPassword('')
      setNewPassword('')
      setConfirmPassword('')
    } catch (error) {
      toast.error('เกิดข้อผิดพลาดในการเปลี่ยนรหัสผ่าน', {
        description: error instanceof Error ? error.message : 'กรุณาลองใหม่อีกครั้ง'
      })
    }
  }

  if (!user) {
    return (
      <div className="flex items-center justify-center py-12">
        <Loader2 className="h-8 w-8 animate-spin" />
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold">โปรไฟล์</h1>

      {/* User Information */}
      <Card>
        <CardHeader>
          <CardTitle>ข้อมูลบัญชี</CardTitle>
          <CardDescription>ข้อมูลผู้ใช้งานในระบบ</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex items-start gap-6">
            {/* Avatar */}
            <Avatar className="h-24 w-24 rounded-lg">
              <AvatarImage src="" alt={getDisplayName()} />
              <AvatarFallback className="rounded-lg text-2xl">
                {getInitials(user.email)}
              </AvatarFallback>
            </Avatar>

            {/* User Details */}
            <div className="flex-1 space-y-4">
              {/* Name/Email */}
              <div className="space-y-1">
                
                <p className="text-lg font-semibold">{getDisplayName()}</p>
                <Badge variant="secondary" className="text-sm">
                  {getRoleDisplay()}
                </Badge>
                <div className="flex items-center gap-2">
                  <span className={cn("h-2 w-2 rounded-full", getStatusDisplay().bgColor)} />
                  <span className={cn("text-sm font-medium", getStatusDisplay().color)}>
                    {getStatusDisplay().label}
                  </span>
                </div>
              </div>

             

            
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Change Password */}
      <Card>
        <CardHeader>
          <CardTitle>เปลี่ยนรหัสผ่าน</CardTitle>
          <CardDescription>อัปเดตรหัสผ่านของคุณ</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleUpdatePassword} className="space-y-4">
            <div className="grid gap-2">
              <Label htmlFor="old-password">รหัสผ่านเดิม</Label>
              <div className="relative">
                <Input
                  id="old-password"
                  type={showOldPassword ? 'text' : 'password'}
                  value={oldPassword}
                  onChange={(e) => setOldPassword(e.target.value)}
                  placeholder="กรอกรหัสผ่านเดิม"
                  disabled={isUpdatePasswordLoading}
                />
                <Button
                  type="button"
                  variant="ghost"
                  size="icon"
                  className="absolute right-2 top-1/2 -translate-y-1/2"
                  onClick={() => setShowOldPassword(!showOldPassword)}
                >
                  {showOldPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                </Button>
              </div>
            </div>

            <div className="grid gap-2">
              <Label htmlFor="new-password">รหัสผ่านใหม่</Label>
              <div className="relative">
                <Input
                  id="new-password"
                  type={showNewPassword ? 'text' : 'password'}
                  value={newPassword}
                  onChange={(e) => setNewPassword(e.target.value)}
                  placeholder="กรอกรหัสผ่านใหม่ (อย่างน้อย 6 ตัวอักษร)"
                  disabled={isUpdatePasswordLoading}
                />
                <Button
                  type="button"
                  variant="ghost"
                  size="icon"
                  className="absolute right-2 top-1/2 -translate-y-1/2"
                  onClick={() => setShowNewPassword(!showNewPassword)}
                >
                  {showNewPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                </Button>
              </div>
            </div>

            <div className="grid gap-2">
              <Label htmlFor="confirm-password">ยืนยันรหัสผ่านใหม่</Label>
              <div className="relative">
                <Input
                  id="confirm-password"
                  type={showConfirmPassword ? 'text' : 'password'}
                  value={confirmPassword}
                  onChange={(e) => setConfirmPassword(e.target.value)}
                  placeholder="กรอกรหัสผ่านใหม่อีกครั้ง"
                  disabled={isUpdatePasswordLoading}
                />
                <Button
                  type="button"
                  variant="ghost"
                  size="icon"
                  className="absolute right-2 top-1/2 -translate-y-1/2"
                  onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                >
                  {showConfirmPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                </Button>
              </div>
            </div>

            <Button type="submit" disabled={isUpdatePasswordLoading}>
              {isUpdatePasswordLoading ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  กำลังเปลี่ยนรหัสผ่าน...
                </>
              ) : (
                'เปลี่ยนรหัสผ่าน'
              )}
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}