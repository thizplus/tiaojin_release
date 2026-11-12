import { useState } from "react";
import type { Banner, BannerPosition } from "@/types";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { BannerSearchAndTable } from "@/components/manage/banner/banner-search-and-table";
import { SortableBannerList } from "@/components/manage/banner/SortableBannerList";
import { CreateBannerDialog } from "@/components/manage/banner/CreateBannerDialog";
import { EditBannerDialog } from "@/components/manage/banner/EditBannerDialog";
import { DeleteBannerDialog } from "@/components/manage/banner/DeleteBannerDialog";
import { useBannersData } from "@/hooks/useBannersData";

function BannerTabContent({ position }: { position: BannerPosition }) {
  const [selectedRows, setSelectedRows] = useState<string[]>([]);
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [selectedBanner, setSelectedBanner] = useState<Banner | null>(null);
  const [viewMode, setViewMode] = useState<'table' | 'reorder'>('table');

  const { banners, refetchBannersList } = useBannersData({ position, initialPageSize: 100 });

  return (
    <>
      {/* Toggle table / reorder */}
      <div className="flex gap-2 mb-4">
        <button
          onClick={() => setViewMode('table')}
          className={`px-3 py-1.5 text-sm rounded-md ${viewMode === 'table' ? 'bg-primary text-primary-foreground' : 'bg-muted text-muted-foreground'}`}
        >
          ตาราง
        </button>
        <button
          onClick={() => setViewMode('reorder')}
          className={`px-3 py-1.5 text-sm rounded-md ${viewMode === 'reorder' ? 'bg-primary text-primary-foreground' : 'bg-muted text-muted-foreground'}`}
        >
          จัดลำดับ
        </button>
      </div>

      {viewMode === 'table' ? (
        <BannerSearchAndTable
          position={position}
          onEdit={(banner) => { setSelectedBanner(banner); setIsEditDialogOpen(true); }}
          onDelete={(banner) => { setSelectedBanner(banner); setIsDeleteDialogOpen(true); }}
          onBulkDelete={(ids) => setSelectedRows(ids)}
          onAdd={() => setIsCreateDialogOpen(true)}
          selectedRows={selectedRows}
          onRowSelectionChange={setSelectedRows}
        />
      ) : (
        <div>
          <p className="text-sm text-muted-foreground mb-4">ลากแบนเนอร์เพื่อจัดลำดับ แล้วกด "บันทึกลำดับ"</p>
          <SortableBannerList banners={banners} onOrderSaved={refetchBannersList} />
        </div>
      )}

      <CreateBannerDialog
        open={isCreateDialogOpen}
        onOpenChange={setIsCreateDialogOpen}
        onSuccess={() => refetchBannersList()}
      />
      <EditBannerDialog
        open={isEditDialogOpen}
        onOpenChange={setIsEditDialogOpen}
        banner={selectedBanner}
        onSuccess={() => { setSelectedBanner(null); refetchBannersList(); }}
      />
      <DeleteBannerDialog
        open={isDeleteDialogOpen}
        onOpenChange={setIsDeleteDialogOpen}
        banner={selectedBanner}
        onSuccess={() => { setSelectedBanner(null); refetchBannersList(); }}
      />
    </>
  );
}

export default function BannersPage() {
  const [activeTab, setActiveTab] = useState<BannerPosition>('internal');

  return (
    <div className="p-6">
      <div className="mb-6">
        <h1 className="text-2xl font-bold">จัดการแบนเนอร์</h1>
        <p className="text-muted-foreground mt-1">จัดการข้อมูลแบนเนอร์สำหรับการแสดงผลบนเว็บไซต์</p>
      </div>

      <Tabs value={activeTab} onValueChange={(v) => setActiveTab(v as BannerPosition)}>
        <TabsList>
          <TabsTrigger value="internal">แบนเนอร์ข้างใน (Dashboard)</TabsTrigger>
          <TabsTrigger value="public">แบนเนอร์หน้าแรก (Public)</TabsTrigger>
        </TabsList>

        <TabsContent value="internal">
          <BannerTabContent position="internal" />
        </TabsContent>

        <TabsContent value="public">
          <BannerTabContent position="public" />
        </TabsContent>
      </Tabs>
    </div>
  );
}
