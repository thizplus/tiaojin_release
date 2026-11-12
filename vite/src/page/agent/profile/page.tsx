import { useState } from "react";
import {  RefreshCw, AlertCircle, Building2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { Alert, AlertDescription } from "@/components/ui/alert";

import { ProfileInfoCard } from "@/components/agent/profile/ProfileInfoCard";
import { CreateProfileForm } from "@/components/agent/profile/CreateProfileForm";
import { EditProfileDialog } from "@/components/agent/profile/EditProfileDialog";
import { useAgentProfile } from "@/hooks/agent/useAgentProfile";
import type { CreateAgentProfileRequest, UpdateAgentProfileRequest } from "@/types";

const MESSAGES = {
  FETCH_FAILED: "ไม่สามารถโหลดข้อมูลได้ กรุณาลองใหม่อีกครั้ง"
};

export default function AgentProfilePage() {
  const [editDialogOpen, setEditDialogOpen] = useState(false);

  const {
    profile,
    profileStatus,
    isInitialized,
    hasProfile,
    //companyName,
    isCheckingProfile,
    isLoadingProfile,
    isCreatingProfile,
    isUpdatingProfile,
    //canCreateProfile,
    createProfile,
    updateProfile,
    refetchProfileCheck,
    refetchProfile,
    checkError,
    profileError,
  } = useAgentProfile({
    enableCheck: true,
    enableProfile: false,
    autoFetchProfile: true,
    showToast: true
  });

  const handleCreateProfile = async (data: CreateAgentProfileRequest) => {
    await createProfile(data);
  };

  const handleUpdateProfile = async (data: UpdateAgentProfileRequest) => {
    await updateProfile(data);
    setEditDialogOpen(false);
  };


  const handleRetry = () => {
    refetchProfileCheck();
    if (hasProfile) {
      refetchProfile();
    }
  };

  // Loading state while checking profile
  if (isCheckingProfile && !isInitialized) {
    return (
      <div className="min-h-screen bg-background">
        <div className="container mx-auto p-6 max-w-4xl">
          <div className="space-y-6">
            {/* Loading Header */}
            <div className="space-y-3">
              <Skeleton className="h-8 w-64" />
              <Skeleton className="h-4 w-96" />
            </div>

            {/* Loading Cards */}
            <div className="space-y-6">
              <Skeleton className="h-32 w-full rounded-lg" />
              <Skeleton className="h-64 w-full rounded-lg" />
            </div>
          </div>
        </div>
      </div>
    );
  }

  // Error state - Unable to check profile status
  if (checkError && !profileStatus) {
    return (
      <div className="min-h-screen bg-background">
        <div className="container mx-auto p-6 max-w-4xl">
          <div className="space-y-6">
            {/* Header */}
            <div className="text-center">
              <Building2 className="h-12 w-12 mx-auto mb-4 text-muted-foreground" />
              <h1 className="text-3xl font-bold text-foreground mb-2">
                โปรไฟล์บริษัท
              </h1>
              <p className="text-lg text-muted-foreground">
                จัดการข้อมูลบริษัททัวร์ของคุณ
              </p>
            </div>

            {/* Error Alert */}
            <Alert variant="destructive" className="max-w-2xl mx-auto">
              <AlertCircle className="h-4 w-4" />
              <AlertDescription>
                {checkError.message || MESSAGES.FETCH_FAILED}
              </AlertDescription>
            </Alert>

            {/* Retry Button */}
            <div className="flex justify-center">
              <Button onClick={handleRetry} variant="outline" size="lg">
                <RefreshCw className="h-4 w-4 mr-2" />
                ลองใหม่
              </Button>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background">
      <div className="container mx-auto max-w-2xl">
        <div className="space-y-12">
        


          {/* Main Content - Clean Layout */}
          <div className="space-y-12">
            {!hasProfile ? (
              // Create Profile Form
              <div className="flex justify-center">
                <CreateProfileForm
                  onSubmit={handleCreateProfile}
                  isLoading={isCreatingProfile}
                  className="w-full max-w-4xl"
                />
              </div>
            ) : profile ? (
              // Profile Info - Clean Design
              <div className="flex justify-center">
                <div className="w-full max-w-5xl">
                  <ProfileInfoCard
                    profile={profile}
                    onEdit={() => setEditDialogOpen(true)}
                    className=""
                  />
                </div>

                {/* Edit Profile Dialog */}
                <EditProfileDialog
                  open={editDialogOpen}
                  onOpenChange={setEditDialogOpen}
                  profile={profile}
                  onSubmit={handleUpdateProfile}
                  isLoading={isUpdatingProfile}
                />
              </div>
            ) : isLoadingProfile ? (
              // Loading Profile - Clean Skeletons
              <div className="flex justify-center">
                <div className="w-full max-w-5xl space-y-8">
                  <div className="p-8 bg-card border border-border rounded-xl shadow-lg">
                    <div className="animate-pulse space-y-6">
                      <div className="flex items-center space-x-6">
                        <Skeleton className="h-20 w-20 rounded-xl" />
                        <div className="space-y-3 flex-1">
                          <Skeleton className="h-6 w-64" />
                          <Skeleton className="h-4 w-40" />
                          <Skeleton className="h-4 w-32" />
                        </div>
                      </div>
                      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                        <Skeleton className="h-24 rounded-lg" />
                        <Skeleton className="h-24 rounded-lg" />
                      </div>
                    </div>
                  </div>
                  <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
                    <Skeleton className="h-80 rounded-xl" />
                    <Skeleton className="h-80 rounded-xl" />
                  </div>
                </div>
              </div>
            ) : profileError ? (
              // Profile Error - Clean Error Card
              <div className="flex justify-center">
                <div className="w-full max-w-2xl">
                  <div className="p-10 bg-card border-2 border-destructive rounded-xl text-center">
                    <div className="mx-auto w-20 h-20 rounded-full bg-muted border border-border flex items-center justify-center mb-6">
                      <AlertCircle className="h-10 w-10 text-destructive" />
                    </div>
                    <h3 className="text-2xl font-bold text-foreground mb-3">
                      เกิดข้อผิดพลาด
                    </h3>
                    <p className="text-lg text-muted-foreground mb-8">
                      {profileError.message || MESSAGES.FETCH_FAILED}
                    </p>
                    <Button
                      onClick={handleRetry}
                      variant="outline"
                      size="lg"
                      className="px-8 py-3 text-base font-medium"
                    >
                      <RefreshCw className="h-5 w-5 mr-3" />
                      ลองใหม่อีกครั้ง
                    </Button>
                  </div>
                </div>
              </div>
            ) : (
              // Fallback - No profile data - Clean Design
              <div className="flex justify-center">
                <div className="w-full max-w-2xl">
                  <div className="p-12 bg-card border border-border rounded-xl text-center">
                    <div className="mx-auto w-24 h-24 rounded-xl bg-muted border border-border flex items-center justify-center mb-8">
                      <AlertCircle className="h-12 w-12 text-muted-foreground" />
                    </div>
                    <h3 className="text-2xl font-bold text-foreground mb-4">
                      ไม่พบข้อมูลโปรไฟล์
                    </h3>
                    <p className="text-lg text-muted-foreground mb-10">
                      ไม่สามารถโหลดข้อมูลโปรไฟล์บริษัทได้ กรุณาลองใหม่อีกครั้ง
                    </p>
                    <Button
                      onClick={handleRetry}
                      variant="outline"
                      size="lg"
                      className="px-10 py-4 text-base font-medium"
                    >
                      <RefreshCw className="h-5 w-5 mr-3" />
                      ลองใหม่
                    </Button>
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}