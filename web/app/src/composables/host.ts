import type { Host } from '~/types/host';

import { createInjectionState } from '@vueuse/core';
import { useAxios } from '@vueuse/integrations/useAxios';

export const [provideHosts, useHosts] = createInjectionState(() =>
    useAxios<Record<string, Host>>('/api/host'),
);
