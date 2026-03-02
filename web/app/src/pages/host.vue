<script setup lang="ts">
    definePage({
        name: 'host-layout',
        meta: {
            breadcrumb: 'Hosts',
        },
    });

    provideHosts();
</script>

<template>
    <BreadcrumbList />
    <RouterView />
</template>
