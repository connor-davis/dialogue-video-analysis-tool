import { putApiV1RolesByIdMutation } from '@/api-client/@tanstack/react-query.gen';
import { useMutation } from '@tanstack/react-query';
import { createFileRoute, useRouter } from '@tanstack/react-router';
import { useForm } from 'react-hook-form';

import { zodResolver } from '@hookform/resolvers/zod';
import { toast } from 'sonner';
import z from 'zod';

import { type ErrorResponse, type Role, getApiV1RolesById } from '@/api-client';
import { zUpdateRolePayload } from '@/api-client/zod.gen';
import PageLoader from '@/components/messages/page-loader';
import { Button } from '@/components/ui/button';
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { apiClient } from '@/lib/api-client';

export const Route = createFileRoute('/_auth/roles/$id/')({
  pendingComponent: () => <PageLoader label="Loading role..." />,
  loader: async ({ params: { id } }) => {
    const { data: roleData } = await getApiV1RolesById({
      client: apiClient,
      path: {
        id,
      },
      throwOnError: true,
    });

    return {
      role: (roleData.item ?? {}) as Role,
    };
  },
  component: RouteComponent,
});

function RouteComponent() {
  const router = useRouter();

  const { id } = Route.useParams();
  const { role } = Route.useLoaderData();

  const updateForm = useForm<z.infer<typeof zUpdateRolePayload>>({
    resolver: zodResolver(zUpdateRolePayload),
    defaultValues: {
      name: role.name,
      description: role.description,
    },
  });

  const updateRole = useMutation({
    ...putApiV1RolesByIdMutation({
      client: apiClient,
    }),
  });

  return (
    <div className="flex flex-col w-full h-auto gap-3">
      <Form {...updateForm}>
        <form
          onSubmit={updateForm.handleSubmit((values) =>
            toast.promise(
              updateRole.mutateAsync({
                path: {
                  id,
                },
                body: values,
              }),
              {
                loading: 'Updating the role...',
                success: () => {
                  router.invalidate();

                  return 'The role has been updated.';
                },
                error: (error: ErrorResponse) => error.message,
              }
            )
          )}
          className="flex flex-col gap-3"
        >
          <FormField
            control={updateForm.control}
            name="name"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Name</FormLabel>
                <FormControl>
                  <Input type="text" placeholder="Name" {...field} />
                </FormControl>
                <FormDescription>
                  The name of the role (e.g., Administrator, Editor, Viewer).
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={updateForm.control}
            name="description"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Description</FormLabel>
                <FormControl>
                  <Textarea placeholder="Description" {...field} />
                </FormControl>
                <FormDescription>
                  A brief description of the role's responsibilities.
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <Button type="submit" disabled={updateRole.isPending}>
            {updateRole.isPending ? 'Updating...' : 'Update Role'}
          </Button>
        </form>
      </Form>
    </div>
  );
}
