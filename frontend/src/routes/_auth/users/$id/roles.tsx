import { postApiV1UsersAssignRoleByUserIdByRoleIdMutation } from '@/api-client/@tanstack/react-query.gen';
import { useMutation } from '@tanstack/react-query';
import {
  Link,
  createFileRoute,
  redirect,
  useNavigate,
  useRouter,
} from '@tanstack/react-router';
import {
  flexRender,
  getCoreRowModel,
  useReactTable,
} from '@tanstack/react-table';
import {
  CheckIcon,
  ChevronDownIcon,
  FlagIcon,
  InfoIcon,
  PlusIcon,
  SearchIcon,
} from 'lucide-react';

import { toast } from 'sonner';
import z from 'zod';

import {
  type ErrorResponse,
  type Pagination,
  type Roles,
  getApiV1Roles,
  getApiV1UsersByUserIdListRoles,
} from '@/api-client';
import PageLoader from '@/components/messages/page-loader';
import { DataTablePagination } from '@/components/table/pagination';
import { Button } from '@/components/ui/button';
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuTrigger,
} from '@/components/ui/context-menu';
import { DebounceInputGroupInput } from '@/components/ui/debounce-input';
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { InputGroup, InputGroupAddon } from '@/components/ui/input-group';
import { Label } from '@/components/ui/label';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip';
import { apiClient } from '@/lib/api-client';
import { userRoleColumns } from '@/lib/table-columns/userRoles';

export const Route = createFileRoute('/_auth/users/$id/roles')({
  validateSearch: z.object({
    userRolesPage: z.coerce.number().min(1).default(1),
    userRolesPageSize: z.coerce.number().min(1).max(100).default(10),
    userRolesSearchTerm: z.string().optional(),
    rolesPage: z.coerce.number().min(1).default(1),
    rolesPageSize: z.coerce.number().min(1).max(100).default(10),
    rolesSearchTerm: z.string().optional(),
  }),
  pendingComponent: () => <PageLoader label="Loading roles..." />,
  loaderDeps: ({
    search: {
      userRolesPage,
      userRolesPageSize,
      userRolesSearchTerm,
      rolesPage,
      rolesPageSize,
      rolesSearchTerm,
    },
  }) => ({
    userRolesPage,
    userRolesPageSize,
    userRolesSearchTerm,
    rolesPage,
    rolesPageSize,
    rolesSearchTerm,
  }),
  loader: async ({
    params: { id },
    deps: {
      userRolesPage,
      userRolesPageSize,
      userRolesSearchTerm,
      rolesPage,
      rolesPageSize,
      rolesSearchTerm,
    },
  }) => {
    const { data: userRolesData } = await getApiV1UsersByUserIdListRoles({
      client: apiClient,
      path: {
        userId: id,
      },
      query: {
        page: userRolesPage,
        pageSize: userRolesPageSize,
        preload: [],
        searchTerm: userRolesSearchTerm,
        searchColumn: ['name', 'description'],
      },
      throwOnError: true,
    });

    if (userRolesData.pagination) {
      if (userRolesPage > userRolesData.pagination.pages) {
        throw redirect({
          to: '/roles',
          search: (old) => ({
            ...old,
            userRolesPage: userRolesData.pagination?.pages,
          }),
        });
      }
    }

    const { data: rolesData } = await getApiV1Roles({
      client: apiClient,
      query: {
        page: rolesPage,
        pageSize: rolesPageSize,
        preload: [],
        searchTerm: rolesSearchTerm,
        searchColumn: ['name', 'description'],
      },
      throwOnError: true,
    });

    if (rolesData.pagination) {
      if (rolesPage > rolesData.pagination.pages) {
        throw redirect({
          to: '/roles',
          search: (old) => ({
            ...old,
            rolesPage: rolesData.pagination?.pages,
          }),
        });
      }
    }

    return {
      userRoles: (userRolesData.items ?? []) as Roles,
      userRolesPagination: (userRolesData.pagination ?? {}) as Pagination,
      roles: (rolesData.items ?? []) as Roles,
      rolesPagination: (rolesData.pagination ?? {}) as Pagination,
    };
  },
  component: RouteComponent,
});

