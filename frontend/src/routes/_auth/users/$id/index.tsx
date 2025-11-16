import { putApiV1UsersByIdMutation } from '@/api-client/@tanstack/react-query.gen';
import { useMutation } from '@tanstack/react-query';
import { createFileRoute, useRouter } from '@tanstack/react-router';
import { useForm } from 'react-hook-form';

import { zodResolver } from '@hookform/resolvers/zod';
import { toast } from 'sonner';
import z from 'zod';

import { type ErrorResponse, type User, getApiV1UsersById } from '@/api-client';
import { zUpdateUserPayload } from '@/api-client/zod.gen';
import { AvatarCropper } from '@/components/avatar-cropper';
import PageLoader from '@/components/messages/page-loader';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogTrigger } from '@/components/ui/dialog';
import {
  Field,
  FieldDescription,
  FieldGroup,
  FieldLabel,
} from '@/components/ui/field';
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

export const Route = createFileRoute('/_auth/users/$id/')({
  pendingComponent: () => <PageLoader label="Loading user..." />,
  loader: async ({ params: { id } }) => {
    const { data: userData } = await getApiV1UsersById({
      client: apiClient,
      path: {
        id,
      },
      throwOnError: true,
    });

    return {
      user: (userData.item ?? {}) as User,
    };
  },
  component: RouteComponent,
});

function RouteComponent() {
  const router = useRouter();

  const { id } = Route.useParams();
  const { user } = Route.useLoaderData();

  const updateForm = useForm<z.infer<typeof zUpdateUserPayload>>({
    resolver: zodResolver(zUpdateUserPayload),
    defaultValues: {
      image: user.image,
      name: user.name,
      username: user.username,
      bio: user.bio,
    },
  });

  const updateUser = useMutation({
    ...putApiV1UsersByIdMutation({
      client: apiClient,
    }),
  });

  return (
    <div className="flex flex-col w-full h-auto gap-3">
      <Form {...updateForm}>
        <form
          onSubmit={updateForm.handleSubmit((values) =>
            toast.promise(
              updateUser.mutateAsync({
                path: {
                  id,
                },
                body: values,
              }),
              {
                loading: 'Updating your profile...',
                success: () => {
                  router.invalidate();

                  return 'Your profile has been updated.';
                },
                error: (error: ErrorResponse) => error.message,
              }
            )
          )}
          className="grid grid-cols-1 md:grid-cols-2 gap-3"
        >
          <div className="flex flex-col w-full items-center justify-center h-auto gap-3">
            <FormField
              control={updateForm.control}
              name="image"
              render={({ field }) => (
                <FormItem className="flex flex-col w-full h-auto items-center justify-center">
                  <FormLabel>Profile Image</FormLabel>
                  <FormControl>
                    <Dialog>
                      <DialogTrigger asChild>
                        <Avatar className="w-full h-auto max-w-64 min-w-32 min-h-32">
                          <AvatarImage
                            src={field.value}
                            alt={user.name.charAt(0)}
                            className=" hover:cursor-pointer"
                          />
                          <AvatarFallback>{user.name.charAt(0)}</AvatarFallback>
                        </Avatar>
                      </DialogTrigger>
                      <DialogContent>
                        <div className="flex flex-col w-full h-auto gap-3">
                          <FieldGroup>
                            <AvatarCropper
                              src={field.value || ''}
                              onComplete={(dataUrl) => {
                                field.onChange(dataUrl);
                              }}
                            />
                            <Field>
                              <FieldLabel>Choose Image</FieldLabel>
                              <Input
                                placeholder="https://example.com/my-avatar.png"
                                type="file"
                                accept="image/png, image/jpeg, image/jpg, image/webp, image/gif"
                                onChange={(e) => {
                                  const file = e.target.files?.[0];
                                  if (file) {
                                    const reader = new FileReader();
                                    reader.onloadend = () => {
                                      field.onChange(reader.result as string);
                                    };
                                    reader.readAsDataURL(file);
                                  }
                                }}
                              />
                              <FieldDescription>
                                Choose a profile image from your device.
                                Supported formats: PNG, JPEG, JPG, WEBP, GIF.
                              </FieldDescription>
                            </Field>
                          </FieldGroup>
                        </div>
                      </DialogContent>
                    </Dialog>
                  </FormControl>
                  <FormDescription>
                    Click on the avatar to change your profile image.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
          </div>
          <div className="flex flex-col w-full h-auto gap-5">
            <FormField
              control={updateForm.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Name</FormLabel>
                  <FormControl>
                    <Input placeholder="John Doe" {...field} />
                  </FormControl>
                  <FormDescription>
                    This is your display name visible to other users.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={updateForm.control}
              name="username"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Username</FormLabel>
                  <FormControl>
                    <Input placeholder="user@domain.com" {...field} />
                  </FormControl>
                  <FormDescription>
                    This is your unique username used for logging in.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={updateForm.control}
              name="bio"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Bio</FormLabel>
                  <FormControl>
                    <Textarea
                      placeholder="Tell us about yourself..."
                      {...field}
                      value={field.value ?? undefined}
                    />
                  </FormControl>
                  <FormDescription>
                    This is your bio, visible to other users.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <Button type="submit" className="w-full">
              Update Details
            </Button>
          </div>
        </form>
      </Form>
    </div>
  );
}
