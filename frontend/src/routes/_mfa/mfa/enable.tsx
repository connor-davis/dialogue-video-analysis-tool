import { createFileRoute, redirect } from '@tanstack/react-router';

import z from 'zod';

import { EnableMfaForm } from '@/components/forms/mfa/enable';

export const Route = createFileRoute('/_mfa/mfa/enable')({
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
  component: RouteComponent,
});

function RouteComponent() {
  const { to } = Route.useSearch();

  return (
    <div className="flex min-h-svh flex-col items-center justify-center p-6 md:p-10">
      <EnableMfaForm to={to} />
    </div>
  );
}
