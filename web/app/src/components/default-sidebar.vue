<script setup lang="ts">
    import type { NavItem } from '../types/nav';

    defineProps<{
        items: NavItem[];
    }>();
</script>

<template>
    <div class="drawer-side">
        <label
            for="drawer-toggle"
            aria-label="close sidebar"
            class="drawer-overlay"
        ></label>

        <ul class="menu bg-base-200 min-h-full w-80 p-4">
            <li
                v-for="item of items"
                :key="item.id"
            >
                <RouterLink
                    :to="item.to"
                    class="flex gap-1"
                >
                    <i
                        v-if="item.icon"
                        class="iconify"
                        :class="item.icon"
                    />
                    <span v-text="item.text" />
                </RouterLink>
            </li>
        </ul>
    </div>
</template>
