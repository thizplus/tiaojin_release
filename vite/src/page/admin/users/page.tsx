import { useCallback, useState } from "react";
import type { User } from "@/types";
import { UserSearchAndTable } from "@/components/manage/user/user-search-and-table";
import { CreateUserDialog } from "@/components/manage/user/CreateUserDialog";
import { DeleteUserDialog } from "@/components/manage/user/DeleteUserDialog";
import { EditUserDialog } from "@/components/manage/user/EditUserDialog";
import { ChangeUserStatusDialog } from "@/components/manage/user/ChangeUserStatusDialog";
import { ChangeUserRoleDialog } from "@/components/manage/user/ChangeUserRoleDialog";
import { ResetPasswordDialog } from "@/components/manage/user/ResetPasswordDialog";

export default function UsersPage() {
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [isChangeStatusDialogOpen, setIsChangeStatusDialogOpen] = useState(false);
  const [isChangeRoleDialogOpen, setIsChangeRoleDialogOpen] = useState(false);
  const [isResetPasswordDialogOpen, setIsResetPasswordDialogOpen] = useState(false);
  const [selectedUser, setSelectedUser] = useState<User | null>(null);

  const handleEdit = useCallback((user: User) => {
    setSelectedUser(user);
    setIsEditDialogOpen(true);
  }, []);

  const handleDelete = useCallback((user: User) => {
    setSelectedUser(user);
    setIsDeleteDialogOpen(true);
  }, []);

  const handleChangeStatus = useCallback((user: User) => {
    setSelectedUser(user);
    setIsChangeStatusDialogOpen(true);
  }, []);

  const handleChangeRole = useCallback((user: User) => {
    setSelectedUser(user);
    setIsChangeRoleDialogOpen(true);
  }, []);

  const handleResetPassword = useCallback((user: User) => {
    setSelectedUser(user);
    setIsResetPasswordDialogOpen(true);
  }, []);

  const handleAdd = useCallback(() => {
    setIsCreateDialogOpen(true);
  }, []);

  const handleCreateSuccess = useCallback(() => {
    console.log("User created successfully");
  }, []);

  const handleDeleteSuccess = useCallback(() => {
    console.log("User deleted successfully");
    setSelectedUser(null);
  }, []);

  const handleEditSuccess = useCallback(() => {
    console.log("User updated successfully");
    setSelectedUser(null);
  }, []);

  const handleStatusChangeSuccess = useCallback(() => {
    console.log("User status changed successfully");
    setSelectedUser(null);
  }, []);

  const handleRoleChangeSuccess = useCallback(() => {
    console.log("User role changed successfully");
    setSelectedUser(null);
  }, []);

  const handleResetPasswordSuccess = useCallback(() => {
    console.log("Password reset successfully");
    setSelectedUser(null);
  }, []);

  return (
    <div className="p-6">
      <UserSearchAndTable
        title="จัดการผู้ใช้"
        description="จัดการข้อมูลผู้ใช้และสิทธิ์การเข้าถึงระบบ"
        onEdit={handleEdit}
        onDelete={handleDelete}
        onChangeStatus={handleChangeStatus}
        onChangeRole={handleChangeRole}
        onResetPassword={handleResetPassword}
        onAdd={handleAdd}
      />

      <CreateUserDialog
        open={isCreateDialogOpen}
        onOpenChange={setIsCreateDialogOpen}
        onSuccess={handleCreateSuccess}
      />

      <DeleteUserDialog
        open={isDeleteDialogOpen}
        onOpenChange={setIsDeleteDialogOpen}
        user={selectedUser}
        onSuccess={handleDeleteSuccess}
      />

      <EditUserDialog
        open={isEditDialogOpen}
        onOpenChange={setIsEditDialogOpen}
        user={selectedUser}
        onSuccess={handleEditSuccess}
      />

      <ChangeUserStatusDialog
        open={isChangeStatusDialogOpen}
        onOpenChange={setIsChangeStatusDialogOpen}
        user={selectedUser}
        onSuccess={handleStatusChangeSuccess}
      />

      <ChangeUserRoleDialog
        open={isChangeRoleDialogOpen}
        onOpenChange={setIsChangeRoleDialogOpen}
        user={selectedUser}
        onSuccess={handleRoleChangeSuccess}
      />

      <ResetPasswordDialog
        open={isResetPasswordDialogOpen}
        onOpenChange={setIsResetPasswordDialogOpen}
        user={selectedUser}
        onSuccess={handleResetPasswordSuccess}
      />
    </div>
  );
}
