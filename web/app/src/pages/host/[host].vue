<script setup lang="ts">
    import type { RouteLocationRaw } from 'vue-router';
    import type { Item } from '~/components/FloatingActionButton.vue';

    definePage({
        name: 'host-layout',
        meta: {
            breadcrumb: route => route.params['host'] as string,
        },
    });

    const hostname = useRouteParams<string>('host');

    const actionItems: Item[] = [
        {
            tooltip: 'Start',
            icon: 'solar--play-line-duotone',
        },
        {
            tooltip: 'Stop',
            icon: 'solar--stop-line-duotone',
        },
        {
            tooltip: 'Restart',
            icon: 'solar--restart-line-duotone',
        },
    ];
</script>

<template>
    <h3 class="text-4xl font-subheading font-bold mb-3">
        {{ hostname }}
    </h3>

    <SimpleCard class="bg-base-200 min-h-48">
        <template #title>
            <div class="tabs tabs-box">
                <RouterLink
                    class="tab"
                    :to="{ name: 'host', params: { host: hostname } }"
                >
                    Information
                </RouterLink>

                <RouterLink
                    class="tab"
                    :to="{ name: 'host-instances', params: { host: hostname } }"
                >
                    Instances
                </RouterLink>

                <RouterLink
                    class="tab"
                    :to="{ name: 'host-logs', params: { host: hostname } }"
                >
                    Logs
                </RouterLink>

                <RouterLink
                    class="tab"
                    :to="{ name: 'host-metrics', params: { host: hostname } }"
                >
                    Metrics
                </RouterLink>
            </div>
        </template>

        <RouterView />
    </SimpleCard>

    <FloatingActionButton :items="actionItems">
        <i class="iconify solar--power-line-duotone" />
    </FloatingActionButton>
</template>
