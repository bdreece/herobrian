<script setup lang="ts">
    import type { RouteLocationRaw } from 'vue-router';

    export interface Item {
        text: string;
        to?: RouteLocationRaw;
    }

    const route = useRoute();

    const items = computed<Item[]>(() =>
        route.matched
            .filter(r => r.meta.breadcrumb)
            .map(r => ({
                text: r.meta.breadcrumb!,
                to: r.path,
            })),
    );
</script>

<template>
    <div class="breadcrumbs text-sm px-2">
        <ul>
            <li
                v-for="item of items"
                :key="item.text"
            >
                <RouterLink
                    v-if="item.to"
                    :to="item.to"
                >
                    {{ item.text }}
                </RouterLink>

                <span v-else>
                    {{ item.text }}
                </span>
            </li>
        </ul>
    </div>
</template>
