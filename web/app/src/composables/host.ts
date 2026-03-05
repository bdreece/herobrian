import type { Host } from '~/types/host';

import { useQuery } from '@tanstack/vue-query';
import { createInjectionState } from '@vueuse/core';
import axios from 'axios';

export const [provideHosts, useHosts] = createInjectionState(() =>
    useQuery({
        queryKey: ['hosts'],
        queryFn: getHosts,
    }),
);

async function getHosts() {
    const res = await axios.get<Record<string, Host>>('/api/host');
    if (res.status !== 200) {
        throw new Error('Error fetching hosts!');
    }

    return res.data;
}
