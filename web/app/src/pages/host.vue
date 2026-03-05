<script setup lang="ts">
    definePage({
        name: 'host-layout',
        meta: {
            breadcrumb: 'Hosts',
            seo: {
                title: 'herobrian - Hosts',
            },
        },
    });

    provideHosts();
</script>

<template>
    <BreadcrumbList />
    <RouterView />
</template>
