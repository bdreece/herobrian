<script setup lang="ts">
    import type { Notification } from '~/composables/notifications';

    const notifications = ref<Notification[]>([]);
    const bus = useNotificationBus();

    bus.on(n => notifications.value.push(n));
</script>

<template>
    <details class="dropdown dropdown-end">
        <summary class="btn m-1">
            <i class="iconify solar--bell-line-duotone" />
        </summary>

        <ul
            class="menu dropdown-content bg-base-200 rounded-box z-1 w-52 p-2 shadow-sm"
        >
            <li
                v-for="notif of notifications"
                :key="notif.timestamp"
            >
                <span>{{ notif.title }}</span>
            </li>

            <li v-if="!notifications.length">
                <span>No notifications</span>
            </li>
        </ul>
    </details>
</template>
