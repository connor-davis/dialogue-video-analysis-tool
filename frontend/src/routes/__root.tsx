import { TanStackDevtools } from '@tanstack/react-devtools';
import { QueryClientProvider } from '@tanstack/react-query';
import { Outlet, createRootRouteWithContext } from '@tanstack/react-router';
import { TanStackRouterDevtoolsPanel } from '@tanstack/react-router-devtools';

import { ThemeProvider } from '@/components/providers/theme-provider';
import { Toaster } from '@/components/ui/sonner';
import { queryClient } from '@/lib/query-client';
import type { getRoles, getUser } from '@/lib/user';

export const Route = createRootRouteWithContext<{
  getUser: typeof getUser;
  getRoles: typeof getRoles;
}>()({
  component: () => (
    <ThemeProvider defaultAppearance="one" defaultTheme="system">
      <QueryClientProvider client={queryClient}>
        <div className="flex flex-col w-screen h-screen text-foreground bg-background font-montserrat">
          <Outlet />
          {import.meta.env.MODE === 'development' && (
            <TanStackDevtools
              config={{
                position: 'bottom-right',
              }}
              plugins={[
                {
                  name: 'Tanstack Router',
                  render: <TanStackRouterDevtoolsPanel />,
                },
              ]}
            />
          )}
        </div>

        <Toaster />
      </QueryClientProvider>
    </ThemeProvider>
  ),
});
