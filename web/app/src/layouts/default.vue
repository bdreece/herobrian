<script setup lang="ts">
    import type { NavItem } from '../types/nav';

    const items: NavItem[] = [
        {
            id: 'hosts',
            text: 'Hosts',
            to: '/host',
        },
        {
            id: 'users',
            text: 'Users',
            to: '/user',
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
