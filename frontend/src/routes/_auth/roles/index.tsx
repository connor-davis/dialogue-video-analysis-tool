import {
  Link,
  createFileRoute,
  redirect,
  useNavigate,
} from '@tanstack/react-router';
import {
  flexRender,
  getCoreRowModel,
  useReactTable,
} from '@tanstack/react-table';
import {
  ChevronDownIcon,
  FlagIcon,
  InfoIcon,
  PlusIcon,
  SearchIcon,
} from 'lucide-react';

import z from 'zod';

import { type Pagination, type Roles, getApiV1Roles } from '@/api-client';
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
import { roleColumns } from '@/lib/table-columns/roles';

export const Route = createFileRoute('/_auth/roles/')({
  validateSearch: z.object({
    page: z.coerce.number().min(1).default(1),
    pageSize: z.coerce.number().min(1).max(100).default(10),
    searchTerm: z.string().optional(),
  }),
  pendingComponent: () => <PageLoader label="Loading roles..." />,
  loaderDeps: ({ search: { page, pageSize, searchTerm } }) => ({
    page,
    pageSize,
    searchTerm,
  }),
  loader: async ({ deps: { page, pageSize, searchTerm } }) => {
    const { data: rolesData } = await getApiV1Roles({
      client: apiClient,
      query: {
        page,
        pageSize,
        preload: [],
        searchTerm,
        searchColumn: ['name', 'rolename'],
      },
      throwOnError: true,
    });

    if (rolesData.pagination) {
      if (page > rolesData.pagination.pages) {
        throw redirect({
          to: '/roles',
          search: (old) => ({
            ...old,
            page: rolesData.pagination?.pages,
          }),
        });
      }
    }

    return {
      roles: (rolesData.items ?? []) as Roles,
      pagination: (rolesData.pagination ?? {}) as Pagination,
    };
  },
  component: RouteComponent,
});

function RouteComponent() {
  const navigate = useNavigate();

  const { searchTerm } = Route.useLoaderDeps();
  const { roles, pagination } = Route.useLoaderData();

  const table = useReactTable({
    data: roles,
    columns: roleColumns,
    getCoreRowModel: getCoreRowModel(),
    pageCount: pagination.pages,
    autoResetPageIndex: false,
    manualPagination: true,
  });

  return (
    <div className="flex flex-col w-full h-full overflow-hidden p-3 bg-background rounded-2xl border gap-3">
      <div className="flex items-center w-full h-auto gap-3">
        <div className="flex items-center gap-3">
          <Label className="text-lg">Roles</Label>
        </div>
        <div className="flex items-center gap-3 ml-auto"></div>
      </div>

      <div className="flex flex-col w-full h-full gap-3">
        <div className="flex items-center gap-3">
          <InputGroup className="w-auto">
            <DebounceInputGroupInput
              placeholder="Search roles..."
              defaultValue={searchTerm ?? ''}
              onChange={(event) =>
                navigate({
                  to: '/roles',
                  search: (old) => ({
                    ...old,
                    searchTerm: event.target.value || undefined,
                  }),
                })
              }
              className="max-w-sm"
            />
            <InputGroupAddon>
              <SearchIcon />
            </InputGroupAddon>
            <InputGroupAddon align="inline-end">
              <Label className="shrink-0">{pagination.count} results...</Label>
            </InputGroupAddon>
          </InputGroup>

          <div className="flex items-center gap-3 ml-auto">
            <Tooltip>
              <TooltipTrigger>
                <Button variant="outline" size="icon" asChild>
                  <Link to="/roles/create">
                    <PlusIcon />
                  </Link>
                </Button>
              </TooltipTrigger>
              <TooltipContent side="bottom">Create Role</TooltipContent>
            </Tooltip>

            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="outline" className="ml-auto">
                  Columns <ChevronDownIcon />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                {table
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
              {table.getHeaderGroups().map((headerGroup) => (
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
              {table.getRowModel().rows?.length ? (
                table.getRowModel().rows.map((row) => (
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
                    colSpan={roleColumns.length}
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
          pagination={pagination}
          onPageChange={(page) =>
            navigate({
              to: '/roles',
              search: (old) => ({ ...old, page }),
            })
          }
          onPageSizeChange={(pageSize) =>
            navigate({
              to: '/roles',
              search: (old) => ({ ...old, pageSize }),
            })
          }
          className="mt-auto"
        />
      </div>
    </div>
  );
}
