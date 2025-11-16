import { deleteApiV1UsersByIdMutation } from '@/api-client/@tanstack/react-query.gen';
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

export default function DeleteUserByIdDialog({
  id,
  children,
}: {
  id: string;
  children: React.ReactNode;
}) {
  const router = useRouter();

  const deleteUser = useMutation({
    ...deleteApiV1UsersByIdMutation({
      client: apiClient,
    }),
  });

  return (
    <AlertDialog>
      <AlertDialogTrigger>{children}</AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>Delete User</AlertDialogHeader>
        <AlertDialogDescription>
          Are you sure you want to delete this user? This action cannot be
          undone.
        </AlertDialogDescription>
        <AlertDialogFooter>
          <AlertDialogCancel>Cancel</AlertDialogCancel>
          <AlertDialogAction
            onClick={() =>
              toast.promise(
                deleteUser.mutateAsync({
                  path: {
                    id,
                  },
                }),
                {
                  loading: 'Removing user from the system...',
                  success: () => {
                    router.invalidate();

                    return 'The user has been removed from the system.';
                  },
                  error: (error: ErrorResponse) => error.message,
                }
              )
            }
          >
            Delete User
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
