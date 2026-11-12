import { useState, useCallback, useEffect } from 'react';
import {
  DndContext,
  closestCenter,
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors,
  type DragEndEvent,
} from '@dnd-kit/core';
import {
  arrayMove,
  SortableContext,
  sortableKeyboardCoordinates,
  useSortable,
  verticalListSortingStrategy,
} from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import { GripVertical, Loader2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { getStorageUrl } from '@/lib/utils';
import { useToast } from '@/hooks/use-toast';
import bannerService from '@/services/admin/bannerService';
import type { Banner } from '@/types';

interface SortableBannerListProps {
  banners: Banner[];
  onOrderSaved?: () => void;
}

function SortableBannerItem({ banner }: { banner: Banner }) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: banner.id });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
  };

  return (
    <div
      ref={setNodeRef}
      style={style}
      className={`flex items-center gap-3 p-3 rounded-lg border bg-card ${isDragging ? 'shadow-lg z-10' : 'hover:bg-muted/50'}`}
    >
      <button {...attributes} {...listeners} className="cursor-grab active:cursor-grabbing text-muted-foreground hover:text-foreground">
        <GripVertical className="h-5 w-5" />
      </button>

      <div className="h-12 w-12 rounded overflow-hidden flex-shrink-0 bg-muted">
        {banner.image && (
          <img src={getStorageUrl(banner.image)} alt={banner.title} className="h-full w-full object-cover" />
        )}
      </div>

      <div className="flex-1 min-w-0">
        <p className="font-medium text-sm truncate">{banner.title}</p>
        <div className="flex items-center gap-2 mt-1">
          <Badge variant={banner.status === 'publish' ? 'default' : 'secondary'} className="text-xs">
            {banner.status === 'publish' ? 'เผยแพร่' : 'ร่าง'}
          </Badge>
          <Badge variant="outline" className="text-xs">
            {banner.position === 'public' ? 'หน้าแรก' : 'ข้างใน'}
          </Badge>
        </div>
      </div>

      <span className="text-sm text-muted-foreground font-mono w-8 text-center">
        {banner.display_order}
      </span>
    </div>
  );
}

export function SortableBannerList({ banners: initialBanners, onOrderSaved }: SortableBannerListProps) {
  const [items, setItems] = useState(initialBanners);
  const [hasChanges, setHasChanges] = useState(false);
  const [isSaving, setIsSaving] = useState(false);
  const { toast } = useToast();

  const sensors = useSensors(
    useSensor(PointerSensor),
    useSensor(KeyboardSensor, { coordinateGetter: sortableKeyboardCoordinates })
  );

  // Sync เมื่อ prop เปลี่ยน + sort ตาม display_order
  useEffect(() => {
    const sorted = [...initialBanners].sort((a, b) => a.display_order - b.display_order);
    setItems(sorted);
    setHasChanges(false);
  }, [initialBanners]);

  const handleDragEnd = useCallback((event: DragEndEvent) => {
    const { active, over } = event;
    if (!over || active.id === over.id) return;

    setItems((prev) => {
      const oldIndex = prev.findIndex((i) => i.id === active.id);
      const newIndex = prev.findIndex((i) => i.id === over.id);
      return arrayMove(prev, oldIndex, newIndex);
    });
    setHasChanges(true);
  }, []);

  const handleSave = async () => {
    setIsSaving(true);
    try {
      // อัปเดตทุก banner ด้วยลำดับใหม่
      const promises = items.map((banner, index) =>
        bannerService.updateBannerOrder(banner.id, { display_order: index + 1 })
      );
      await Promise.all(promises);

      toast({ title: 'บันทึกลำดับแบนเนอร์สำเร็จ' });
      setHasChanges(false);
      onOrderSaved?.();
    } catch {
      toast({ title: 'เกิดข้อผิดพลาด', description: 'ไม่สามารถบันทึกลำดับได้', variant: 'destructive' });
    } finally {
      setIsSaving(false);
    }
  };

  if (items.length === 0) {
    return <p className="text-center text-muted-foreground py-8">ไม่มีแบนเนอร์</p>;
  }

  return (
    <div className="space-y-4">
      <DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={handleDragEnd}>
        <SortableContext items={items.map((i) => i.id)} strategy={verticalListSortingStrategy}>
          <div className="space-y-2">
            {items.map((banner) => (
              <SortableBannerItem key={banner.id} banner={banner} />
            ))}
          </div>
        </SortableContext>
      </DndContext>

      {hasChanges && (
        <div className="flex justify-end">
          <Button onClick={handleSave} disabled={isSaving}>
            {isSaving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            บันทึกลำดับ
          </Button>
        </div>
      )}
    </div>
  );
}
