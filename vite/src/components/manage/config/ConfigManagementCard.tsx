import { useState } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import {
  Database,
  Plus,
  Search,
  Filter,
  MoreHorizontal,
  Edit3,
  Trash2,
  Eye,
  Loader2,
  RefreshCw
} from 'lucide-react';
import { useConfigs, useDeleteConfig } from '@/stores/configStore';
import { ConfigDialog } from './ConfigDialog';
import { useToast } from '@/hooks/use-toast';
import { format } from 'date-fns';
import { th } from 'date-fns/locale';
import type { ConfigType, SystemConfig } from '@/types/config';

const CONFIG_TYPES: { value: ConfigType; label: string }[] = [
  { value: 'expiration', label: 'การหมดอายุ' },
  { value: 'booking_limits', label: 'การจำกัดการจอง' },
  { value: 'payment', label: 'การชำระเงิน' },
  { value: 'notification', label: 'การแจ้งเตือน' },
  { value: 'security', label: 'ความปลอดภัย' },
  { value: 'integration', label: 'การเชื่อมต่อ' },
  { value: 'general', label: 'ทั่วไป' },
];

export function ConfigManagementCard() {
  const { toast } = useToast();
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedType, setSelectedType] = useState<string>('all');
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [dialogMode, setDialogMode] = useState<'create' | 'edit' | 'view'>('create');
  const [selectedConfig, setSelectedConfig] = useState<SystemConfig | null>(null);

  const { data, isLoading, error, refetch } = useConfigs({
    config_type: selectedType === 'all' ? undefined : selectedType,
    limit: 50,
  });

  const deleteMutation = useDeleteConfig();

  const filteredConfigs = data?.configs?.filter(config =>
    config.config_key.toLowerCase().includes(searchTerm.toLowerCase()) ||
    config.description.toLowerCase().includes(searchTerm.toLowerCase())
  ) || [];

  const handleView = (config: SystemConfig) => {
    setSelectedConfig(config);
    setDialogMode('view');
    setIsDialogOpen(true);
  };

  const handleEdit = (config: SystemConfig) => {
    setSelectedConfig(config);
    setDialogMode('edit');
    setIsDialogOpen(true);
  };

  const handleDelete = async (config: SystemConfig) => {
    if (!confirm(`คุณต้องการลบการตั้งค่า "${config.config_key}" หรือไม่?`)) {
      return;
    }

    try {
      await deleteMutation.mutateAsync(config.id);
      toast({
        title: "ลบสำเร็จ",
        description: `ลบการตั้งค่า ${config.config_key} แล้ว`,
      });
    } catch (error) {
      toast({
        title: "เกิดข้อผิดพลาด",
        description: error instanceof Error ? error.message : "ไม่สามารถลบได้",
        variant: "destructive",
      });
    }
  };

  const handleCreate = () => {
    setSelectedConfig(null);
    setDialogMode('create');
    setIsDialogOpen(true);
  };

  const getTypeColor = (type: string) => {
    switch (type) {
      case 'expiration': return 'destructive';
      case 'booking_limits': return 'default';
      case 'payment': return 'default';
      case 'notification': return 'secondary';
      case 'security': return 'outline';
      case 'integration': return 'destructive';
      default: return 'secondary';
    }
  };

  if (error) {
    return (
      <Card>
        <CardContent className="p-6">
          <div className="text-center text-destructive">
            เกิดข้อผิดพลาด: {error.message}
            <Button variant="outline" onClick={() => refetch()} className="ml-2">
              <RefreshCw className="h-4 w-4" />
            </Button>
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <>
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-4">
          <div className="flex items-center space-x-2">
            <Database className="h-5 w-5" />
            <div>
              <CardTitle className="text-lg">การจัดการการตั้งค่าระบบ</CardTitle>
              <CardDescription>
                จัดการการตั้งค่าต่างๆ ของระบบ
              </CardDescription>
            </div>
          </div>
          <Button onClick={handleCreate}>
            <Plus className="h-4 w-4 mr-2" />
            เพิ่มการตั้งค่า
          </Button>
        </CardHeader>

        <CardContent className="space-y-4">
          {/* Filters */}
          <div className="flex flex-col sm:flex-row gap-4">
            <div className="flex-1">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                <Input
                  placeholder="ค้นหาด้วยชื่อหรือคำอธิบาย..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="pl-10"
                />
              </div>
            </div>
            <div className="w-full sm:w-48">
              <Select value={selectedType} onValueChange={setSelectedType}>
                <SelectTrigger>
                  <Filter className="h-4 w-4 mr-2" />
                  <SelectValue placeholder="ประเภท" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">ทั้งหมด</SelectItem>
                  {CONFIG_TYPES.map((type) => (
                    <SelectItem key={type.value} value={type.value}>
                      {type.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>

          {/* Table */}
          {isLoading ? (
            <div className="flex items-center justify-center py-8">
              <Loader2 className="h-6 w-6 animate-spin" />
              <span className="ml-2">กำลังโหลดข้อมูล...</span>
            </div>
          ) : (
            <div className="border rounded-lg">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Key</TableHead>
                    <TableHead>ประเภท</TableHead>
                    <TableHead>คำอธิบาย</TableHead>
                    <TableHead>สถานะ</TableHead>
                    <TableHead>อัพเดทล่าสุด</TableHead>
                    <TableHead className="text-right">การดำเนินการ</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {filteredConfigs.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={6} className="text-center py-8 text-muted-foreground">
                        ไม่พบข้อมูลการตั้งค่า
                      </TableCell>
                    </TableRow>
                  ) : (
                    filteredConfigs.map((config) => (
                      <TableRow key={config.id}>
                        <TableCell className="font-mono text-sm">
                          {config.config_key}
                        </TableCell>
                        <TableCell>
                          <Badge variant={getTypeColor(config.config_type)}>
                            {CONFIG_TYPES.find(t => t.value === config.config_type)?.label || config.config_type}
                          </Badge>
                        </TableCell>
                        <TableCell className="max-w-xs truncate">
                          {config.description}
                        </TableCell>
                        <TableCell>
                          <Badge variant={config.is_active ? 'default' : 'secondary'}>
                            {config.is_active ? 'ใช้งาน' : 'ปิดใช้งาน'}
                          </Badge>
                        </TableCell>
                        <TableCell className="text-sm text-muted-foreground">
                          {format(new Date(config.updated_at), 'dd MMM yyyy HH:mm', { locale: th })}
                        </TableCell>
                        <TableCell className="text-right">
                          <DropdownMenu>
                            <DropdownMenuTrigger asChild>
                              <Button variant="ghost" className="h-8 w-8 p-0">
                                <MoreHorizontal className="h-4 w-4" />
                              </Button>
                            </DropdownMenuTrigger>
                            <DropdownMenuContent align="end">
                              <DropdownMenuLabel>การดำเนินการ</DropdownMenuLabel>
                              <DropdownMenuSeparator />
                              <DropdownMenuItem onClick={() => handleView(config)}>
                                <Eye className="h-4 w-4 mr-2" />
                                ดู
                              </DropdownMenuItem>
                              <DropdownMenuItem onClick={() => handleEdit(config)}>
                                <Edit3 className="h-4 w-4 mr-2" />
                                แก้ไข
                              </DropdownMenuItem>
                              <DropdownMenuSeparator />
                              <DropdownMenuItem
                                onClick={() => handleDelete(config)}
                                className="text-destructive"
                                disabled={deleteMutation.isPending}
                              >
                                <Trash2 className="h-4 w-4 mr-2" />
                                ลบ
                              </DropdownMenuItem>
                            </DropdownMenuContent>
                          </DropdownMenu>
                        </TableCell>
                      </TableRow>
                    ))
                  )}
                </TableBody>
              </Table>
            </div>
          )}

          {/* Pagination Info */}
          {data?.pagination && (
            <div className="flex items-center justify-between text-sm text-muted-foreground">
              <div>
                แสดง {filteredConfigs.length} รายการ จากทั้งหมด {data.pagination.total} รายการ
              </div>
              <div>
                หน้า {data.pagination.current_page} จาก {data.pagination.total_pages}
              </div>
            </div>
          )}
        </CardContent>
      </Card>

      <ConfigDialog
        open={isDialogOpen}
        onOpenChange={setIsDialogOpen}
        mode={dialogMode}
        config={selectedConfig}
      />
    </>
  );
}