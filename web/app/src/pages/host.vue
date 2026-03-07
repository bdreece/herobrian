<script setup lang="ts">
    definePage({
        name: 'hosts-layout',
        meta: {
            breadcrumb: 'Hosts',
        },
    });

    provideHostQuery();
</script>

<template>
    <RouterView />
</template>
