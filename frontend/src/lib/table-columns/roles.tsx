import { type ColumnDef } from '@tanstack/react-table';
import { TrashIcon } from 'lucide-react';

import type { Role } from '@/api-client';
import DeleteRoleByIdDialog from '@/components/dialogs/roles/delete';
import { Button } from '@/components/ui/button';

export const roleColumns: ColumnDef<Role>[] = [
  {
    accessorKey: 'name',
    header: 'Name',
    cell: ({ row }) => <div>{row.getValue('name')}</div>,
  },
  {
    accessorKey: 'description',
    header: 'Description',
    cell: ({ row }) => (
      <div className="text-ellipsis">{row.getValue('description')}</div>
    ),
  },
  {
    accessorKey: 'id',
    header: 'Actions',
    cell: ({ row }) => (
      <div className="flex items-center gap-3">
        <DeleteRoleByIdDialog id={row.original.id}>
          <Button variant="destructive" size="icon">
            <TrashIcon />
          </Button>
        </DeleteRoleByIdDialog>
      </div>
    ),
  },
];
