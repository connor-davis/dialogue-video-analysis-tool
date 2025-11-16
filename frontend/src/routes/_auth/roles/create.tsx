import { postApiV1RolesMutation } from '@/api-client/@tanstack/react-query.gen';
import { useMutation } from '@tanstack/react-query';
import { Link, createFileRoute, useRouter } from '@tanstack/react-router';
import { ChevronLeftIcon } from 'lucide-react';
import { useForm } from 'react-hook-form';

import { zodResolver } from '@hookform/resolvers/zod';
import { toast } from 'sonner';
import z from 'zod';

import type { ErrorResponse } from '@/api-client';
import { zCreateRolePayload } from '@/api-client/zod.gen';
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
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { apiClient } from '@/lib/api-client';

export const Route = createFileRoute('/_auth/roles/create')({
  component: RouteComponent,
});

function RouteComponent() {
  const router = useRouter();

  const createForm = useForm<z.infer<typeof zCreateRolePayload>>({
    resolver: zodResolver(zCreateRolePayload),
    defaultValues: {
      name: undefined,
      description: undefined,
    },
  });

  const createRole = useMutation({
    ...postApiV1RolesMutation({
      client: apiClient,
    }),
  });

  return (
    <div className="flex flex-col w-full h-full overflow-hidden p-3 bg-background rounded-2xl border gap-3">
      <div className="flex w-full h-auto items-center gap-3">
        <div className="flex items-center gap-3">
          <Button variant="ghost" size="icon" asChild>
            <Link to="/roles">
              <ChevronLeftIcon />
            </Link>
          </Button>

          <Label className="text-lg">Create Role</Label>
        </div>
        <div className="flex items-center gap-3"></div>
      </div>

      <Form {...createForm}>
        <form
          onSubmit={createForm.handleSubmit((values) =>
            toast.promise(createRole.mutateAsync({ body: values }), {
              loading: 'Creating the new role...',
              success: () => {
                router.invalidate();

                return 'The role has been created.';
              },
              error: (error: ErrorResponse) => error.message,
            })
          )}
          className="flex flex-col gap-3"
        >
          <FormField
            control={createForm.control}
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
            control={createForm.control}
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

          <Button type="submit" disabled={createRole.isPending}>
            {createRole.isPending ? 'Creating...' : 'Create Role'}
          </Button>
        </form>
      </Form>
    </div>
  );
}
