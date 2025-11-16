import { createFileRoute, redirect } from '@tanstack/react-router';

import z from 'zod';

import { VerifyMfaForm } from '@/components/forms/mfa/verify';

export const Route = createFileRoute('/_mfa/mfa/verify')({
  validateSearch: z.object({
    to: z.string().default('/'),
  }),
  beforeLoad: async ({ context: { getUser }, search: { to } }) => {
    const { user } = await getUser();

    if (user && user.mfaEnabled && user.mfaVerified) {
      throw redirect({
        to,
      });
    }

    return {};
  },
  loader: async ({ context: { getUser } }) => {
    const { user } = await getUser();
    return { user };
  },
  component: RouteComponent,
});

function RouteComponent() {
  const { to } = Route.useSearch();
  const { user } = Route.useLoaderData();

  if (!user) return null;

  return (
    <div className="flex min-h-svh flex-col items-center justify-center p-6 md:p-10">
      <VerifyMfaForm to={to} user={user} />
    </div>
  );
}
