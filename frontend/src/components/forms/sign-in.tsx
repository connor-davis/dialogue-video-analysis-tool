import { useForm } from 'react-hook-form';

import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import {
  Field,
  FieldDescription,
  FieldGroup,
  FieldSeparator,
} from '@/components/ui/field';
import { cn } from '@/lib/utils';

import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '../ui/form';
import { Input } from '../ui/input';

export function SignInForm({
  to = encodeURIComponent(location.href),
  className,
  ...props
}: React.ComponentProps<'div'> & { to?: string }) {
  const signInForm = useForm({});

  return (
    <div className={cn('flex flex-col gap-6', className)} {...props}>
      <Card className="overflow-hidden p-0">
        <CardContent className="grid p-0 md:grid-cols-1">
          <form className="p-6 md:p-8">
            <FieldGroup>
              <div className="flex flex-col items-center gap-2 text-center">
                <h1 className="text-2xl font-bold">Welcome back</h1>
                <p className="text-muted-foreground text-balance">
                  Login to your{' '}
                  <span className="font-bold">
                    Dialogue Video Analysis Tool
                  </span>{' '}
                  account to continue.
                </p>
              </div>

              <Form {...signInForm}>
                <form
                  onSubmit={signInForm.handleSubmit((values) =>
                    console.log(values)
                  )}
                  className="flex flex-col w-full h-auto gap-5"
                >
                  <div className="flex flex-col w-full h-auto gap-3">
                    <FormField
                      control={signInForm.control}
                      name="email"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Email</FormLabel>
                          <FormControl>
                            <Input placeholder="Email" {...field} />
                          </FormControl>
                          <FormDescription>
                            This is your email address.
                          </FormDescription>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={signInForm.control}
                      name="password"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Password</FormLabel>
                          <FormControl>
                            <Input placeholder="Password" {...field} />
                          </FormControl>
                          <FormDescription>
                            This is your password.
                          </FormDescription>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </div>

                  <Button>Sign In</Button>
                </form>
              </Form>

              <FieldSeparator className="*:data-[slot=field-separator-content]:bg-card">
                Available Methods
              </FieldSeparator>
              <Field className="grid grid-cols-1 gap-4">
                <a
                  href={`${import.meta.env.VITE_API_BASE_URL}/api/v1/authentication/microsoft/redirect?to=${to}`}
                  className="w-full"
                >
                  <Button variant="outline" type="button" className="w-full">
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      viewBox="0 0 448 512"
                      fill="var(--foreground)"
                    >
                      <path d="M0 32l214.6 0 0 214.6-214.6 0 0-214.6zm233.4 0l214.6 0 0 214.6-214.6 0 0-214.6zM0 265.4l214.6 0 0 214.6-214.6 0 0-214.6zm233.4 0l214.6 0 0 214.6-214.6 0 0-214.6z" />
                    </svg>
                    <span className="sr-only">Login with Microsoft</span>
                    <p>Microsoft</p>
                  </Button>
                </a>
              </Field>
            </FieldGroup>
          </form>
        </CardContent>
      </Card>
      <FieldDescription className="px-6 text-center">
        By continuing, you agree to our <a href="#">Terms of Service</a> and{' '}
        <a href="#">Privacy Policy</a>.
      </FieldDescription>
    </div>
  );
}
