<script setup lang="ts">
    definePage({
        name: 'hosts',
        meta: {},
    });

    const { data } = useHosts() ?? {};
</script>

<template>
    <ul class="list bg-base-100 rounded-box shadow-md">
        <li
            v-for="[name, host] in Object.entries(data ?? {})"
            :key="name"
            class="list-row"
        >
            <div>{{ name }}</div>
            <div>{{ JSON.stringify(host) }}</div>
        </li>
    </ul>
</template>
