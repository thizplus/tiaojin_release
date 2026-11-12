import { useState, useEffect, useRef } from 'react';
import { Search, User, X } from 'lucide-react';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Loader2 } from 'lucide-react';
import salesService from '@/services/admin/salesService';
import type { Sale } from '@/types/models/sales';

interface SalesSelectorProps {
  selectedSales: Sale | null;
  onSalesSelect: (sales: Sale | null) => void;
}

export function SalesSelector({ selectedSales, onSalesSelect }: SalesSelectorProps) {
  const [sales, setSales] = useState<Sale[]>([]);
  const [searchTerm, setSearchTerm] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  // โหลด active sales
  useEffect(() => {
    const loadActiveSales = async () => {
      setIsLoading(true);
      try {
        const response = await salesService.getSalesActive();
        if (response.success && response.data.sales) {
          setSales(response.data.sales);
        }
      } catch (error) {
        console.error('Failed to load sales:', error);
        setSales([]);
      } finally {
        setIsLoading(false);
      }
    };

    loadActiveSales();
  }, []);

  // Click outside handler
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    }

    if (isOpen) {
      document.addEventListener('mousedown', handleClickOutside);
      return () => {
        document.removeEventListener('mousedown', handleClickOutside);
      };
    }
  }, [isOpen]);

  // Filter sales based on search term
  const filteredSales = sales.filter((sale) => {
    const searchLower = searchTerm.toLowerCase();
    return (
      sale.code.toLowerCase().includes(searchLower) ||
      sale.first_name.toLowerCase().includes(searchLower) ||
      sale.last_name.toLowerCase().includes(searchLower) ||
      sale.nick_name?.toLowerCase().includes(searchLower)
    );
  });

  const handleSalesSelect = (sale: Sale) => {
    onSalesSelect(sale);
    setIsOpen(false);
    setSearchTerm('');
  };

  const handleClear = () => {
    onSalesSelect(null);
    setSearchTerm('');
  };

  const displayName = (sale: Sale) => {
    return `${sale.code} - ${sale.first_name} ${sale.last_name}${sale.nick_name ? ` (${sale.nick_name})` : ''}`;
  };

  return (
    <div className="space-y-2">
      <label className="text-sm font-medium">เลือก Sales</label>
      <div className="relative" ref={dropdownRef}>
        <div className="flex gap-2">
          <div className="flex-1 relative">
            <Search className="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="ค้นหา Sales ด้วยรหัสหรือชื่อ..."
              value={selectedSales ? displayName(selectedSales) : searchTerm}
              onChange={(e) => {
                if (!selectedSales) {
                  setSearchTerm(e.target.value);
                  setIsOpen(true);
                }
              }}
              onFocus={() => setIsOpen(true)}
              className="pl-10"
              readOnly={!!selectedSales}
            />
          </div>
          {selectedSales && (
            <Button variant="outline" onClick={handleClear}>
              <X className="h-4 w-4" />
            </Button>
          )}
        </div>

        {/* Dropdown รายการ sales */}
        {isOpen && !selectedSales && (
          <div className="absolute top-full left-0 right-0 z-10 mt-1 bg-background border rounded-md shadow-lg">
            <ScrollArea className="max-h-64">
              {isLoading ? (
                <div className="flex items-center justify-center p-4">
                  <Loader2 className="h-4 w-4 animate-spin" />
                  <span className="ml-2 text-sm text-muted-foreground">กำลังโหลด...</span>
                </div>
              ) : filteredSales.length === 0 ? (
                <div className="p-4 text-center text-sm text-muted-foreground">
                  {searchTerm ? 'ไม่พบ Sales ที่ค้นหา' : 'ไม่มี Sales ที่ Active'}
                </div>
              ) : (
                <div className="p-2">
                  {filteredSales.map((sale) => (
                    <div
                      key={sale.id}
                      className="flex items-center gap-3 p-3 hover:bg-muted cursor-pointer rounded"
                      onClick={() => handleSalesSelect(sale)}
                    >
                      <User className="h-4 w-4 text-muted-foreground flex-shrink-0" />
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2">
                          <Badge variant="outline" className="text-xs">
                            {sale.code}
                          </Badge>
                        </div>
                        <div className="text-sm font-medium truncate">
                          {sale.first_name} {sale.last_name}
                          {sale.nick_name && (
                            <span className="text-muted-foreground ml-1">({sale.nick_name})</span>
                          )}
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </ScrollArea>
          </div>
        )}
      </div>
    </div>
  );
}