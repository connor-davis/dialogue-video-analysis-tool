import { postApiV1UsersUnassignRoleByUserIdByRoleIdMutation } from '@/api-client/@tanstack/react-query.gen';
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

export default function UnassignRoleByUserIdAndRoleId({
  userId,
  roleId,
  children,
}: {
  userId: string;
  roleId: string;
  children: React.ReactNode;
}) {
  const router = useRouter();

  const unassignRole = useMutation({
    ...postApiV1UsersUnassignRoleByUserIdByRoleIdMutation({
      client: apiClient,
    }),
  });

  return (
    <AlertDialog>
      <AlertDialogTrigger>{children}</AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>Unassign Role</AlertDialogHeader>
        <AlertDialogDescription>
          Are you sure you want to unassign this role? You can always assign it
          again later.
        </AlertDialogDescription>
        <AlertDialogFooter>
          <AlertDialogCancel>Cancel</AlertDialogCancel>
          <AlertDialogAction
            onClick={() =>
              toast.promise(
                unassignRole.mutateAsync({
                  path: {
                    userId,
                    roleId,
                  },
                }),
                {
                  loading: 'Unassigning the role from the user...',
                  success: () => {
                    router.invalidate();

                    return 'The role has been unassigned from the user.';
                  },
                  error: (error: ErrorResponse) => error.message,
                }
              )
            }
          >
            Unassign Role
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
