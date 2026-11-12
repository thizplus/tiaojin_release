import { useState, useCallback } from 'react';
import { useToast } from '@/hooks/use-toast';
import { useCompanyProfile } from '@/hooks/useCompanyProfile';
import { CompanyProfileSearchAndTable } from '@/components/manage/profile/company-profile-search-and-table';
import { ViewCompanyProfileDialog } from '@/components/manage/profile/ViewCompanyProfileDialog';
import { EditCompanyProfileDialog } from '@/components/manage/profile/EditCompanyProfileDialog';
import { ApproveCompanyProfileDialog } from '@/components/manage/profile/ApproveCompanyProfileDialog';
import type { CompanyProfile } from '@/types';

export default function ProfilesPage() {
  const [selectedProfile, setSelectedProfile] = useState<CompanyProfile | null>(null);
  const [profileToApprove, setProfileToApprove] = useState<CompanyProfile | null>(null);
  const [profileToEdit, setProfileToEdit] = useState<CompanyProfile | null>(null);
  const [viewDialogOpen, setViewDialogOpen] = useState(false);
  const [editDialogOpen, setEditDialogOpen] = useState(false);
  const [approveDialogOpen, setApproveDialogOpen] = useState(false);

  const { approveCompanyProfile, isApproveLoading, refetchCompanyProfilesList } = useCompanyProfile();
  const { toast } = useToast();

  const handleView = useCallback((profile: CompanyProfile) => {
    setSelectedProfile(profile);
    setViewDialogOpen(true);
  }, []);

  const handleEdit = useCallback((profile: CompanyProfile) => {
    setProfileToEdit(profile);
    setEditDialogOpen(true);
  }, []);

  const handleApprove = useCallback((profile: CompanyProfile) => {
    setProfileToApprove(profile);
    setApproveDialogOpen(true);
  }, []);

  const handleConfirmApprove = useCallback(async (profile: CompanyProfile) => {
    try {
      await approveCompanyProfile({ id: profile.id });
      await refetchCompanyProfilesList();
      toast({
        title: "อนุมัติสำเร็จ",
        description: `อนุมัติโปรไฟล์บริษัท "${profile.name}" เรียบร้อยแล้ว`,
      });
    } catch (error) {
      toast({
        title: "เกิดข้อผิดพลาด",
        description: error instanceof Error ? error.message : "ไม่สามารถอนุมัติโปรไฟล์ได้",
        variant: "destructive",
      });
      throw error;
    }
  }, [approveCompanyProfile, refetchCompanyProfilesList, toast]);

  const handleEditSuccess = useCallback(() => {
    refetchCompanyProfilesList();
  }, [refetchCompanyProfilesList]);

  return (
    <div className="space-y-6">
      <CompanyProfileSearchAndTable
        title="จัดการโปรไฟล์บริษัท"
        description="ตรวจสอบ แก้ไข และอนุมัติโปรไฟล์บริษัท"
        onView={handleView}
        onEdit={handleEdit}
        onApprove={handleApprove}
      />

      <ViewCompanyProfileDialog
        profile={selectedProfile}
        open={viewDialogOpen}
        onOpenChange={setViewDialogOpen}
      />

      <EditCompanyProfileDialog
        profile={profileToEdit}
        open={editDialogOpen}
        onOpenChange={setEditDialogOpen}
        onSuccess={handleEditSuccess}
      />

      <ApproveCompanyProfileDialog
        profile={profileToApprove}
        open={approveDialogOpen}
        onOpenChange={setApproveDialogOpen}
        onConfirm={handleConfirmApprove}
        isLoading={isApproveLoading}
      />
    </div>
  );
}
