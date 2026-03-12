<script setup lang="ts">
    import type { NavItem } from '../types/nav';

    const { items = [] } = defineProps<{
        items?: NavItem[];
        sidebarHidden?: boolean;
    }>();

    const toggleId = useId();
</script>

<template>
    <div class="drawer">
        <input
            :id="toggleId"
            type="checkbox"
            class="drawer-toggle"
        />

        <section class="drawer-content flex flex-col h-screen">
            <AppShellHeader
                :items
                :toggle-id
            />

            <slot />

            <AppShellFooter />
        </section>

        <AppShellSidebar
            v-if="!sidebarHidden"
            :items
            :toggle-id
        />
    </div>
</template>
