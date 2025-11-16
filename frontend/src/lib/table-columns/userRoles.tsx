import type { ColumnDef } from '@tanstack/react-table';
import { MinusCircleIcon } from 'lucide-react';

import type { Role } from '@/api-client';
import UnassignRoleByUserIdAndRoleId from '@/components/dialogs/users/roles/unassign';
import { Button } from '@/components/ui/button';

export const userRoleColumns: ({
  userId,
}: {
  userId: string;
}) => ColumnDef<Role>[] = ({ userId }: { userId: string }) => [
  {
    accessorKey: 'name',
    header: 'Name',
    cell: ({ row }) => <div>{row.getValue('name')}</div>,
  },
  {
    accessorKey: 'description',
    header: 'Description',
    cell: ({ row }) => <div>{row.getValue('description')}</div>,
  },
  {
    accessorKey: 'id',
    header: 'Actions',
    cell: ({ row }) => (
      <div className="flex items-center gap-3">
        <UnassignRoleByUserIdAndRoleId userId={userId} roleId={row.original.id}>
          <Button variant="destructive" size="icon">
            <MinusCircleIcon />
          </Button>
        </UnassignRoleByUserIdAndRoleId>
      </div>
    ),
  },
];
