<script setup lang="ts">
    definePage({
        name: 'hosts',
        meta: {
            seo: {
                title: 'herobrian \u2014 Hosts',
            },
        },
    });

    const hosts = useHosts() ?? {};
</script>

<template>
    <h3 class="text-4xl font-subheading font-bold mb-3">Hosts</h3>

    <ul class="list bg-base-200 rounded-box shadow-md">
        <li
            v-for="[hostname] in Object.entries(hosts ?? {})"
            :key="hostname"
            class="list-row"
        >
            <RouterLink
                :to="{
                    name: 'host',
                    params: { host: hostname },
                }"
            >
                {{ hostname }}
            </RouterLink>

            <span class="flex-1" />

            <div class="flex items-center gap-2">
                <StatusIndicator />
                <span>Online</span>
            </div>
        </li>
    </ul>
</template>
