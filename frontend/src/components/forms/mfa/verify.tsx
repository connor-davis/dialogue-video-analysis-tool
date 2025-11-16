import { postApiV1AuthenticationMfaTotpVerifyMutation } from '@/api-client/@tanstack/react-query.gen';
import { useMutation } from '@tanstack/react-query';
import { useNavigate } from '@tanstack/react-router';
import { useForm } from 'react-hook-form';

import { toast } from 'sonner';

import type { ErrorResponse, User } from '@/api-client';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import {
  InputOTP,
  InputOTPGroup,
  InputOTPSeparator,
  InputOTPSlot,
} from '@/components/ui/input-otp';
import { apiClient } from '@/lib/api-client';
import { cn } from '@/lib/utils';

export function VerifyMfaForm({
  to,
  user,
  className,
  ...props
}: React.ComponentProps<'div'> & { to: string; user: User }) {
  const navigate = useNavigate();

  const verifyMfaForm = useForm({
    defaultValues: {
      code: '',
    },
  });

  const verifyMfaMutation = useMutation({
    ...postApiV1AuthenticationMfaTotpVerifyMutation({
      client: apiClient,
    }),
  });

  return (
    <div className={cn('flex flex-col gap-6', className)} {...props}>
      <Card className="overflow-hidden p-0">
        <CardContent className="grid p-0">
          <Form {...verifyMfaForm}>
            <form
              className="p-6 md:p-8"
              onSubmit={verifyMfaForm.handleSubmit(({ code }) =>
                toast.promise(
                  verifyMfaMutation.mutateAsync({
                    query: {
                      code,
                    },
                  }),
                  {
                    loading: 'Verifying MFA code...',
                    success: () => {
                      navigate({ to: decodeURIComponent(to) });

                      return 'MFA code verified successfully!';
                    },
                    error: (err: ErrorResponse) => err.message,
                  }
                )
              )}
            >
              <div className="flex flex-col gap-6 justify-between w-full h-full">
                <div className="flex flex-col items-center text-center">
                  <img
                    src={user.image}
                    alt="Image"
                    className="rounded-full h-48 w-48"
                  />
                  <h1 className="text-2xl font-bold">
                    Welcome back, {user.name}
                  </h1>
                </div>

                <FormField
                  control={verifyMfaForm.control}
                  name="code"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>MFA Code</FormLabel>
                      <FormControl>
                        <InputOTP maxLength={6} {...field}>
                          <InputOTPGroup>
                            <InputOTPSlot index={0} />
                            <InputOTPSlot index={1} />
                            <InputOTPSlot index={2} />
                          </InputOTPGroup>
                          <InputOTPSeparator />
                          <InputOTPGroup>
                            <InputOTPSlot index={3} />
                            <InputOTPSlot index={4} />
                            <InputOTPSlot index={5} />
                          </InputOTPGroup>
                        </InputOTP>
                      </FormControl>
                      <FormDescription>
                        Please enter the 6-digit code from your authenticator
                        app.
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <Button type="submit" className="w-full">
                  Verify
                </Button>
              </div>
            </form>
          </Form>
          {/* <div className="bg-muted relative hidden md:block">
            <img
              src="/login-banner.png"
              alt="Image"
              className="absolute inset-0 h-full w-full object-cover dark:brightness-[0.2] dark:grayscale"
            />
          </div> */}
        </CardContent>
      </Card>
    </div>
  );
}
