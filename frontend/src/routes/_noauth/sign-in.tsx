import { createFileRoute } from '@tanstack/react-router';

import z from 'zod';

import { SignInForm } from '@/components/forms/sign-in';

export const Route = createFileRoute('/_noauth/sign-in')({
  validateSearch: z.object({
    to: z.string().default(encodeURIComponent(location.href)),
  }),
  component: RouteComponent,
});

function RouteComponent() {
  const { to } = Route.useSearch();

  return (
    <div className="flex min-h-svh flex-col items-center justify-center p-6 md:p-10">
      <SignInForm to={to} />
    </div>
  );
}
