<script setup lang="ts">
    import type { NavItem } from '../types/nav';

    const items: NavItem[] = [
        {
            id: 'hosts',
            text: 'Hosts',
            to: '/host',
        },
    ];
</script>

<template>
    <AppShell :items>
        <BreadcrumbList class="mb-6" />
        <div class="container mx-auto">
            <RouterView />
        </div>
    </AppShell>
</template>
