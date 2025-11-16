import {
  Outlet,
  createFileRoute,
  redirect,
  useRouter,
} from '@tanstack/react-router';
import { useEffect } from 'react';

import { AppSidebar } from '@/components/sidebar/app';
import { SidebarInset, SidebarProvider } from '@/components/ui/sidebar';

export const Route = createFileRoute('/_auth')({
  beforeLoad: async ({ context: { getUser }, location }) => {
    const { user, error } = await getUser();

    if (!user || error) {
      throw redirect({
        to: '/sign-in',
        search: {
          to: encodeURIComponent(location.url),
        },
      });
    }

    if (!user.mfaEnabled && location.pathname !== '/mfa/enable') {
      throw redirect({
        to: '/mfa/enable',
        search: {
          to: encodeURIComponent(location.href),
        },
      });
    }

    if (
      user.mfaEnabled &&
      !user.mfaVerified &&
      location.pathname !== '/mfa/verify'
    ) {
      throw redirect({
        to: '/mfa/verify',
        search: {
          to: encodeURIComponent(location.href),
        },
      });
    }
  },
  loader: async ({ context: { getUser } }) => {
    const { user } = await getUser();

    return { user };
  },
  component: RouteComponent,
});

function RouteComponent() {
  const router = useRouter();

  const { user } = Route.useLoaderData();

  useEffect(() => {
    const interval = setInterval(
      () => {
        router.invalidate({ filter: (route) => route.id === '/_auth' });
      },
      1 * 60 * 1000
    );

    return () => clearInterval(interval);
  }, []);

  if (!user) {
    return null;
  }

  return (
    <div className="flex w-screen h-screen overflow-hidden">
      <div className="flex flex-col w-full h-full overflow-hidden">
        <SidebarProvider>
          <AppSidebar user={user} />

          <SidebarInset className="flex flex-col w-full h-full overflow-hidden">
            <Outlet />
          </SidebarInset>
        </SidebarProvider>
      </div>
    </div>
  );
}
