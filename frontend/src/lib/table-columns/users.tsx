import { type ColumnDef } from '@tanstack/react-table';
import { TrashIcon } from 'lucide-react';

import type { User } from '@/api-client';
import DeleteUserByIdDialog from '@/components/dialogs/users/delete';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';

export const userColumns: ColumnDef<User>[] = [
  {
    accessorKey: 'image',
    header: 'Image',
    cell: ({ row }) => (
      <div>
        <Avatar>
          <AvatarImage
            src={row.getValue('image') ?? ''}
            alt={row.getValue('name')}
          />
          <AvatarFallback className="rounded-lg">
            {(row.getValue('name') as string).charAt(0)}
          </AvatarFallback>
        </Avatar>
      </div>
    ),
  },
  {
    accessorKey: 'name',
    header: 'Name',
    cell: ({ row }) => <div>{row.getValue('name')}</div>,
  },
  {
    accessorKey: 'username',
    header: 'Username',
    cell: ({ row }) => <div>{row.getValue('username')}</div>,
  },
  {
    accessorKey: 'id',
    header: 'Actions',
    cell: ({ row }) => (
      <div className="flex items-center gap-3">
        <DeleteUserByIdDialog id={row.original.id}>
          <Button variant="destructive" size="icon">
            <TrashIcon />
          </Button>
        </DeleteUserByIdDialog>
      </div>
    ),
  },
];