function RouteComponent() {
  const router = useRouter();
  const navigate = useNavigate();

  const { id } = Route.useParams();
  const { userRolesSearchTerm, rolesSearchTerm } = Route.useLoaderDeps();
  const { userRoles, userRolesPagination, roles, rolesPagination } =
    Route.useLoaderData();

  const userRolesTable = useReactTable({
    data: userRoles,
    columns: userRoleColumns({ userId: id }),
    getCoreRowModel: getCoreRowModel(),
    pageCount: userRolesPagination.pages,
    autoResetPageIndex: false,
    manualPagination: true,
  });

  const assignRole = useMutation({
    ...postApiV1UsersAssignRoleByUserIdByRoleIdMutation({
      client: apiClient,
    }),
  });

  return (
    <div className="flex flex-col w-full h-full gap-3">
      <div className="flex items-center gap-3">
        <InputGroup className="w-auto">
          <DebounceInputGroupInput
            placeholder="Search roles..."
            defaultValue={userRolesSearchTerm ?? ''}
            onChange={(event) =>
              navigate({
                to: '/users/$id/roles',
                params: { id },
                search: (old) => ({
                  ...old,
                  userRolesSearchTerm: event.target.value || undefined,
                }),
              })
            }
            className="max-w-sm"
          />
          <InputGroupAddon>
            <SearchIcon />
          </InputGroupAddon>
          <InputGroupAddon align="inline-end">
            <Label className="shrink-0">
              {userRolesPagination.count} results...
            </Label>
          </InputGroupAddon>
        </InputGroup>

        <div className="flex items-center gap-3 ml-auto">
          <Tooltip>
            <Popover>
              <TooltipTrigger asChild>
                <PopoverTrigger asChild>
                  <Button variant="outline" size="icon">
                    <PlusIcon />
                  </Button>
                </PopoverTrigger>
              </TooltipTrigger>
              <PopoverContent className="w-96 p-3">
                <div className="flex flex-col w-full h-auto gap-3">
                  <div className="flex flex-col w-full h-auto">
                    <Label>Available Roles</Label>
                    <Label className="text-muted-foreground text-sm font-normal">
                      Assign roles to the user by selecting from the list below.
                    </Label>
                  </div>

                  <InputGroup>
                    <DebounceInputGroupInput
                      placeholder="Search roles..."
                      defaultValue={rolesSearchTerm ?? ''}
                      onChange={(event) =>
                        navigate({
                          to: '/users/$id/roles',
                          params: { id },
                          search: (old) => ({
                            ...old,
                            rolesSearchTerm: event.target.value || undefined,
                          }),
                        })
                      }
                      className="max-w-sm"
                    />
                    <InputGroupAddon>
                      <SearchIcon />
                    </InputGroupAddon>
                    <InputGroupAddon align="inline-end">
                      <Label className="shrink-0">
                        {rolesPagination.count} results...
                      </Label>
                    </InputGroupAddon>
                  </InputGroup>

                  <div className="flex flex-col w-full h-auto gap-1 overflow-y-auto max-h-[300px]">
                    {roles.length ? (
                      roles.map((role) => (
                        <ContextMenu key={role.id}>
                          <ContextMenuTrigger asChild>
                            <Button
                              variant="outline"
                              className="justify-between"
                              disabled={userRoles.some(
                                (userRole) => userRole.id === role.id
                              )}
                              onClick={() =>
                                toast.promise(
                                  assignRole.mutateAsync({
                                    path: {
                                      userId: id,
                                      roleId: role.id,
                                    },
                                  }),
                                  {
                                    loading:
                                      'Assigning the role to the user...',
                                    success: () => {
                                      router.invalidate();

                                      return 'The role has been assigned to the user.';
                                    },
                                    error: (error: ErrorResponse) =>
                                      error.message,
                                  }
                                )
                              }
                            >
                              <span className="capitalize">{role.name}</span>
                              {userRoles.some(
                                (userRole) => userRole.id === role.id
                              ) && (
                                <span className="text-sm text-green-500">
                                  <CheckIcon />
                                </span>
                              )}
                            </Button>
                          </ContextMenuTrigger>
                          <ContextMenuContent className="w-auto overflow-hidden">
                            <ContextMenuItem asChild>
                              <Link to="/roles/$id" params={{ id: role.id }}>
                                <InfoIcon />
                                <span>Details</span>
                              </Link>
                            </ContextMenuItem>
                            <ContextMenuItem asChild>
                              <Link
                                to="/roles/$id/permissions"
                                params={{ id: role.id }}
                              >
                                <FlagIcon />
                                <span>Permissions</span>
                              </Link>
                            </ContextMenuItem>
                          </ContextMenuContent>
                        </ContextMenu>
                      ))
                    ) : (
                      <div className="p-4 text-center text-sm text-muted-foreground">
                        No roles found.
                      </div>
                    )}
                  </div>
                </div>
              </PopoverContent>
            </Popover>
            <TooltipContent side="bottom">Assign Role</TooltipContent>
          </Tooltip>

          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline">
                Columns <ChevronDownIcon />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              {userRolesTable
                .getAllColumns()
                .filter((column) => column.getCanHide())
                .map((column) => {
                  return (
                    <DropdownMenuCheckboxItem
                      key={column.id}
                      className="capitalize"
                      checked={column.getIsVisible()}
                      onCheckedChange={(value) =>
                        column.toggleVisibility(!!value)
                      }
                    >
                      {column.id}
                    </DropdownMenuCheckboxItem>
                  );
                })}
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>

      <div className="overflow-hidden rounded-xl border">
        <Table>
          <TableHeader className="bg-accent/30">
            {userRolesTable.getHeaderGroups().map((headerGroup) => (
              <TableRow key={headerGroup.id}>
                {headerGroup.headers.map((header) => {
                  return (
                    <TableHead key={header.id}>
                      {header.isPlaceholder
                        ? null
                        : flexRender(
                            header.column.columnDef.header,
                            header.getContext()
                          )}
                    </TableHead>
                  );
                })}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {userRolesTable.getRowModel().rows?.length ? (
              userRolesTable.getRowModel().rows.map((row) => (
                <ContextMenu key={row.id}>
                  <ContextMenuTrigger asChild>
                    <TableRow>
                      {row.getVisibleCells().map((cell) => (
                        <TableCell key={cell.id}>
                          {flexRender(
                            cell.column.columnDef.cell,
                            cell.getContext()
                          )}
                        </TableCell>
                      ))}
                    </TableRow>
                  </ContextMenuTrigger>
                  <ContextMenuContent className="w-auto overflow-hidden">
                    <ContextMenuItem asChild>
                      <Link to="/roles/$id" params={{ id: row.original.id }}>
                        <InfoIcon />
                        <span>Details</span>
                      </Link>
                    </ContextMenuItem>
                    <ContextMenuItem asChild>
                      <Link
                        to="/roles/$id/permissions"
                        params={{ id: row.original.id }}
                      >
                        <FlagIcon />
                        <span>Permissions</span>
                      </Link>
                    </ContextMenuItem>
                  </ContextMenuContent>
                </ContextMenu>
              ))
            ) : (
              <TableRow>
                <TableCell
                  colSpan={userRoleColumns({ userId: id }).length}
                  className="h-24 text-center"
                >
                  No results.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>

      <DataTablePagination
        pagination={userRolesPagination}
        onPageChange={(page) =>
          navigate({
            to: '/users/$id/roles',
            params: { id },
            search: (old) => ({ ...old, userRolesPage: page }),
          })
        }
        onPageSizeChange={(pageSize) =>
          navigate({
            to: '/users/$id/roles',
            params: { id },
            search: (old) => ({ ...old, userRolesPageSize: pageSize }),
          })
        }
        className="mt-auto"
      />
    </div>
  );
}
