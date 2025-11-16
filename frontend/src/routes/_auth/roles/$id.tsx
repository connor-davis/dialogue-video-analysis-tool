import {
  Link,
  Outlet,
  createFileRoute,
  useMatchRoute,
  useParams,
} from '@tanstack/react-router';
import { ChevronLeftIcon, FlagIcon, InfoIcon } from 'lucide-react';

import { type Role, getApiV1RolesById } from '@/api-client';
import PageLoader from '@/components/messages/page-loader';
import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs';
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip';
import { apiClient } from '@/lib/api-client';

export const Route = createFileRoute('/_auth/roles/$id')({
  pendingComponent: () => <PageLoader label="Loading role..." />,
  loader: async ({ params: { id } }) => {
    const { data: roleData } = await getApiV1RolesById({
      client: apiClient,
      path: {
        id,
      },
      throwOnError: true,
    });

    return {
      role: (roleData.item ?? {}) as Role,
    };
  },
  component: RouteComponent,
});

function RouteComponent() {
  const { role } = Route.useLoaderData();
  const params = useParams({ from: '/_auth/roles/$id' });
  const matchRoute = useMatchRoute();

  // Define routes for tabs
  const tabs = [
    { value: 'details', label: 'Details', to: '/roles/$id' },
    {
      value: 'permissions',
      label: 'Permissions',
      to: '/roles/$id/permissions',
    },
  ];

  // Determine which tab is active by checking matches
  const activeTab =
    tabs.find((tab) =>
      matchRoute({ to: tab.to, params: { id: params.id }, fuzzy: false })
    )?.value ?? 'details';

  return (
    <div className="flex flex-col w-full h-full overflow-hidden p-3 bg-background rounded-2xl border gap-3">
      <Tabs defaultValue={activeTab}>
        <div className="flex items-center w-full h-auto gap-3">
          <div className="flex items-center gap-3">
            <Button variant="ghost" size="icon" asChild>
              <Link to="/roles">
                <ChevronLeftIcon />
              </Link>
            </Button>
            <Label>{role.name}</Label>
          </div>
          <div className="flex items-center gap-3 ml-auto">
            <TabsList>
              <Tooltip>
                <TooltipTrigger>
                  <TabsTrigger value="details" asChild>
                    <Link to="/roles/$id" params={{ id: params.id }}>
                      <InfoIcon />
                    </Link>
                  </TabsTrigger>
                </TooltipTrigger>
                <TooltipContent side="bottom">Details</TooltipContent>
              </Tooltip>
              <Tooltip>
                <TooltipTrigger>
                  <TabsTrigger value="permissions" asChild>
                    <Link
                      to="/roles/$id/permissions"
                      params={{ id: params.id }}
                    >
                      <FlagIcon />
                    </Link>
                  </TabsTrigger>
                </TooltipTrigger>
                <TooltipContent side="bottom">Permissions</TooltipContent>
              </Tooltip>
            </TabsList>
          </div>
        </div>

        <Outlet />
      </Tabs>
    </div>
  );
}
