import { postApiV1AuthenticationMfaTotpVerifyMutation } from '@/api-client/@tanstack/react-query.gen';
import { useMutation } from '@tanstack/react-query';
import { useNavigate } from '@tanstack/react-router';
import { useForm } from 'react-hook-form';

import { toast } from 'sonner';

import type { ErrorResponse } from '@/api-client';
import { AspectRatio } from '@/components/ui/aspect-ratio';
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

export function EnableMfaForm({
  to,
  className,
  ...props
}: React.ComponentProps<'div'> & { to: string }) {
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
        <CardContent className="grid p-0 md:grid-cols-1">
          <div className="bg-muted relative hidden md:block">
            <AspectRatio ratio={1 / 1} className="h-full w-full">
              <img
                src={`${import.meta.env.VITE_API_BASE_URL}/api/v1/authentication/mfa/totp/enable`}
                alt="Image"
                className="absolute inset-0 h-full w-full"
              />
            </AspectRatio>
          </div>
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

                      return 'Your MFA has been enabled successfully!';
                    },
                    error: (error: ErrorResponse) => error.message,
                  }
                )
              )}
            >
              <div className="flex flex-col gap-6 justify-between h-full">
                <div className="flex flex-col items-center text-center">
                  <h1 className="text-2xl font-bold">Welcome back</h1>
                </div>

                <FormField
                  control={verifyMfaForm.control}
                  name="code"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>MFA Code</FormLabel>
                      <FormControl>
                        <InputOTP
                          maxLength={6}
                          autoComplete="one-time-code"
                          {...field}
                        >
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
        </CardContent>
      </Card>
    </div>
  );
}
