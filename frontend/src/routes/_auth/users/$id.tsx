import {
  Link,
  Outlet,
  createFileRoute,
  useMatchRoute,
  useParams,
} from '@tanstack/react-router';
import {
  BrickWallShieldIcon,
  ChevronLeftIcon,
  InfoIcon,
  ShieldIcon,
} from 'lucide-react';

import { type User, getApiV1UsersById } from '@/api-client';
import PageLoader from '@/components/messages/page-loader';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs';
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip';
import { apiClient } from '@/lib/api-client';

export const Route = createFileRoute('/_auth/users/$id')({
  pendingComponent: () => <PageLoader label="Loading user..." />,
  loader: async ({ params: { id } }) => {
    const { data: userData } = await getApiV1UsersById({
      client: apiClient,
      path: {
        id,
      },
      throwOnError: true,
    });

    return {
      user: (userData.item ?? {}) as User,
    };
  },
  component: RouteComponent,
});

function RouteComponent() {
  const { user } = Route.useLoaderData();
  const params = useParams({ from: '/_auth/users/$id' });
  const matchRoute = useMatchRoute();

  // Define routes for tabs
  const tabs = [
    { value: 'profile', label: 'Profile', to: '/users/$id' },
    { value: 'roles', label: 'Roles', to: '/users/$id/roles' },
    { value: 'security', label: 'Security', to: '/users/$id/security' },
  ];

  // Determine which tab is active by checking matches
  const activeTab =
    tabs.find((tab) =>
      matchRoute({ to: tab.to, params: { id: params.id }, fuzzy: false })
    )?.value ?? 'profile';

  return (
    <div className="flex flex-col w-full h-full overflow-hidden p-3 bg-background rounded-2xl border gap-3">
      <Tabs defaultValue={activeTab} className="h-full">
        <div className="flex items-center w-full h-auto gap-3">
          <div className="flex items-center gap-3">
            <Button variant="ghost" size="icon" asChild>
              <Link to="/users">
                <ChevronLeftIcon />
              </Link>
            </Button>
            <Label>
              <Avatar className="h-8 w-8 rounded-lg">
                <AvatarImage src={user.image ?? ''} alt={user.name} />
                <AvatarFallback className="rounded-lg">
                  {user.name.charAt(0)}
                </AvatarFallback>
              </Avatar>
              <div className="grid flex-1 text-left text-sm leading-tight">
                <span className="truncate font-medium text-lg">
                  {user.name}
                </span>
              </div>
            </Label>
          </div>
          <div className="flex items-center gap-3 ml-auto">
            <TabsList>
              <Tooltip>
                <TooltipTrigger>
                  <TabsTrigger value="profile" asChild>
                    <Link to="/users/$id" params={{ id: params.id }}>
                      <InfoIcon />
                    </Link>
                  </TabsTrigger>
                </TooltipTrigger>
                <TooltipContent side="bottom">Profile</TooltipContent>
              </Tooltip>
              <Tooltip>
                <TooltipTrigger>
                  <TabsTrigger value="roles" asChild>
                    <Link to="/users/$id/roles" params={{ id: params.id }}>
                      <BrickWallShieldIcon />
                    </Link>
                  </TabsTrigger>
                </TooltipTrigger>
                <TooltipContent side="bottom">Roles</TooltipContent>
              </Tooltip>
              <Tooltip>
                <TooltipTrigger>
                  <TabsTrigger value="security" asChild>
                    <Link to="/users/$id/security" params={{ id: params.id }}>
                      <ShieldIcon />
                    </Link>
                  </TabsTrigger>
                </TooltipTrigger>
                <TooltipContent side="bottom">Security</TooltipContent>
              </Tooltip>
            </TabsList>
          </div>
        </div>

        <Outlet />
      </Tabs>
    </div>
  );
}
