<script setup lang="ts">
    import type { NavItem } from '../types/nav';

    const items: NavItem[] = [
        {
            id: 'hosts',
            text: 'Hosts',
            href: '/host',
        },
    ];
</script>

<template>
    <div class="drawer">
        <input
            id="drawer-toggle"
            type="checkbox"
            class="drawer-toggle"
        />

        <div class="drawer-content flex flex-col">
            <DefaultHeader :items />
            <main>
                <RouterView />
            </main>
        </div>

        <DefaultSidebar :items />
    </div>
</template>
