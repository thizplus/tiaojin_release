import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { bankAccountService } from '@/services/admin/bankAccountService';
import type { CreateBankAccountRequest, UpdateBankAccountRequest } from '@/types/models/bankAccount';

const keys = {
  all: ['bank-accounts'] as const,
  list: () => [...keys.all, 'list'] as const,
  active: () => [...keys.all, 'active'] as const,
};

export function useBankAccounts() {
  const { data, isLoading } = useQuery({
    queryKey: keys.list(),
    queryFn: () => bankAccountService.list(),
  });

  return {
    bankAccounts: data?.success ? data.data : [],
    isLoading,
  };
}

export function useActiveBankAccounts() {
  const { data, isLoading } = useQuery({
    queryKey: keys.active(),
    queryFn: () => bankAccountService.listActive(),
  });

  return {
    bankAccounts: data?.success ? data.data : [],
    isLoading,
  };
}

export function useCreateBankAccount() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateBankAccountRequest) => bankAccountService.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: keys.all });
    },
  });
}

export function useUpdateBankAccount() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateBankAccountRequest }) => bankAccountService.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: keys.all });
    },
  });
}

export function useDeleteBankAccount() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => bankAccountService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: keys.all });
    },
  });
}
