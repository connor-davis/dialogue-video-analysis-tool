import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import {
  Field,
  FieldDescription,
  FieldGroup,
  FieldSeparator,
} from '@/components/ui/field';
import { cn } from '@/lib/utils';

export function SignInForm({
  to = encodeURIComponent(location.href),
  className,
  ...props
}: React.ComponentProps<'div'> & { to?: string }) {
  return (
    <div className={cn('flex flex-col gap-6', className)} {...props}>
      <Card className="overflow-hidden p-0">
        <CardContent className="grid p-0 md:grid-cols-1">
          <form className="p-6 md:p-8">
            <FieldGroup>
              <div className="flex flex-col items-center gap-2 text-center">
                <h1 className="text-2xl font-bold">Welcome back</h1>
                <p className="text-muted-foreground text-balance">
                  Login to your <span className="font-bold">Thusa One</span>{' '}
                  account to continue.
                </p>
              </div>
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
