import { deleteApiV1RolesByIdMutation } from '@/api-client/@tanstack/react-query.gen';
import { useMutation } from '@tanstack/react-query';
import { useRouter } from '@tanstack/react-router';

import { toast } from 'sonner';

import type { ErrorResponse } from '@/api-client';
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog';
import { apiClient } from '@/lib/api-client';

export default function DeleteRoleByIdDialog({
  id,
  children,
}: {
  id: string;
  children: React.ReactNode;
}) {
  const router = useRouter();

  const deleteRole = useMutation({
    ...deleteApiV1RolesByIdMutation({
      client: apiClient,
    }),
  });

  return (
    <AlertDialog>
      <AlertDialogTrigger>{children}</AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>Delete Role</AlertDialogHeader>
        <AlertDialogDescription>
          Are you sure you want to delete this role? This action cannot be
          undone.
        </AlertDialogDescription>
        <AlertDialogFooter>
          <AlertDialogCancel>Cancel</AlertDialogCancel>
          <AlertDialogAction
            onClick={() =>
              toast.promise(
                deleteRole.mutateAsync({
                  path: {
                    id,
                  },
                }),
                {
                  loading: 'Removing role from the system...',
                  success: () => {
                    router.invalidate();

                    return 'The role has been removed from the system.';
                  },
                  error: (error: ErrorResponse) => error.message,
                }
              )
            }
          >
            Delete Role
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
