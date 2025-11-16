import {
  type ClientOptions,
  createClient,
  createConfig,
} from '@/api-client/client';

export const apiClient = createClient(
  createConfig<ClientOptions>({
    baseUrl: import.meta.env.VITE_API_BASE_URL || 'http://localhost:6173',
    credentials: 'include',
  })
);
