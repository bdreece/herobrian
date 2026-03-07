<script setup lang="ts">
    import type { NavItem } from '../types/nav';

    defineProps<{
        items: NavItem[];
        toggleId: string;
    }>();

    const [authenticated] = useAuth();
</script>

<template>
    <header class="navbar bg-base-300 w-full">
        <div
            v-if="items.length"
            class="flex-none lg:hidden"
        >
            <label
                :for="toggleId"
                aria-label="open sidebar"
                class="btn btn-square btn-ghost"
            >
                <i class="iconify solar--hamburger-menu-line-duotone" />
            </label>
        </div>

        <RouterLink
            :to="{ name: 'home' }"
            class="btn btn-ghost items-center text-xl"
        >
            <i class="iconify solar--bonfire-line-duotone" />

            <span class="font-heading font-bold">herobrian</span>
        </RouterLink>

        <span class="flex-1" />

        <nav class="hidden flex-none lg:block">
            <ul class="menu menu-horizontal">
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
        </nav>

        <NotificationTray class="mr-2" />

        <ProfileMenu v-if="authenticated" />
    </header>
</template>
