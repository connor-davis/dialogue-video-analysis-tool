import {
  type Roles,
  type User,
  getApiV1AuthenticationMe,
  getApiV1UsersByUserIdListRoles,
} from '@/api-client';

import { apiClient } from './api-client';

export async function getUser(): Promise<{
  user?: User;
  error: boolean;
}> {
  const response = await getApiV1AuthenticationMe({
    client: apiClient,
  });

  const { data, error } = response;

  return {
    user: data?.item as User,
    error: error !== undefined,
  };
}

export async function getRoles(userId: string): Promise<{
  roles: Roles;
  error: boolean;
}> {
  const response = await getApiV1UsersByUserIdListRoles({
    client: apiClient,
    path: {
      userId,
    },
  });

  const { data, error } = response;

  return {
    roles: (data?.items ?? []) as Roles,
    error: error !== undefined,
  };
}
